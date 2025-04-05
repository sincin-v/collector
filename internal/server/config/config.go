package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	Address string `env:"ADDRESS"`
}

type ServerConfig struct {
	Host string
}

func GetServerConfig() ServerConfig {
	var cfg ServerConfig
	var envCfg EnvConfig

	var argHostStr = flag.String("a", "localhost:8080", "Listen host and port")
	flag.Parse()

	err := env.Parse(&envCfg)
	if err != nil {
		log.Printf("Env is empty")
	}

	if envCfg.Address != "" {
		cfg.Host = envCfg.Address
	} else {
		cfg.Host = *argHostStr
	}

	return cfg
}
