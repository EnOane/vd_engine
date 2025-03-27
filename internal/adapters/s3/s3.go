package s3

import (
	"context"
	"fmt"
	"github.com/EnOane/vd_engine/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"io"
)

type S3Interface interface {
	UploadStream(ctx context.Context, f string, r io.Reader) (*minio.UploadInfo, error)
}

type Client struct {
	client *minio.Client
}

func NewS3Client() (S3Interface, error) {
	endpoint := fmt.Sprintf("%v:%d",
		config.AppConfig.S3Config.S3Host,
		config.AppConfig.S3Config.S3Port,
	)

	mcl, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AppConfig.S3Config.S3AccessKey, config.AppConfig.S3Config.S3SecretKey, ""),
		Secure: false,
		Region: config.AppConfig.S3Config.S3Region,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error connect to MinIO")
		return nil, err
	}

	log.Info().Msgf("S3: connected")

	return &Client{client: mcl}, nil
}

// UploadStream загружает поток в S3
func (s3 *Client) UploadStream(ctx context.Context, filename string, reader io.Reader) (*minio.UploadInfo, error) {
	m, err := s3.client.PutObject(ctx, config.AppConfig.S3Config.S3Bucket, filename, reader, -1, minio.PutObjectOptions{})
	return &m, err
}
