package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"

	ping "github.com/camartinez04/grpc-ping"
	"google.golang.org/grpc"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.StreamPing(ctx)
	if err != nil {
		log.Fatalf("could not stream: %v", err)
	}

	for i := 0; i < 10; i++ {
		if err := stream.Send(&ping.PingRequest{Message: "Ping"}); err != nil {
			log.Fatalf("could not send ping: %v", err)
			return
		}
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("error receiving pong: %v", err)
			return
		}
		log.Printf("Response: %s", resp.Message)
	}
}
