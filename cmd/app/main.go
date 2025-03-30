package main

import (
	"context"
	tgpb "github.com/EnOane/vd_engine/generated"
	"github.com/EnOane/vd_engine/internal/config"
	"github.com/EnOane/vd_engine/internal/ports/api"
	"github.com/EnOane/vd_engine/internal/service"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

func RunGRPSServer(server *grpc.Server, conf config.GrpcConfig) {
	lis, err := net.Listen("tcp", conf.ApiHost+":"+strconv.Itoa(conf.ApiPort))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}
	log.Info().Msgf("GRPC Server: Running on %s:%d", conf.ApiHost, conf.ApiPort)
	err = server.Serve(lis)

}

func main() {
	conf := config.NewConfig()
	ctx := context.Background()
	app := service.NewApplication(ctx, conf)
	serv := api.NewGrpcService(*app)
	server := grpc.NewServer()
	tgpb.RegisterTgServiceServer(server, serv)
	RunGRPSServer(server, conf.GrpcConfig)

}
