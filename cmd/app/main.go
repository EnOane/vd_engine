package main

import (
	"github.com/EnOane/vd_engine/internal/api"
	"github.com/EnOane/vd_engine/internal/config"
	"github.com/EnOane/vd_engine/internal/di"
	"github.com/EnOane/vd_engine/internal/infr/s3"
)

func main() {
	di.MakeDIContainer()

	config.MustLoad()
	s3.MustConnect()
	api.MustConnect()
}
