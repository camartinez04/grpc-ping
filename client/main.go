package main

import (
	"context"
	"flag"
	"fmt"
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
	serverEnv := os.Getenv("GRPC_PING_SERVER")
	portEnv := os.Getenv("GRPC_PING_PORT")
	delayStrEnv := os.Getenv("GRPC_PING_DELAY")

	server := flag.String("server", serverEnv, "GRPC Server Address")
	port := flag.String("port", portEnv, "GRPC Server Port")
	delayStr := flag.String("delay", delayStrEnv, "Delay in seconds between pings")

	if *delayStr == "" {
		*delayStr = "5"

	}

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *server == "" || *port == "" || *delayStr == "" {
		fmt.Println("Missing required parameters:")
		flag.Usage()
		os.Exit(1)
	}

	delaySeconds, _ := strconv.Atoi(*delayStr)

	log.Printf("GRPC Client tool!")

	log.Printf("Connecting to GRPC Server on %s:%s", *server, *port)

	conn, err := grpc.Dial(*server+":"+*port, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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

// sendAndReceiveMessages sends and receives messages on the stream until an error occurs
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
