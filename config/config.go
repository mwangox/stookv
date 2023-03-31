package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Config struct {
	ServerPort            string `json:"server_port"`
	GrpcPort              string `json:"grpc_port"`
	StorageType           string `json:"storage_type"`
	EncryptKey            string `json:"encrypt_key"`
	EnableDecryptEndpoint bool   `json:"enable_decrypt_endpoint"`
	RdbmsDefaultTable     string `json:"rdbms_default_table"`
	EncryptPrefix         string `json:"encrypt_prefix"`
	Providers             struct {
		Redis struct {
			Host               string `json:"host"`
			Port               string `json:"port"`
			Password           string `json:"password"`
			Database           int    `json:"database"`
			ConnectionPoolSize int    `json:"connection_pool_size"`
			StoreName          string `json:"store_name"`
		} `json:"redis"`
		Mysql struct {
			Host         string `json:"host"`
			Port         string `json:"port"`
			Username     string `json:"username"`
			Password     string `json:"password"`
			DatabaseName string `json:"database_name"`
		} `json:"mysql"`
		Postgres struct {
			Host         string `json:"host"`
			Port         string `json:"port"`
			Username     string `json:"username"`
			Password     string `json:"password"`
			DatabaseName string `json:"database_name"`
			SslMode      string `json:"ssl_mode"`
			TimeZone     string `json:"timezone"`
		} `json:"postgres"`
		Mongo struct {
			MongoUri       string `json:"mongo_uri"`
			DatabaseName   string `json:"database_name"`
			CollectionName string `json:"collection_name"`
		} `json:"mongo"`
		Etcd struct {
			Endpoints   []string `json:"endpoints"`
			Username    string   `json:"username"`
			Password    string   `json:"password"`
			DialTimeout int      `json:"dial_timeout"`
		} `json:"etcd"`
	} `json:"providers"`
}

func NewConfig(configFile string) (*Config, error) {
	config := &Config{}
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	configs, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(configs, config); err != nil {
		return nil, err
	}
	return config, nil
}

func ParseConfig() (*Config, error) {
	var configFile string
	flag.StringVar(&configFile, "config.file", "./conf/stoo_kv.json", "Configuration file")
	flag.Parse()
	cfg, err := NewConfig(configFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
