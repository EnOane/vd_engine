package service

import (
	"context"
	"fmt"
	dl "github.com/EnOane/cli_downloader/pkg/downloader"
	tgpb "github.com/EnOane/vd_engine/generated"
	"github.com/EnOane/vd_engine/internal/infr/s3"
	"github.com/EnOane/vd_engine/internal/util"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net/url"
	"path/filepath"
)

type streamR grpc.ServerStreamingServer[tgpb.DownloadVideoStreamResponse]

type DownloadService struct {
	s3 *s3.Client
}

func NewDownloadService(s3 *s3.Client) *DownloadService {
	return &DownloadService{s3}
}

func (d *DownloadService) DownloadAndSendToClient(
	request *tgpb.DownloadVideoStreamRequest,
	stream grpc.ServerStreamingServer[tgpb.DownloadVideoStreamResponse],
) error {
	// Проверка URL
	uri, err := url.Parse(request.Url)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	in, filenamePath, err := downloadVideoStream(uri)
	if err != nil {
		return fmt.Errorf("failed to download video stream: %w", err)
	}

	// TODO:
	// проверка клиента и лимитов? нарушение, можно слать только поле лимита загрузки на клиент
	// 1. отправка в s3 синхронно до лимита
	// 2. если лимит превышен - отправка после уже ссылки
	// 3. если до лимита - отправка чанками

	// Отправка имени файла
	if err := sendFilename(filenamePath, stream); err != nil {
		return fmt.Errorf("failed to send filename: %w", err)
	}

	// Отправка чанков видео
	out, err := sendChunks(in, stream)
	if err != nil {
		return fmt.Errorf("failed to send video chunks: %w", err)
	}

	// Отправка в S3
	_, err = uploadToS3(d.s3, out, filenamePath)
	if err != nil {
		return fmt.Errorf("failed to upload video to S3: %w", err)
	}

	log.Info().Msgf("Video was successfully send to client: '%v' and uploaded to S3", request.ClientId)

	return nil
}

// downloadVideoStream создание видео потока
func downloadVideoStream(uri *url.URL) (<-chan []byte, string, error) {
	ch, fname, err := dl.DownloadStreamVideo(uri)
	if err != nil {
		return nil, "", fmt.Errorf("downloading video failed: %w", err)
	}
	return ch, fname, nil
}

// sendFilename отправка имени файла
func sendFilename(filenamePath string, stream streamR) error {
	return stream.Send(&tgpb.DownloadVideoStreamResponse{
		Data: &tgpb.DownloadVideoStreamResponse_Filename{
			Filename: filepath.Base(filenamePath),
		},
	})
}

// sendChunks отправка чанков видео в gRPC
func sendChunks(in <-chan []byte, stream streamR) (chan []byte, error) {
	out := make(chan []byte)

	go func() {
		defer close(out)

		for data := range in {
			err := stream.Send(&tgpb.DownloadVideoStreamResponse{
				Data: &tgpb.DownloadVideoStreamResponse_Chunk{
					Chunk: data,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to send chunk")
				return
			}

			out <- data
		}
	}()

	return out, nil
}

// uploadToS3 загрузка потока в S3
func uploadToS3(s3Client *s3.Client, in <-chan []byte, filename string) (*minio.UploadInfo, error) {
	fileName := filepath.Base(filename)
	reader := util.NewChannelReader(in)

	meta, err := s3Client.UploadStream(context.TODO(), fileName, reader)
	if err != nil {
		return nil, fmt.Errorf("uploading stream to S3 failed: %w", err)
	}

	return meta, nil
}
