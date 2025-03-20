package di

import (
	"github.com/EnOane/vd_engine/internal/infr/s3"
	"github.com/EnOane/vd_engine/internal/service"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

var Container *dig.Container

func MakeDIContainer() {
	Container = dig.New()

	makeProviders()
}

func makeProviders() {
	Container.Provide(s3.NewS3Client)
	Container.Provide(service.NewDownloadService)
}

func Inject[T any]() T {
	var dep T

	err := Container.Invoke(func(d T) { dep = d })
	if err != nil {
		log.Fatal().Err(err)
	}

	return dep
}
