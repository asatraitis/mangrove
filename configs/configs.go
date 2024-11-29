package configs

import (
	"os"

	"github.com/rs/zerolog"
)

type MangroveEnvType string

const (
	DEV  MangroveEnvType = "dev"
	PROD MangroveEnvType = "prod"
)

type EnvVariables struct {
	// MangroveEnv is the environment variable that specifies the environment in which the application is running
	// It can be either "dev" or "production"
	MangroveEnv MangroveEnvType
	// MangrovePostgresAddress is the address of the postgres database
	MangrovePostgresAddress string
	// MangrovePostgresPort is the port of the postgres database
	MangrovePostgresPort string
	// MangrovePostgresUser is the user of the postgres database
	MangrovePostgresUser string
	// MangrovePostgresPassword is the password of the postgres database
	MangrovePostgresPassword string
	// MangrovePostgresDBName is the name of the postgres database
	MangrovePostgresDBName string
}

type Conf interface {
	GetEnvironmentVars() *EnvVariables
}
type conf struct {
	logger zerolog.Logger
}

func NewConf(logger zerolog.Logger) Conf {
	logger = logger.With().Str("component", "Conf").Logger()
	return &conf{
		logger: logger,
	}
}

func (c *conf) GetEnvironmentVars() *EnvVariables {
	return &EnvVariables{
		MangroveEnv:              MangroveEnvType(c.getEnvByName("MANGROVE_ENV")),
		MangrovePostgresAddress:  c.getEnvByName("MANGROVE_POSTGRES_ADDRESS"),
		MangrovePostgresPort:     c.getEnvByName("MANGROVE_POSTGRES_PORT"),
		MangrovePostgresUser:     c.getEnvByName("MANGROVE_POSTGRES_USER"),
		MangrovePostgresPassword: c.getEnvByName("MANGROVE_POSTGRES_PASSWORD"),
		MangrovePostgresDBName:   c.getEnvByName("MANGROVE_POSTGRES_DB_NAME"),
	}
}

func (c *conf) getEnvByName(envName string) string {
	envValue, ok := os.LookupEnv(envName)
	if ok {
		return envValue
	}
	c.logger.Warn().Msgf("Environment variable %s was not set", envName)
	return ""
}
