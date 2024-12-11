package bll

import (
	"context"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/rs/zerolog"
)

type BLL interface {
	Config(context.Context) ConfigBLL
}
type bll struct {
	logger zerolog.Logger
	vars   *configs.EnvVariables
	dal    dal.DAL
}

func NewBLL(logger zerolog.Logger, vars *configs.EnvVariables, dal dal.DAL) BLL {
	logger = logger.With().Str("component", "BLL").Logger()
	return &bll{
		logger: logger,
		vars:   vars,
		dal:    dal,
	}
}

func (b *bll) Config(ctx context.Context) ConfigBLL {
	return NewConfigBLL(ctx, b.logger, b.vars, b.dal)
}
