package main

import (
	"os"

	"github.com/asatraitis/mangrove/configs"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	variables := configs.NewConf(logger).GetEnvironmentVars()
	logger.Info().Msgf("MangroveEnv: %s", variables.MangroveEnv)
}
