package main

import (
	"context"
	"github.com/EnOane/vd_engine/internal/config"
	"github.com/EnOane/vd_engine/internal/ports/api"
	"github.com/EnOane/vd_engine/internal/service"
)

func main() {
	service.MakeDIContainer()
	conf := config.NewConfig()

	ctx := context.Background()
	app := service.NewApplication(ctx)
	server := api.NewGrpcServer(conf)
}
