package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type MongoConfig struct {
	MongoURI       string `yaml:"mongo_uri"`
	DatabaseName   string `yaml:"database_name"`
	CollectionName string `yaml:"collection_name"`
}

type Config struct {
	ServerHostPort string `env:"SERVER_HOST_PORT" envDefault:"localhost:8080"`
	Mongo          MongoConfig
}

func LoadConfig(path string) (*Config, error) {
	hostPort := flag.String("server-host-port", "", "Server host and port (e.g., localhost:8080)")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if *hostPort != "" {
		cfg.ServerHostPort = *hostPort
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var mongoConfig MongoConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&mongoConfig); err != nil {
		return nil, err
	}

	cfg.Mongo = mongoConfig

	return cfg, nil
}
