package main

import (
	"context"
	ping "github.com/camartinez04/grpc-ping"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	// get interface to the server from OS ENV variable
	server := os.Getenv("GRPC-SERVER")
	if server == "" {
		server = "localhost"

	}

	port := os.Getenv("GRPC-PORT")
	if port == "" {
		port = "50051"

	}

	delayStr := os.Getenv("DELAY_SECONDS")
	if delayStr == "" {
		delayStr = "5"
	}

	delaySeconds, _ := strconv.Atoi(delayStr)

	log.Printf("GRPC Client tool!")

	log.Printf("Connecting to GRPC Server on %s:%s", server, port)

	conn, err := grpc.Dial(server+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(conn)
	c := ping.NewPingServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Continuously try to send and receive messages
	for {
		stream, err := c.StreamPing(ctx)
		if err != nil {
			log.Printf("could not get grpc stream: %v", err)
			time.Sleep(5 * time.Second) // wait before trying to reconnect
			continue
		}

		err = sendAndReceiveMessages(stream, delaySeconds)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unavailable {
				log.Printf("grpc server unavailable: %v", err)
				time.Sleep(5 * time.Second) // wait before trying to reconnect
				continue
			}
			log.Fatalf("grpc stream closed from the server: %v", err)
		}
	}
}

func sendAndReceiveMessages(stream ping.PingService_StreamPingClient, delaySeconds int) error {
	for {
		if err := stream.Send(&ping.PingRequest{Message: "Pong"}); err != nil {
			return err // Return error to handle reconnect
		}
		resp, err := stream.Recv()
		if err != nil {
			return err // Return error to handle reconnect
		}
		log.Printf("Response: %s", resp.Message)

		// Wait for delaySeconds before sending the next ping
		time.Sleep(time.Duration(delaySeconds) * time.Second)
	}
}
