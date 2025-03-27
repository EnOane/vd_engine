package api

import (
	"fmt"
	tgpb "github.com/EnOane/vd_engine/generated"
	"github.com/EnOane/vd_engine/internal/config"
	"github.com/EnOane/vd_engine/internal/domain/media_fetcher"
	"github.com/EnOane/vd_engine/internal/service"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
)

type grpcServer struct {
	tgpb.UnimplementedTgServiceServer
}

func NewGrpcServer(config *config.Config) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcConfig.ApiPort))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	tgpb.RegisterTgServiceServer(server, &grpcServer{})

	log.Info().Msgf("Starting server on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}

func (s *grpcServer) DownloadVideoStream(
	request *tgpb.DownloadVideoStreamRequest,
	stream grpc.ServerStreamingServer[tgpb.DownloadVideoStreamResponse],
) error {
	sr := service.Inject[media_fetcher.DownloadServiceInterface]()
	return sr.DownloadAndSendToClient(request, stream)
}
