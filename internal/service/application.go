package service

import (
	"context"
	s3client "github.com/EnOane/vd_engine/internal/adapters/s3"
	"github.com/EnOane/vd_engine/internal/config"
	"github.com/EnOane/vd_engine/internal/domain/media_fetcher"
)

type Application struct {
	Downloader media_fetcher.DownloadServiceInterface
}

func NewApplication(ctx context.Context, conf *config.Config) *Application {
	s3, err := s3client.NewS3Client(conf.S3Config)
	if err != nil {
		panic(err)
	}
	downloader := media_fetcher.NewDownloadService(s3)
	return &Application{downloader}
}
