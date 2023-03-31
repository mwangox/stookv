package cmd

import (
	"log"
	"stoo-kv/api"
	"stoo-kv/api/rpc"
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
	if err := rpc.RunGrpcServer(cfg, storage); err != nil {
		return err
	}
	log.Println("Initialize REST API routes...")
	return api.InitRoutes(storage, cfg)
}
