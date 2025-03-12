package main

import (
	tgpb "github.com/EnOane/vd_engine/generated"
	"github.com/EnOane/vd_engine/internal/api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	// Создаем gRPC-сервер
	server := grpc.NewServer()
	tgpb.RegisterTgServiceServer(server, &api.GrpcServer{})

	log.Info().Msgf("Starting server on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}
