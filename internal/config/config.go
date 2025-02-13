package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                string `yaml:"env" env-default:"local"`
	MongoURL           string `yaml:"mongo_url" env-required:"true"`
	DbName             string `yaml:"db_name" env-required:"true"`
	HTTPServer         `yaml:"http_server"`
	ProxyServer        `yaml:"proxy_server"`
	SpotifyCredentials `yaml:"spotify_credentials"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" enc-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" enc-default:"60s"`
}

type ProxyServer struct {
	Address  string        `yaml:"address" env-required:"true"`
	Username string        `yaml:"username" env-required:"true"`
	Password string        `yaml:"password" env-required:"true"`
	Timeout  time.Duration `yaml:"timeout" enc-default:"10s"`
}

type SpotifyCredentials struct {
	ClientID     string `yaml:"client_id" env-required:"true"`
	ClientSecret string `yaml:"client_secret" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
