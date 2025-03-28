package di

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
	"vd_engine/internal/core/interfaces"
	"vd_engine/internal/infr/brokers"
	"vd_engine/internal/infr/s3"
	"vd_engine/internal/services"
	"vd_engine/internal/services/providers/rutube"
	"vd_engine/internal/services/providers/vkvideo"
	"vd_engine/internal/services/providers/youtube"
)

var Container *dig.Container

func MakeDIContainer() {
	Container = dig.New()

	makeProviders()
}

func makeProviders() {
	Container.Provide(func() interfaces.S3Interface {
		return s3.NewS3Client()
	})

	Container.Provide(func() interfaces.Broker {
		return brokers.NewNatsBroker()
	})
	Container.Provide(func(l interfaces.Broker) interfaces.DownloaderProvider {
		return youtube.NewYoutubeService(l)
	})
	Container.Provide(func(l interfaces.Broker) interfaces.DownloaderProvider {
		return vkvideo.NewVkVideoService(l)
	})
	Container.Provide(func(l interfaces.Broker) interfaces.DownloaderProvider {
		return rutube.NewRutubeService(l)
	})
	Container.Provide(func(yt, vk, rt interfaces.DownloaderProvider) interfaces.Downloader {
		return services.NewDownloader(yt, vk, rt)
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
