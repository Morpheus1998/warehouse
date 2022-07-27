package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/warehouse/app/server"
)

func main() {
	log.Info().Msg("starting warehouse service")

	cfg, err := server.GetConfigurationFromEnv()
	if err != nil {
		log.Error().AnErr("error", err).Msg("reading configuration failed")
		os.Exit(1)
	}

	err = server.StartServer(cfg)
	if err != nil {
		log.Warn().AnErr("error", err).Msg("stopped warehouse service")
		os.Exit(1)
	}
	log.Info().Msg("stopped warehouse service")
}
