package api

import (
	tgpb "github.com/EnOane/vd_engine/generated"
	"github.com/EnOane/vd_engine/internal/service"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net/url"
	"path/filepath"
)

type GrpcServer struct {
	tgpb.UnimplementedTgServiceServer
}

func (s *GrpcServer) DownloadVideo(
	request *tgpb.DownloadVideoRequest,
	stream grpc.ServerStreamingServer[tgpb.DownloadVideoResponse],
) error {
	uri, err := url.Parse(request.Url)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return err
	}

	chunksCh, filenamePath, err := service.Execute(uri)
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

	err = sendChunks(chunksCh, stream)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return err
	}

	return nil
}

func sendFilename(
	filenamePath string,
	stream grpc.ServerStreamingServer[tgpb.DownloadVideoResponse],
) error {
	return stream.Send(&tgpb.DownloadVideoResponse{
		Data: &tgpb.DownloadVideoResponse_Filename{
			Filename: filepath.Base(filenamePath),
		},
	})
}

func sendChunks(
	in <-chan []byte,
	stream grpc.ServerStreamingServer[tgpb.DownloadVideoResponse],
) error {
	for data := range in {
		err := stream.Send(&tgpb.DownloadVideoResponse{
			Data: &tgpb.DownloadVideoResponse_Chunk{
				Chunk: data,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("error")
			return err
		}
	}

	return nil
}
