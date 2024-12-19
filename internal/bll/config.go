package bll

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/utils"
)

//go:generate mockgen -destination=./mocks/mock_config.go -package=mocks github.com/asatraitis/mangrove/internal/bll ConfigBLL
type ConfigBLL interface {
	GetAll() (dal.Configs, error)
	Set(dal.ConfigKey, string) error
	InitRegistrationCode() (string, error)
	ValidateRegistrationCode(string) error
}
type configBLL struct {
	ctx    context.Context
	hasher utils.Crypto
	*BaseBLL
}

func NewConfigBLL(ctx context.Context, baseBLL *BaseBLL) ConfigBLL {
	cbll := &configBLL{
		ctx:     ctx,
		hasher:  utils.NewCrypto(1, []byte(baseBLL.vars.MangroveSalt), 64*1024, 4, 32),
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

func (c *configBLL) Set(conf dal.ConfigKey, value string) error {
	const funcName string = "Set"

	if conf == "" || value == "" {
		err := errors.New("missing config")
		c.logger.Err(err).Str("func", funcName).Str(string(conf), value).Msg("missing conf")
		return err
	}

	err := c.dal.Config(c.ctx).Set(conf, value)
	if err != nil {
		c.logger.Err(err).Str("func", funcName).Str(string(conf), value).Msg("failed to retrieve config")
		return errors.New("failed to update config")
	}

	err = c.appConfig.UpdateConfig(conf, value)
	if err != nil {
		c.logger.Err(err).Str("func", funcName).Str(string(conf), value).Msg("failed to update config cache")
		return err
	}

	return nil
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
		initCode = utils.EncodeToString(6)
		hashedInitCode := c.hasher.GenerateBase64String([]byte(initCode))

		err = c.Set(dal.CONFIG_INIT_SA_CODE, hashedInitCode)
		if err != nil {
			c.logger.Err(err).Str("func", funcName).Msg("failed to update config cache")
			return "", err
		}
		err = c.Set(dal.CONFIG_INIT_ATTEMPTS, "0")
		if err != nil {
			c.logger.Err(err).Str("func", funcName).Msg("failed to update config cache")
			return "", err
		}
	}
	return initCode, nil
}

// TODO: Consolidate duplication
func (c *configBLL) ValidateRegistrationCode(registrationCode string) error {
	const funcName string = "ValidateRegistrationCode"

	if registrationCode == "" || len(registrationCode) != 6 {
		err := errors.New("no registration code provided")
		c.logger.Err(err).Str("func", funcName).Msg("missing code")
		return err
	}
	attempts, err := c.appConfig.GetConfig(dal.CONFIG_INIT_ATTEMPTS)
	if err != nil {
		c.logger.Err(err).Str("func", funcName).Msg("failed to retrieve config from cache")
		return err
	}
	intAttempts, err := strconv.Atoi(attempts)
	if err != nil {
		c.logger.Err(err).Str("func", funcName).Msg("failed to typeconv string to int (initAttempts)")
		return errors.New("faield to validate code")
	}
	if intAttempts >= 3 {
		err := errors.New("max registration attempts reached")
		c.logger.Err(err).Str("func", funcName).Msg("Superadmin registration attempts reached 3; rejecting until new code is generated on application start")
		return err
	}

	initCode, err := c.appConfig.GetConfig(dal.CONFIG_INIT_SA_CODE)
	if err != nil {
		c.logger.Err(err).Str("func", funcName).Msg("failed to retrieve config from cache")
		return err
	}

	decodedInitCode, err := c.hasher.DecodeBase64String(initCode)
	if err != nil {
		c.logger.Err(err).Str("func", funcName).Msg("failed to decode init registration code")
		return err
	}

	err = c.hasher.CompareValueToHash(registrationCode, decodedInitCode)
	if err != nil {
		c.logger.Err(err).Str("func", funcName).Msg("failed to validate code")
		attempts, err := c.appConfig.GetConfig(dal.CONFIG_INIT_ATTEMPTS)
		if err != nil {
			c.logger.Err(err).Str("func", funcName).Msg("failed to retrieve config from cache")
		}
		intAttempts, err := strconv.Atoi(attempts)
		if err != nil {
			c.logger.Err(err).Str("func", funcName).Msg("failed to typeconv string to int (initAttempts)")
		}
		err = c.Set(dal.CONFIG_INIT_ATTEMPTS, fmt.Sprintf("%d", intAttempts+1))
		if err != nil {
			c.logger.Fatal().Err(err).Str("func", funcName).Msg("failed to increment initAttempts!")
		}

		return errors.New("failed to validate code")
	}

	return err
}
