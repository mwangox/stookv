package provider

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"stoo-kv/config"
	"strings"
	"time"
)

type EtcdClient struct {
	client *clientv3.Client
	ctx    context.Context
}

func NewEtcdClient(ctx context.Context, config *config.Config) (*EtcdClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Providers.Etcd.Endpoints,
		DialTimeout: time.Duration(config.Providers.Etcd.DialTimeout) * time.Second,
		Username:    config.Providers.Etcd.Username,
		Password:    config.Providers.Etcd.Password,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdClient{client: client, ctx: ctx}, nil

}
func (e *EtcdClient) Set(key string, value any) error {
	_, err := e.client.Put(e.ctx, key, value.(string))
	return err
}

func (e *EtcdClient) Get(key string) (string, error) {
	resp, err := e.client.Get(e.ctx, key)
	if err != nil {
		return "", err
	}
	for _, v := range resp.Kvs {
		return string(v.Value), nil
	}
	return "", nil
}

func (e *EtcdClient) Delete(key string) error {
	_, err := e.client.Delete(e.ctx, key)
	return err
}

//func (e *EtcdClient) GetAll() (map[string]string, error) {
//	return e.findAll("")
//}

func (e *EtcdClient) GetByNameSpaceAndProfile(namespace, profile string) (map[string]string, error) {
	return e.findAll(fmt.Sprintf("%s::%s", namespace, profile))
}

func (e *EtcdClient) findAll(prefix string) (map[string]string, error) {
	keyValues := make(map[string]string)
	result, err := e.client.Get(e.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, v := range result.Kvs {
		key := string(v.Key)
		if strings.Contains(key, "::") {
			keyValues[strings.Split(key, "::")[2]] = string(v.Value)
			continue
		}
		keyValues[key] = string(v.Value)
	}
	return keyValues, nil
}
