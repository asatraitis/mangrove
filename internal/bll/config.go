package bll

import (
	"context"
	"errors"

	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/utils"
)

//go:generate mockgen -destination=./mocks/mock_config.go -package=mocks github.com/asatraitis/mangrove/internal/bll ConfigBLL
type ConfigBLL interface {
	GetAll() (dal.Configs, error)
	InitRegistrationCode() (string, error)
}
type configBLL struct {
	ctx context.Context
	*BaseBLL
}

func NewConfigBLL(ctx context.Context, baseBLL *BaseBLL) ConfigBLL {
	cbll := &configBLL{
		ctx:     ctx,
		BaseBLL: baseBLL,
	}
	cbll.logger = baseBLL.logger.With().Str("subcomponent", "ConfigBLL").Logger()
	return cbll
}

func (c *configBLL) GetAll() (dal.Configs, error) {
	configs, err := c.dal.Config(c.ctx).GetAll()
	if err != nil {
		c.logger.Err(err).Str("func", "GetAll")
		return nil, errors.New("failed to get configs")
	}
	c.appConfig.SetAll(configs)
	return configs, nil
}

func (c *configBLL) InitRegistrationCode() (string, error) {
	const funcName string = "InitRegistrationCode"
	var initCode string
	configs, err := c.GetAll()
	if err != nil || configs == nil {
		c.logger.Err(err).Str("func", funcName).Msg("failed to retrieve config")
		return "", errors.New("config instanceReady was not retrieved")
	}
	instanceReady, ok := configs[dal.CONFIG_INSTANCE_READY]
	if !ok || instanceReady.Value == nil {
		c.logger.Err(err).Str("func", funcName).Msgf("config for %s not set", dal.CONFIG_INSTANCE_READY)
		return "", errors.New("config instanceReady not set")
	}
	if *instanceReady.Value != "true" {
		hasher := utils.NewCrypto(1, []byte(c.vars.MangroveSalt), 64*1024, 4, 32)
		initCode = utils.EncodeToString(6)
		hashedInitCode := hasher.GenerateBase64String([]byte(initCode))

		configDAL := c.dal.Config(c.ctx)
		err := configDAL.Set(dal.CONFIG_INIT_SA_CODE, hashedInitCode)
		if err != nil {
			c.logger.Err(err).Str("func", funcName).Msg("failed to set init code")
			return "", err
		}
	}
	return initCode, nil
}
