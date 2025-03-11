package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "vd_engine/internal/generated"
)

type server struct {
	pb.UnimplementedTgServiceServer
}

func (s *server) DownloadVideo(ctx context.Context, req *pb.DownloadVideoRequest) (*pb.DownloadVideoResponse, error) {
	return &pb.DownloadVideoResponse{Filename: "Hello, "}, nil
}

func (s *server) mustEmbedUnimplementedTgServiceServer() {
	panic("implement me")
}

func main() {
	// Создаем TCP-слушатель
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем gRPC-сервер
	s := grpc.NewServer()
	pb.RegisterTgServiceServer(s, &server{})

	log.Println("Starting server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
