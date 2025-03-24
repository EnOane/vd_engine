package api

import (
	"fmt"
	tgpb "github.com/EnOane/vd_engine/generated"
	"github.com/EnOane/vd_engine/internal/config"
	"github.com/EnOane/vd_engine/internal/core/interfaces"
	"github.com/EnOane/vd_engine/internal/di"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
)

type grpcServer struct {
	tgpb.UnimplementedTgServiceServer
}

func MustConnect() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.AppConfig.GrpcConfig.ApiPort))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	// Создаем gRPC-сервер
	server := grpc.NewServer()
	tgpb.RegisterTgServiceServer(server, &grpcServer{})

	log.Info().Msgf("Starting server on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}

// DownloadVideoStream скачивание видео потоком
func (s *grpcServer) DownloadVideoStream(
	request *tgpb.DownloadVideoStreamRequest,
	stream grpc.ServerStreamingServer[tgpb.DownloadVideoStreamResponse],
) error {
	sr := di.Inject[interfaces.DownloadServiceInterface]()
	return sr.DownloadAndSendToClient(request, stream)
}
