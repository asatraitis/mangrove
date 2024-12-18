package dal

import (
	"context"
)

//go:generate mockgen -destination=./mocks/mock_config.go -package=mocks github.com/asatraitis/mangrove/internal/dal ConfigDAL
type ConfigDAL interface {
	GetAll() (Configs, error)
	Set(ConfigKey, string) error
}
type configDAL struct {
	ctx context.Context
	*BaseDAL
}
type Config struct {
	Key         string  `json:"key"`
	Label       string  `json:"label"`
	Value       *string `json:"value"`
	Type        string  `json:"type"`
	Description *string `json:"description"`
}
type ConfigKey string

const (
	CONFIG_INSTANCE_READY ConfigKey = "instanceReady"
	CONFIG_INIT_SA_CODE   ConfigKey = "initSACode"
)

type Configs map[ConfigKey]Config

func NewConfigDAL(ctx context.Context, baseDAL *BaseDAL) ConfigDAL {
	cdal := &configDAL{
		ctx:     ctx,
		BaseDAL: baseDAL,
	}
	cdal.logger = baseDAL.logger.With().Str("subcomponent", "ConfigDAL").Logger()
	return cdal
}

func (c *configDAL) GetAll() (Configs, error) {
	const funcName = "GetAll"
	rows, err := c.db.Query(c.ctx, "SELECT key, label, value, type, description FROM config")
	if err != nil {
		c.logger.Err(err).Str("func", funcName)
		return nil, err
	}
	defer rows.Close()

	configs := make(Configs)
	for rows.Next() {
		var config Config
		err = rows.Scan(&config.Key, &config.Label, &config.Value, &config.Type, &config.Description)
		if err != nil {
			c.logger.Err(err).Str("func", funcName)
			return nil, err
		}
		configs[ConfigKey(config.Key)] = config
	}

	return configs, nil
}

func (c *configDAL) Set(key ConfigKey, value string) error {
	_, err := c.db.Exec(c.ctx, "UPDATE config SET value=$1 WHERE key=$2", value, key)
	if err != nil {
		c.logger.Err(err).Str("func", "Set")
		return err
	}
	return nil
}
