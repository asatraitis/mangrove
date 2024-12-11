package bll

import (
	"context"
	"errors"
	"fmt"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/rs/zerolog"
)

type ConfigBLL interface {
	GetAll() (dal.Configs, error)
	InitRegistrationCode() error
}
type configBLL struct {
	ctx    context.Context
	logger zerolog.Logger
	vars   *configs.EnvVariables
	dal    dal.DAL

	configDAL dal.ConfigDAL
}

func NewConfigBLL(ctx context.Context, logger zerolog.Logger, vars *configs.EnvVariables, dal dal.DAL) ConfigBLL {
	logger = logger.With().Str("subcomponent", "ConfigBLL").Logger()
	return &configBLL{
		ctx:       ctx,
		logger:    logger,
		vars:      vars,
		dal:       dal,
		configDAL: dal.Config(ctx),
	}
}

func (c *configBLL) GetAll() (dal.Configs, error) {
	configs, err := c.configDAL.GetAll()
	if err != nil {
		c.logger.Err(err).Str("func", "GetAll")
		return nil, errors.New("failed to get configs")
	}
	return configs, nil
}

func (c *configBLL) InitRegistrationCode() error {
	const funcName string = "InitRegistrationCode"
	configs, err := c.GetAll()
	if err != nil {
		c.logger.Err(err).Str("func", funcName).Msg("failed to retrieve config")
		return err
	}
	instanceReady, ok := configs[dal.CONFIG_INSTANCE_READY]
	if !ok || instanceReady.Value == nil {
		c.logger.Err(err).Str("func", funcName).Msgf("config for %s not set", dal.CONFIG_INSTANCE_READY)
		return err
	}
	if *instanceReady.Value != "true" {
		hasher := utils.NewCrypto(1, []byte(c.vars.MangroveSalt), 64*1024, 4, 32)
		initCode := utils.EncodeToString(6)
		fmt.Printf("================================= [REGISTRATION CODE: %s] =================================\n", initCode)
		hashedInitCode := hasher.GenerateBase64String([]byte(initCode))

		err := c.configDAL.Set(dal.CONFIG_INIT_SA_CODE, hashedInitCode)
		if err != nil {
			c.logger.Err(err).Str("func", funcName).Msg("failed to set init code")
			return err
		}
	}
	return nil
}
