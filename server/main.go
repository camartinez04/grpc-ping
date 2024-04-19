package main

import (
	_ "context"
	"flag"
	"fmt"
	"google.golang.org/grpc/peer"
	"io"
	"log"
	"net"
	"os"

	ping "github.com/camartinez04/grpc-ping"
	"google.golang.org/grpc"
)

// grpcServer is used to implement ping.PingServiceServer
type grpcServer struct {
	ping.UnimplementedPingServiceServer
}

// StreamPing is a server side streaming RPC to receive messages from the client
func (s *grpcServer) StreamPing(stream ping.PingService_StreamPingServer) error {
	p, ok := peer.FromContext(stream.Context())
	if ok {
		log.Printf("Client connected from: %s", p.Addr)
	}

	defer func() {
		if ok {
			log.Printf("Client disconnected at: %s", p.Addr)
		}
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Client closed the connection
			return nil
		}
		if err != nil {
			log.Printf("Failed to receive a message from the client: %v", err)
			return err // Return error to close the stream and handle disconnection
		}

		response := &ping.PongResponse{Message: "Ping: " + req.Message}
		if err := stream.Send(response); err != nil {
			log.Printf("Failed to send a message to the client: %v", err)
			return err // Return error to close the stream and handle disconnection
		}
	}
}

func main() {

	// get interface to the server from OS ENV variable
	serverEnv := os.Getenv("GRPC_PING_SERVER")
	portEnv := os.Getenv("GRPC_PING_PORT")

	server := flag.String("server", serverEnv, "GRPC Server Address")
	port := flag.String("port", portEnv, "GRPC Server Port")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *server == "" || *port == "" {
		fmt.Println("Missing required parameters:")
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("Streaming Ping GRPC server on %s:%s", *server, *port)

	lis, err := net.Listen("tcp", *server+":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ping.RegisterPingServiceServer(s, &grpcServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
