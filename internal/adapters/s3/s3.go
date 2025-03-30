package s3

import (
	"context"
	"fmt"
	conf "github.com/EnOane/vd_engine/internal/config"
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
	s3Conf conf.S3Config
}

func NewS3Client(s3Conf conf.S3Config) (S3Interface, error) {
	endpoint := fmt.Sprintf("%v:%d",
		s3Conf.S3Host,
		s3Conf.S3Port,
	)

	mcl, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3Conf.S3AccessKey, s3Conf.S3SecretKey, ""),
		Secure: false,
		Region: s3Conf.S3Region,
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
	m, err := s3.client.PutObject(ctx, s3.s3Conf.S3Bucket, filename, reader, -1, minio.PutObjectOptions{})
	return &m, err
}
