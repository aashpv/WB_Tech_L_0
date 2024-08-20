package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	StorageConn   string              `yaml:"storage_conn"`
	HttpServer    HttpServerConfig    `yaml:"http_server"`
	NatsStreaming NatsStreamingConfig `yaml:"nats_streaming"`
}

type HttpServerConfig struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type NatsStreamingConfig struct {
	Url       string `yaml:"url"`
	ClusterId string `yaml:"cluster_id"`
	ClientId  string `yaml:"client_id"`
	Subject   string `yaml:"subject"`
}

func MustLoad() *Config {
	configPath := "config/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file doesn't exist: ", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config: ", err)
	}

	return &cfg
}
