package main

import (
	"os"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/migrations"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	variables := configs.NewConf(logger).GetEnvironmentVars()
	logger.Info().Msgf("MangroveEnv: %s", variables.MangroveEnv)

	migrator, err := migrations.NewMigrator(variables, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create migrator")
	}

	if err := migrator.Run(); err != nil {
		logger.Fatal().Err(err).Msg("could not run migrator")
	}
}
