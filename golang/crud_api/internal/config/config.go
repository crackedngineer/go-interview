package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string `yaml:"address"`
}

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-default:"dev" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true" env:"STORAGE_PATH"`
	HttpServer  HttpServer `yaml:"http_server" env-required:"true" env:"HTTP_SERVER_ADDRESS"`
}

func MustLoadConfig() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()
		configPath = *flags
	}
	if configPath == "" {
		log.Fatal("config path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist at path: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	return &cfg

}
