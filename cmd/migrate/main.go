package main

import (
	"os"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/migrations"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel).Output(zerolog.ConsoleWriter{Out: os.Stderr})
	variables := configs.NewConf(logger).GetEnvironmentVars()

	migrator, err := migrations.NewMigrator(variables, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create migrator")
		return
	}
	if err := migrator.Run(); err != nil {
		logger.Fatal().Err(err).Msg("could not run migrator")
		return
	}

	logger.Info().Msg("Done!")
}
