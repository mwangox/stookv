package provider

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"stoo-kv/config"
	"strings"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
	cfg    config.Config
}

func NewRedisClient(ctx context.Context, config *config.Config) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     config.Providers.Redis.Host + ":" + config.Providers.Redis.Port,
			Password: config.Providers.Redis.Password,
			DB:       config.Providers.Redis.Database,
			PoolSize: config.Providers.Redis.ConnectionPoolSize,
		}),
		ctx: ctx,
	}
}
func (r *RedisClient) Set(key string, value any) error {
	return r.client.HSet(r.ctx, r.cfg.Providers.Redis.StoreName, key, value).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.HGet(r.ctx, r.cfg.Providers.Redis.StoreName, key).Result()
}

func (r *RedisClient) Delete(key string) error {
	return r.client.HDel(r.ctx, r.cfg.Providers.Redis.StoreName, key).Err()
}

func (r *RedisClient) GetAll() (map[string]string, error) {
	keyValues := make(map[string]string)
	result, err := r.client.HGetAll(r.ctx, r.cfg.Providers.Redis.StoreName).Result()
	if err != nil {
		return nil, err
	}
	for k, v := range result {
		if strings.Contains(k, "::") {
			keyValues[strings.Split(k, "::")[2]] = v
			continue
		}
		keyValues[k] = v
	}
	return keyValues, nil
}

func (r *RedisClient) GetByNameSpaceAndProfile(namespace, profile string) (map[string]string, error) {
	keyValues := make(map[string]string)
	result, err := r.client.HGetAll(r.ctx, r.cfg.Providers.Redis.StoreName).Result()
	if err != nil {
		return nil, err
	}
	for k, v := range result {
		if strings.HasPrefix(k, fmt.Sprintf("%s::%s", namespace, profile)) {
			keyValues[strings.Split(k, "::")[2]] = v
		}
	}
	return keyValues, nil
}
