package api

import (
	tgpb "github.com/EnOane/vd_engine/generated"
	"github.com/EnOane/vd_engine/internal/service"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	tgpb.UnimplementedTgServiceServer
	Application service.Application
}

func NewGrpcService(application service.Application) *GrpcServer {
	return &GrpcServer{Application: application}
}

func (s *GrpcServer) DownloadVideoStream(
	request *tgpb.DownloadVideoStreamRequest,
	stream grpc.ServerStreamingServer[tgpb.DownloadVideoStreamResponse],
) error {
	return s.Application.Downloader.DownloadAndSendToClient(request, stream)
}
