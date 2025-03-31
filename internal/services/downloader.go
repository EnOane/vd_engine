package services

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
	"time"
	"vd_engine/internal/core/interfaces"
)

// TODO: custom errors type
// TODO: metadata file
// TODO: вынести в const провайдеров

type Downloader struct {
	yt, vk, rt interfaces.DownloaderProvider
}

func NewDownloader(yt, vk, rt interfaces.DownloaderProvider) interfaces.Downloader {
	return &Downloader{yt, vk, rt}
}

// DownloadVideo загрузка видео с rutube, vk, youtube с сохранением файла
func (d *Downloader) DownloadVideo(videoUrl string, destPath string) (string, error) {
	provider, err := prepareProviderData(videoUrl)
	if err != nil {
		return "", err
	}
	return downloadAndSave(d, videoUrl, destPath, provider)
}

// DownloadStreamVideo загрузка видео с rutube, vk, youtube потоком
func (d *Downloader) DownloadStreamVideo(videoUrl string) (<-chan []byte, string, error) {
	provider, err := prepareProviderData(videoUrl)
	if err != nil {
		return nil, "", err
	}
	return downloadStream(d, videoUrl, provider)
}

// prepareProviderData возвращает наименование провайдера
func prepareProviderData(videoUrl string) (string, error) {
	u, err := url.Parse(videoUrl)
	if err != nil {
		return "", fmt.Errorf("error parsing url %w", err)
	}

	host := u.Host
	host = strings.ReplaceAll(host, "www.", "")

	provider := strings.Split(host, ".")[0]
	return provider, nil
}

// downloadAndSave логика скачивания и сохранения
func downloadAndSave(dl *Downloader, videoUrl string, destPath string, provider string) (string, error) {
	// время выполнения
	exStart := time.Now()

	log.Info().Msg(fmt.Sprintf("download video from '%v' has been started", provider))

	// имя сохраненного файла
	var filenamePath string
	var err error

	// эмуляция разной логики провайдеров
	switch provider {
	case "rutube":
		filenamePath, err = dl.rt.DownloadAndSave(videoUrl, destPath)
	case "vk", "vkvideo":
		filenamePath, err = dl.vk.DownloadAndSave(videoUrl, destPath)
	case "youtube":
		filenamePath, err = dl.yt.DownloadAndSave(videoUrl, destPath)
	default:
		return "", fmt.Errorf("download video from provider %v not supported %w", provider, err)
	}

	// обработка ошибок
	if err != nil {
		return "", err
	}

	log.Info().Msg(fmt.Sprintf("video was downloaded in %v to path '%v'", time.Since(exStart), filenamePath))

	return filenamePath, err
}

// downloadStream логика скачивания потоком
func downloadStream(dl *Downloader, videoUrl string, provider string) (<-chan []byte, string, error) {
	log.Info().Msg(fmt.Sprintf("download video from '%v' has been started", provider))

	var filename string
	var in <-chan []byte

	// эмуляция разной логики провайдеров
	switch provider {
	case "rutube":
		in, filename = dl.rt.DownloadStream(videoUrl)
	case "vk", "vkvideo":
		in, filename = dl.vk.DownloadStream(videoUrl)
	case "youtube":
		in, filename = dl.yt.DownloadStream(videoUrl)
	default:
		return nil, "", fmt.Errorf("download video from provider %v not supported", provider)
	}

	return in, filename, nil
}
