package cmd

import (
	"log"
	"stoo-kv/api"
	"stoo-kv/api/grpc"
	"stoo-kv/config"
	"stoo-kv/internal"
)

func Start() error {
	log.Println("Parse the configurations...")
	cfg, err := config.ParseConfig()
	if err != nil {
		return err
	}
	log.Println("Initialize key value pairs storage...")
	storage, err := internal.NewStorage(cfg)
	if err != nil {
		return err
	}

	log.Println("Start GRPC server asynchronously...")
	if err := grpc.RunGrpcServer(cfg, storage); err != nil {
		return err
	}
	log.Println("Initialize REST API routes...")
	return api.InitializeRoutes(storage, cfg)
}
