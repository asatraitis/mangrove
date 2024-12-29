package bll

import (
	"context"

	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/go-webauthn/webauthn/protocol"
)

type UserBLL interface {
	CreateUserSession() (*protocol.CredentialCreation, error)
}
type userBLL struct {
	ctx    context.Context
	hasher utils.Crypto
	*BaseBLL
}

func NewUserBLL(ctx context.Context, baseBLL *BaseBLL) UserBLL {
	uBll := &userBLL{
		ctx:     ctx,
		hasher:  utils.NewStandardCrypto([]byte(baseBLL.vars.MangroveSalt)),
		BaseBLL: baseBLL,
	}
	uBll.logger = baseBLL.logger.With().Str("subcomponent", "UserBLL").Logger()
	return uBll
}

// TODO: refactor the function to seperate token from here into a new function in config
func (u *userBLL) CreateUserSession() (*protocol.CredentialCreation, error) {
	const funcName string = "CreateUserSession"
	creds, err := u.webauthn.BeginRegistration()
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to generate user registration credentials")
		return nil, err
	}

	return creds, nil
}
