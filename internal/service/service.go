package service

import (
	"context"
	dl "github.com/EnOane/cli_downloader/pkg/downloader"
	tgpb "github.com/EnOane/vd_engine/generated"
	"github.com/EnOane/vd_engine/internal/infr/s3"
	"github.com/EnOane/vd_engine/internal/util"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

// 50 мб лимит отправки видео в ТГ, как быть с остальными клиентами?
// отправлять ссылку на файл в s3
// если формат mp4 все-таки не поддерживается - отправлять как документ (только тг)

func Execute(
	request *tgpb.DownloadVideoStreamRequest,
	stream grpc.ServerStreamingServer[tgpb.DownloadVideoStreamResponse],
) error {
	// проверка url
	uri, err := url.Parse(request.Url)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return err
	}

	in, filenamePath, err := DownloadVideoStream(uri)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return err
	}

	// первое сообщение в потоке - имя файла
	err = sendFilename(filenamePath, stream)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return err
	}

	// второе сообщение в потоке - chunk файла
	out, err := sendChunks(in, stream)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return err
	}

	// конвейрная загрузка в s3
	err = uploadToS3(out, filenamePath)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return err
	}

	log.Info().Msg("video was uploaded to s3")

	return nil
}

// sendFilename отправка имени файла
func sendFilename(
	filenamePath string,
	stream grpc.ServerStreamingServer[tgpb.DownloadVideoStreamResponse],
) error {
	return stream.Send(&tgpb.DownloadVideoStreamResponse{
		Data: &tgpb.DownloadVideoStreamResponse_Filename{
			Filename: filepath.Base(filenamePath),
		},
	})
}

// sendChunks отправка чанками в стрим grpc
func sendChunks(in <-chan []byte, stream grpc.ServerStreamingServer[tgpb.DownloadVideoStreamResponse]) (chan []byte, error) {
	out := make(chan []byte)

	go func() {
		defer close(out)

		for data := range in {
			err := stream.Send(&tgpb.DownloadVideoStreamResponse{
				Data: &tgpb.DownloadVideoStreamResponse_Chunk{
					Chunk: data,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("error")
				return
			}

			out <- data
		}
	}()

	return out, nil
}

// uploadToS3 отправка чанками в стрим grpc
func uploadToS3(in <-chan []byte, filename string) error {
	fileName := filepath.Base(filename)

	reader := util.NewChannelReader(in)

	_, err := s3.UploadFile(context.TODO(), fileName, reader)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return err
	}

	return nil
}

// DownloadVideoStream создание потока из stdout
func DownloadVideoStream(uri *url.URL) (<-chan []byte, string, error) {
	// скачивание файла в папку temp
	ch, fname, err := dl.DownloadStreamVideo(uri)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return nil, "", err
	}
	return ch, fname, nil
}

// DownloadVideoFile создание потока из чтения файла
func DownloadVideoFile(uri *url.URL) (<-chan []byte, string, error) {
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
