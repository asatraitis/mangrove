package configs

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type MangroveEnvType string

const (
	DEV  MangroveEnvType = "dev"
	PROD MangroveEnvType = "prod"
)

type HttpConf struct {
	MangroveHost string
	MangrovePort string
}

type WebauthnConf struct {
	MangroveWebauthnRPDisplayName string
	MangroveWebauthnRPID          string
	MangroveWebauthnRPOrigins     []string
}

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
	// MangroveSal is salt used for hashing emails and init codes
	MangroveSalt string
	// http conf
	HttpConf

	WebauthnConf
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
		MangroveSalt:             c.getEnvByName("MANGROVE_SALT"),
		HttpConf: HttpConf{
			MangroveHost: c.getEnvByName("MANGROVE_HOST"),
			MangrovePort: c.getEnvByName("MANGROVE_PORT"),
		},
		WebauthnConf: WebauthnConf{
			MangroveWebauthnRPDisplayName: c.getEnvByName("MANGROVE_WEBAUTHN_RPDISPLAY_NAME"),
			MangroveWebauthnRPID:          c.getEnvByName("MANGROVE_WEBAUTHN_RPID"),
			MangroveWebauthnRPOrigins:     c.parseEnvListByName("MANGROVE_WEBAUTHN_RP_ORIGINS"),
		},
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
func (c *conf) parseEnvListByName(envName string) []string {
	envValue, ok := os.LookupEnv(envName)
	if ok {
		return strings.Split(envValue, ",")
	}
	c.logger.Warn().Msgf("Environment variable %s was not set", envName)
	return []string{}
}
