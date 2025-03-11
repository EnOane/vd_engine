package main

import (
	"context"
	tgpb "github.com/EnOane/vd_engine/internal/generated"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	tgpb.UnimplementedTgServiceServer
}

func (s *server) DownloadVideo(ctx context.Context, req *tgpb.DownloadVideoRequest) (*tgpb.DownloadVideoResponse, error) {
	return &tgpb.DownloadVideoResponse{Filename: "Hello, "}, nil
}

func main() {
	// Создаем TCP-слушатель
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем gRPC-сервер
	s := grpc.NewServer()
	tgpb.RegisterTgServiceServer(s, &server{})

	log.Println("Starting server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
