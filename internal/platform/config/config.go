package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env             string
	Service         string
	HTTPAddr        string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Env:             getEnv("APP_ENV", "local"),
		Service:         getEnv("APP_SERVICE", "orders"),
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 5*time.Second),
	}

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.Service == "" {
		return errors.New("APP_SERVICE must not be empty")
	}

	if cfg.HTTPAddr == "" {
		return errors.New("HTTP_ADDR must not be empty")
	}

	if len(cfg.HTTPAddr) > 1 && cfg.HTTPAddr[0] == ':' {
		port, err := strconv.Atoi(cfg.HTTPAddr[1:])

		if err != nil || port <= 0 || port > 65535 {
			return errors.New("invalid HTTP_ADDR port")
		}
	}

	if cfg.ShutdownTimeout <= 0 {
		return errors.New("SHUTDOWN_TIMEOUT must be > 0")
	}

	return nil
}

func getEnv(key, def string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}
	return v
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
