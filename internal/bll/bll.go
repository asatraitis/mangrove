package bll

import (
	"context"

	"github.com/asatraitis/mangrove/configs"
	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/service/config"
	"github.com/asatraitis/mangrove/internal/service/webauthn"
	"github.com/rs/zerolog"
)

//go:generate mockgen -destination=./mocks/mock_bll.go -package=mocks github.com/asatraitis/mangrove/internal/bll BLL
type BLL interface {
	Config(context.Context) ConfigBLL
	User(context.Context) UserBLL
	Client(context.Context) ClientBLL
}
type BaseBLL struct {
	logger    zerolog.Logger
	vars      *configs.EnvVariables
	appConfig config.Configs
	webauthn  webauthn.WebAuthN
	dal       dal.DAL
}
type bll struct {
	*BaseBLL
}

func NewBLL(logger zerolog.Logger, vars *configs.EnvVariables, appConfig config.Configs, webauthn webauthn.WebAuthN, dal dal.DAL) BLL {
	logger = logger.With().Str("component", "BLL").Logger()
	return &bll{
		BaseBLL: &BaseBLL{
			logger:    logger,
			vars:      vars,
			appConfig: appConfig,
			webauthn:  webauthn,
			dal:       dal,
		},
	}
}

func (b *bll) Config(ctx context.Context) ConfigBLL {
	return NewConfigBLL(ctx, b.BaseBLL)
}
func (b *bll) User(ctx context.Context) UserBLL {
	return NewUserBLL(ctx, b.BaseBLL)
}
func (b *bll) Client(ctx context.Context) ClientBLL {
	return NewClientBLL(ctx, b.BaseBLL)
}
