package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/zalando/go-keyring"
)

const (
	Version                = "0.1.0"
	APIPrefix              = "/api/v1"
	CredentialService      = "com.lisitede.backstage.hub"
	VikunjaURLCredential   = "BACKSTAGE_VIKUNJA_URL"
	VikunjaTokenCredential = "BACKSTAGE_VIKUNJA_TOKEN"
)

type Config struct {
	Host         string
	Port         int
	VikunjaURL   string
	VikunjaToken string
}

func Default() Config {
	return Config{
		Host: "127.0.0.1",
		Port: 20080,
	}
}

func Load() (Config, error) {
	cfg := Default()

	var err error
	cfg.VikunjaURL, err = loadCredential(VikunjaURLCredential)
	if err != nil {
		return Config{}, fmt.Errorf("load Vikunja URL: %w", err)
	}
	cfg.VikunjaToken, err = loadCredential(VikunjaTokenCredential)
	if err != nil {
		return Config{}, fmt.Errorf("load Vikunja token: %w", err)
	}

	if v := os.Getenv("HUB_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("HUB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Port = port
		}
	}
	return cfg, nil
}

func loadCredential(account string) (string, error) {
	value, err := keyring.Get(CredentialService, account)
	if err == nil {
		if value == "" {
			return "", fmt.Errorf("credential %q is empty", account)
		}
		return value, nil
	}

	if !errors.Is(err, keyring.ErrNotFound) {
		return "", fmt.Errorf("read service %q account %q: %w", CredentialService, account, err)
	}

	return "", fmt.Errorf("credential %q not found in keychain", account)
}
