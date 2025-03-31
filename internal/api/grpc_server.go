package api

import (
	"fmt"
	"github.com/EnOane/transports/pkg/grpc/vd_engine"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"vd_engine/internal/config"
	"vd_engine/internal/core/interfaces"
	"vd_engine/internal/di"
	"vd_engine/internal/domain"
)

type grpcServer struct {
	vd_engine.UnimplementedVdEngineServiceServer
}

func MustServe() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.AppConfig.GrpcConfig.ApiPort))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	vd_engine.RegisterVdEngineServiceServer(server, &grpcServer{})

	log.Info().Msgf("Starting server on :%d", config.AppConfig.GrpcConfig.ApiPort)
	if err := server.Serve(listener); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}

// DownloadVideoStream скачивание видео потоком
func (s *grpcServer) DownloadVideoStream(
	request *vd_engine.DownloadVideoStreamRequest,
	stream grpc.ServerStreamingServer[vd_engine.DownloadVideoStreamResponse],
) error {
	ls := domain.NewDownloadService(di.Inject[interfaces.S3Interface]())
	return ls.DownloadAndSendToClient(request, stream)
}
