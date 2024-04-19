package main

import (
	_ "context"
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
	server := os.Getenv("GRPC_PING_SERVER")
	if server == "" {
		server = "localhost"

	}

	port := os.Getenv("GRPC_PING_PORT")
	if port == "" {
		port = "50051"

	}

	log.Printf("Streaming Ping GRPC server on %s:%s", server, port)

	lis, err := net.Listen("tcp", server+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ping.RegisterPingServiceServer(s, &grpcServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
