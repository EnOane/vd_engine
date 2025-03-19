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

var MinioClient *minio.Client

func MustConnect() {
	endpoint := fmt.Sprintf("%v:%d",
		config.AppConfig.S3Config.S3Host,
		config.AppConfig.S3Config.S3Port,
	)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AppConfig.S3Config.S3AccessKey, config.AppConfig.S3Config.S3SecretKey, ""),
		Secure: false,
		Region: config.AppConfig.S3Config.S3Region,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка при создании клиента MinIO")
	}

	MinioClient = minioClient

	log.Info().Msgf("S3: connected")
}

// UploadFile загрузка файла в хранилище
func UploadFile(ctx context.Context, filename string, reader io.Reader) (minio.UploadInfo, error) {
	return MinioClient.PutObject(ctx, config.AppConfig.S3Config.S3Bucket, filename, reader, -1, minio.PutObjectOptions{})
}
