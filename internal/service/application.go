package service

import (
	"context"
	s3client "github.com/EnOane/vd_engine/internal/adapters/s3"
	"github.com/EnOane/vd_engine/internal/domain/media_fetcher"
)

type Application struct {
	downloader media_fetcher.DownloadServiceInterface
}

func NewApplication(ctx context.Context) *Application {
	s3, err := s3client.NewS3Client()
	if err != nil {
		panic(err)
	}
	downloader := media_fetcher.NewDownloadService(s3)
	return &Application{downloader}
}
