package config

import (
	"context"
	"errors"
	"sync"

	"github.com/asatraitis/mangrove/internal/bll"
	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/rs/zerolog"
)

//go:generate mockgen -destination=./mocks/mock_configs.go -package=mocks github.com/asatraitis/mangrove/internal/service/config Configs
type Configs interface {
	GetConfig(dal.ConfigKey) (string, error)
	GetAll() (dal.Configs, error)
	Reload() error
}
type configs struct {
	ctx    context.Context
	logger zerolog.Logger
	bll    bll.BLL
	mu     sync.RWMutex

	currentConfigs dal.Configs
}

func NewConfig(ctx context.Context, logger zerolog.Logger, bll bll.BLL) Configs {
	logger = logger.With().Str("component", "Config").Logger()
	return &configs{
		ctx:    ctx,
		logger: logger,
		bll:    bll,
	}
}

func (c *configs) GetAll() (dal.Configs, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.currentConfigs == nil {
		all, err := c.bll.Config(c.ctx).GetAll()
		if err != nil {
			c.logger.Err(err).Str("func", "GetAll")
			return nil, err
		}
		c.currentConfigs = all
	}
	return c.currentConfigs, nil
}

func (c *configs) GetConfig(key dal.ConfigKey) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	const funcName string = "GetConfig"
	all, err := c.GetAll()
	if err != nil {
		c.logger.Err(err).Str("func", funcName)
		return "", err
	}
	if all == nil {
		err = errors.New("failed to retrieve config")
		c.logger.Err(err).Str("func", funcName)
		return "", err
	}
	config, ok := all[key]
	if !ok {
		c.logger.Info().Str("func", funcName).Msgf("config with key %s was not found", key)
		return "", errors.New("config not found")
	}
	if config.Value == nil {
		c.logger.Info().Str("func", funcName).Msgf("config with key %s is not set (nil)", key)
		return "", errors.New("config not set (nil)")
	}
	return *config.Value, nil
}

func (c *configs) Reload() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	all, err := c.bll.Config(c.ctx).GetAll()
	if err != nil {
		c.logger.Err(err).Str("func", "Reload")
		return err
	}
	c.currentConfigs = all
	return nil
}
