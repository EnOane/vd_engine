package youtube

import (
	"github.com/google/uuid"
	"vd_engine/internal/core/interfaces"
)

type Service struct {
	b interfaces.Broker
}

func NewYoutubeService(b interfaces.Broker) interfaces.DownloaderProvider {
	return &Service{b}
}

func (r *Service) DownloadAndSave(videoUrl, destPath string) (string, error) {
	id := uuid.New().String()

	return id, nil
}

func (r *Service) DownloadStream(videoUrl string) (<-chan []byte, string) {
	id := uuid.New().String()

	a, _ := r.b.GetStream(videoUrl)

	return a, id + ".mp4"
}
