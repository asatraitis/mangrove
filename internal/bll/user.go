package bll

import (
	"context"

	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/go-webauthn/webauthn/protocol"
)

type UserBLL interface {
	CreateUserSession() (*protocol.CredentialCreation, string, error)
}
type userBLL struct {
	ctx    context.Context
	hasher utils.Crypto
	*BaseBLL
}

func NewUserBLL(ctx context.Context, baseBLL *BaseBLL) UserBLL {
	uBll := &userBLL{
		ctx:     ctx,
		hasher:  utils.NewCrypto(1, []byte(baseBLL.vars.MangroveSalt), 64*1024, 4, 32),
		BaseBLL: baseBLL,
	}
	uBll.logger = baseBLL.logger.With().Str("subcomponent", "UserBLL").Logger()
	return uBll
}

func (u *userBLL) CreateUserSession() (*protocol.CredentialCreation, string, error) {
	const funcName string = "CreateUserSession"
	creds, err := u.webauthn.BeginRegistration()
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to generate user registration credentials")
		return nil, "", err
	}

	token, sig, err := u.hasher.GenerateTokenHMAC()
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to generate user registration credentials")
		return nil, "", err
	}

	return creds, token + "." + sig, nil
}
