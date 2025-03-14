package service

import (
	dl "github.com/EnOane/cli_downloader/pkg/downloader"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
	"os"
)

func Execute(uri *url.URL) (<-chan []byte, string, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("error")
		return nil, "", err
	}

	// скачивание файла в папку temp
	fname, err := dl.DownloadVideo(uri, rootPath+"/temp")
	if err != nil {
		log.Error().Err(err).Msg("error")
		return nil, "", err
	}

	file, err := os.Open(fname)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return nil, "", err
	}

	out := make(chan []byte)

	go func() {
		defer file.Close()
		defer close(out)

		buffer := make([]byte, 1024*64)
		for {
			n, err := file.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Error().Err(err).Msg("error")
				break
			}

			// при создании слайса выделяется память общая под данные
			// и при отправке в канал слайс так и ссылается на одну область
			// соответсвенно изменение вызывает состояние гонки
			chunk := make([]byte, n)
			copy(chunk, buffer[:n])
			out <- chunk
		}
	}()

	return out, fname, nil
}
