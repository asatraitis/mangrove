package config

import (
	"context"
	"errors"
	"sync"

	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/rs/zerolog"
)

//go:generate mockgen -destination=./mocks/mock_configs.go -package=mocks github.com/asatraitis/mangrove/internal/service/config Configs
type Configs interface {
	GetConfig(dal.ConfigKey) (string, error)
	SetAll(dal.Configs)
	UpdateConfig(dal.ConfigKey, string) error
	GetAll() dal.Configs
}
type BaseConfig struct {
	logger         zerolog.Logger
	mu             sync.RWMutex
	currentConfigs dal.Configs
}
type configs struct {
	ctx context.Context
	*BaseConfig
}

func NewConfig(ctx context.Context, logger zerolog.Logger) Configs {
	return &configs{
		ctx: ctx,
		BaseConfig: &BaseConfig{
			logger:         logger.With().Str("component", "Config").Logger(),
			currentConfigs: make(dal.Configs),
		},
	}
}

func (c *configs) GetAll() dal.Configs {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentConfigs
}

func (c *configs) SetAll(conf dal.Configs) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.currentConfigs = conf
}

func (c *configs) UpdateConfig(key dal.ConfigKey, value string) error {
	var funcName string = "UpdateConfig"
	if key == "" || value == "" {
		c.logger.Error().Str("func", funcName).Str("key", string(key)).Str("value", value).Msg("invalid config - empty key or value")
		return errors.New("invalid config - empty key or value")
	}
	c.mu.RLock()
	conf, ok := c.currentConfigs[key]
	c.mu.RUnlock()
	if !ok {
		c.logger.Error().Str("func", funcName).Msgf("config with key %s does not exist", key)
		return errors.New("config does not exist")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	conf.Value = &value
	c.currentConfigs[key] = conf

	return nil
}

func (c *configs) GetConfig(key dal.ConfigKey) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	const funcName string = "GetConfig"

	config, ok := c.currentConfigs[key]
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
