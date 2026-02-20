package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env             string
	Service         string
	HTTPAddr        string
	ShutdownTimeout time.Duration

	DBDSN        string
	KafkaBrokers string
	OTLPEndpoint string
	LogLevel     string
}

func Load() (Config, error) {
	cfg := Config{
		Service:         mustGetEnv("APP_SERVICE"),
		Env:             getEnv("APP_ENV", "dev"),
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 5*time.Second),

		DBDSN:        mustGetEnv("DB_DSN"),
		KafkaBrokers: mustGetEnv("KAFKA_BROKERS"),
		OTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
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

func mustGetEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		panic(fmt.Sprintf("%s must not be empty", key))
	}
	return val
}
