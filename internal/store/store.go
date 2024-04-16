package store

import (
	"context"
	"stoo-kv/config"
	"stoo-kv/internal/provider"
)

type Store interface {
	Set(key string, value any) error
	Get(key string) (string, error)
	Delete(key string) error
	//GetAll() (map[string]string, error)
	GetByNameSpaceAndProfile(namespace, profile string) (map[string]string, error)
}

func NewStorage(config *config.Config) (Store, error) {
	switch config.Application.StorageType {
	case "redis":
		return provider.NewRedisClient(context.Background(), config), nil
	case "mysql":
		return provider.NewMySql(config)
	case "postgres":
		return provider.NewPostgres(config)
	case "mongo":
		return provider.NewMongoClient(context.TODO(), config)
	case "etcd":
		return provider.NewEtcdClient(context.Background(), config)
	default:
		return provider.NewMemory(), nil
	}
}
