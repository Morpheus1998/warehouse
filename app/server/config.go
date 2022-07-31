package server

import (
	"errors"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const serverGracefulShutdownTime = 5 * time.Second

var ErrInvalidTypeForStore = errors.New("invalid type asserted for store")

type Configuration struct {
	HTTP struct {
		Port    int   `envconfig:"HTTP_PORT" default:"8080"`
		Timeout int64 `envconfig:"HTTP_TIMEOUT" default:"2000"`
	}
	PostgresConfiguration struct {
		Host                string `envconfig:"POSTGRES_HOST" default:"localhost"`
		Port                int    `envconfig:"POSTGRES_PORT" default:"5432"`
		DB                  string `envconfig:"POSTGRES_DATABASE" default:"warehouse"`
		CredentialsFileName string `envconfig:"POSTGRES_CREDENTIALS_FILE" default:"creds.json"`
	}
}

func GetConfigurationFromEnv() (Configuration, error) {
	var cfg Configuration
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Configuration{}, err
	}
	return cfg, nil
}
