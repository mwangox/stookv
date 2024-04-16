package config

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

type Config struct {
	Application *ApplicationConfig
	Providers   *ProviderConfig
}

type ProviderConfig struct {
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
}

type ApplicationConfig struct {
	ServerLogLevel        string `json:"server_log_level"`
	ServerPort            string `json:"server_port"`
	ServerBindingHost     string `json:"server_binding_host"`
	GrpcPort              string `json:"grpc_port"`
	GrpcUseTls            bool   `json:"grpc_use_tls"`
	GrpcServerKey         string `json:"grpc_server_key"`
	GrpcServerCert        string `json:"grpc_server_cert"`
	StorageType           string `json:"storage_type"`
	EncryptKey            string `json:"encrypt_key"`
	EnableDecryptEndpoint bool   `json:"enable_decrypt_endpoint"`
	RdbmsDefaultTable     string `json:"rdbms_default_table"`
	EncryptPrefix         string `json:"encrypt_prefix"`
	ProviderPath          string `json:"provider_path"`
}

func NewApplicationConfig(configFile string) (*ApplicationConfig, error) {
	config := &ApplicationConfig{}
	configs, err := readFile(configFile)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(configs, config); err != nil {
		return nil, err
	}
	if config.ServerLogLevel == "" {
		config.ServerLogLevel = gin.ReleaseMode
	}

	return config, nil
}

func NewProviderConfig(providerConfigFile string) (*ProviderConfig, error) {
	config := &ProviderConfig{}
	configs, err := readFile(providerConfigFile)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(configs, config); err != nil {
		return nil, err
	}
	return config, nil
}

func readFile(configFile string) ([]byte, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ParseConfig() (*Config, error) {
	var configFile string
	flag.StringVar(&configFile, "config.file", "./conf/stoo_kv.json", "Configuration file")
	flag.Parse()
	applicationCfg, err := NewApplicationConfig(configFile)
	if err != nil {
		return nil, err
	}
	providerFile := applicationCfg.ProviderPath

	providerConfig, err := NewProviderConfig(providerFile)
	if err != nil {
		return nil, err
	}
	return &Config{
		Application: applicationCfg,
		Providers:   providerConfig,
	}, nil
}
