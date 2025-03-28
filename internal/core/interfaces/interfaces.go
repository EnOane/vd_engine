package interfaces

import (
	"context"
	"github.com/EnOane/transports/pkg/grpc/vd_engine"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"io"
)

type DownloaderGrpcClientInterface interface {
	Stream() error
}

type DownloaderProvider interface {
	DownloadAndSave(videoUrl, destPath string) (string, error)
	DownloadStream(videoUrl string) (<-chan []byte, string)
}

type Downloader interface {
	DownloadVideo(videoUrl string, destPath string) (string, error)
	DownloadStreamVideo(videoUrl string) (<-chan []byte, string, error)
}

type Broker interface {
	GetStream(url string) (<-chan []byte, error)
}

type S3Interface interface {
	UploadStream(ctx context.Context, f string, r io.Reader) (*minio.UploadInfo, error)
}

type DownloadServiceInterface interface {
	DownloadAndSendToClient(
		r *vd_engine.DownloadVideoStreamRequest,
		stream grpc.ServerStreamingServer[vd_engine.DownloadVideoStreamResponse],
	) error
}
