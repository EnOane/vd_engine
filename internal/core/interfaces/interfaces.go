package interfaces

import (
	"context"
	tgpb "github.com/EnOane/vd_engine/generated"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"io"
)

type S3Interface interface {
	UploadStream(ctx context.Context, f string, r io.Reader) (*minio.UploadInfo, error)
}

type DownloadServiceInterface interface {
	DownloadAndSendToClient(r *tgpb.DownloadVideoStreamRequest, stream grpc.ServerStreamingServer[tgpb.DownloadVideoStreamResponse]) error
}
