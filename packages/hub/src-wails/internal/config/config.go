package config

import (
	"os"
	"strconv"
)

const (
	Version   = "0.1.0"
	APIPrefix = "/api/v1"
)

type Config struct {
	Host string
	Port int
}

func Default() Config {
	return Config{
		Host: "127.0.0.1",
		Port: 20080,
	}
}

func Load() Config {
	cfg := Default()
	if v := os.Getenv("HUB_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("HUB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Port = port
		}
	}
	return cfg
}
