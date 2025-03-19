package api

import (
	"fmt"
	tgpb "github.com/EnOane/vd_engine/generated"
	"github.com/EnOane/vd_engine/internal/config"
	"github.com/EnOane/vd_engine/internal/service"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
)

type GrpcServer struct {
	tgpb.UnimplementedTgServiceServer
}

func MustConnect() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.AppConfig.GrpcConfig.ApiPort))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	// Создаем gRPC-сервер
	server := grpc.NewServer()
	tgpb.RegisterTgServiceServer(server, &GrpcServer{})

	log.Info().Msgf("Starting server on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}

// DownloadVideoStream скачивание видео потоком
func (s *GrpcServer) DownloadVideoStream(
	request *tgpb.DownloadVideoStreamRequest,
	stream grpc.ServerStreamingServer[tgpb.DownloadVideoStreamResponse],
) error {
	return service.Execute(request, stream)
}
