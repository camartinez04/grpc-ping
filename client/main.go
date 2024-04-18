package main

import (
	_ "context"
	"log"
	"net"
	"os"

	ping "github.com/camartinez04/grpc-ping"
	"google.golang.org/grpc"
)

type grpcServer struct {
	ping.UnimplementedPingServiceServer
}

func (s *grpcServer) StreamPing(stream ping.PingService_StreamPingServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		err = stream.Send(&ping.PongResponse{Message: "Pong: " + req.Message})
		if err != nil {
			return err
		}
	}
}

func main() {

	// get interface to the server from OS ENV variable
	client := os.Getenv("GRPC-CLIENT")
	if client == "" {
		client = "localhost"

	}

	port := os.Getenv("GRPC-PORT")
	if port == "" {
		port = "50051"

	}

	lis, err := net.Listen("tcp", client+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ping.RegisterPingServiceServer(s, &grpcServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
