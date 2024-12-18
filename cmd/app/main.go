package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/bll"
	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/handler"
	"github.com/asatraitis/mangrove/internal/migrations"
	"github.com/asatraitis/mangrove/internal/service/config"
	"github.com/asatraitis/mangrove/internal/service/router"
	"github.com/asatraitis/mangrove/internal/service/webauthn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	variables := configs.NewConf(logger).GetEnvironmentVars()
	logger.Info().Msgf("MangroveEnv: %s", variables.MangroveEnv)

	if variables.MangroveEnv == configs.DEV {
		startDev(ctx, variables, logger)
	} else {
		// TODO: Add prod start
	}
}

func startDev(ctx context.Context, variables *configs.EnvVariables, logger zerolog.Logger) {
	logger = logger.Level(zerolog.DebugLevel).Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Run migrations
	migrator, err := migrations.NewMigrator(variables, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not create migrator")
		return
	}
	if err := migrator.Run(); err != nil {
		logger.Fatal().Err(err).Msg("could not run migrator")
	}

	// init db connection pool
	dbpool, err := initDbPool(ctx, variables)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not connect to the database")
		return
	}
	defer dbpool.Close()

	// init webauthn
	wauthn, err := webauthn.NewWebAuthN(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not init webAuthn")
		return
	}

	appConfig := config.NewConfig(ctx, logger)

	DAL := dal.NewDAL(logger, dbpool)
	BLL := bll.NewBLL(logger, variables, appConfig, wauthn, DAL)

	initCode, err := BLL.Config(ctx).InitRegistrationCode()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not init registration code")
		return
	}

	ro := router.NewRouter(
		logger,
		appConfig,
		handler.NewHandler(logger, BLL, variables, appConfig),
	)

	httpServer := &http.Server{
		Addr:    ":3030", // TODO: Add port config
		Handler: ro,
	}
	fmt.Printf("============================================ [REGISTRATION CODE: %s] ============================================\n", initCode)
	logger.Info().Msgf("Starting http server on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error().Err(err).Msg("Failed to start http server")
	}

}

// TODO: consolidate w/ getConnection() in migrator.go
func initDbPool(ctx context.Context, vars *configs.EnvVariables) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		vars.MangrovePostgresUser,
		vars.MangrovePostgresPassword,
		vars.MangrovePostgresAddress,
		vars.MangrovePostgresPort,
		vars.MangrovePostgresDBName,
	)
	return pgxpool.New(ctx, connStr)
}
