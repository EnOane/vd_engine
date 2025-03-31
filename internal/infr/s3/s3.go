package s3

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"io"
	"vd_engine/internal/config"
	"vd_engine/internal/core/interfaces"
)

type Client struct {
	client *minio.Client
}

var minioClient *minio.Client

func NewS3Client() interfaces.S3Interface {
	return &Client{client: minioClient}
}

func MustConnect() {
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
	}

	log.Info().Msgf("S3: connected")

	minioClient = mcl
}

// UploadStream загружает поток в S3
func (s3 *Client) UploadStream(ctx context.Context, filename string, reader io.Reader) (*minio.UploadInfo, error) {
	m, err := s3.client.PutObject(ctx, config.AppConfig.S3Config.S3Bucket, filename, reader, -1, minio.PutObjectOptions{})
	return &m, err
}
