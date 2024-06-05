package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

const (
	devMode  = "development"
	testMode = "test"
	prodMode = "production"
)

type (
	Address string

	App struct {
		Mode string `env:"APP_MODE" validate:"required,oneof=development production test"`
	}

	Shortener struct {
		Address Address `env:"SERVER_ADDRESS" validate:"url"`
		Domain  string  `env:"BASE_URL"       validate:"url"`
	}

	Store struct {
		FileStoragePath string `env:"FILE_STORAGE_PATH"`
		DatabaseDSN     string `env:"DATABASE_DSN"`
	}

	Config struct {
		App       App
		Shortener Shortener
		Store     Store
	}
)

func NewConfig() (Config, error) {
	var config Config

	err := config.Shortener.Address.Set("localhost:8080")
	if err != nil {
		return Config{}, fmt.Errorf("failed to set default value: %w", err)
	}

	flag.Var(&config.Shortener.Address, "a", "server address host:port")
	flag.StringVar(&config.Shortener.Domain, "b", "http://localhost:8080", "domain url")
	flag.StringVar(&config.Store.FileStoragePath, "f", "file_store.json", "file storage path")
	flag.StringVar(&config.Store.DatabaseDSN, "d", "", "database connection string")

	flag.Parse()

	if err = env.Parse(&config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	if err = config.validate(); err != nil {
		return Config{}, fmt.Errorf("failed to validate config: %w", err)
	}

	return config, nil
}

func (c Config) validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(c)
	if err != nil {
		return fmt.Errorf("failed to validate config %v: %w", c, err)
	}

	return nil
}
