package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

type AgentConfig struct {
	ServerHost     string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func GetAgentConfig() AgentConfig {
	var cfg AgentConfig
	var envCfg EnvConfig

	var argServerHost = flag.String("a", "localhost:8080", "Metric server host and port")
	var argReportInterval = flag.Int("r", 10, "Report interval")
	var argPollInterval = flag.Int("p", 2, "Poll interval")

	flag.Parse()

	err := env.Parse(&envCfg)
	if err != nil {
		log.Printf("Env is empty")
	}

	if envCfg.Address != "" {
		cfg.ServerHost = envCfg.Address
	} else {
		cfg.ServerHost = *argServerHost
	}
	if envCfg.ReportInterval != 0 {
		cfg.ReportInterval = envCfg.ReportInterval
	} else {
		cfg.ReportInterval = time.Duration(time.Duration(*argReportInterval) * time.Second)
	}

	if envCfg.PollInterval != 0 {
		cfg.PollInterval = envCfg.PollInterval
	} else {
		cfg.PollInterval = time.Duration(time.Duration(*argPollInterval) * time.Second)

	}

	return cfg
}
