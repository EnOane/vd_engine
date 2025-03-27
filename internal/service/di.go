package service

import (
	"github.com/EnOane/vd_engine/internal/adapters/s3"
	"github.com/EnOane/vd_engine/internal/domain/media_fetcher"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

var Container *dig.Container

func MakeDIContainer() {
	Container = dig.New()

	makeProviders()
}

func makeProviders() {
	Container.Provide(func() s3.S3Interface {
		return s3.NewS3Client()
	})
	Container.Provide(func(s s3.S3Interface) media_fetcher.DownloadServiceInterface {
		return media_fetcher.NewDownloadService(s)
	})
}

func Inject[T any]() T {
	var dep T

	err := Container.Invoke(func(d T) { dep = d })
	if err != nil {
		log.Fatal().Err(err)
	}

	return dep
}
