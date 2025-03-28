package main

import (
	"vd_engine/internal/api"
	"vd_engine/internal/config"
	"vd_engine/internal/di"
	"vd_engine/internal/infr/brokers"
	"vd_engine/internal/infr/s3"
)

func main() {
	config.MustLoad()

	di.MakeDIContainer()

	brokers.MustConnect()
	s3.MustConnect()
	api.MustServe()
}
