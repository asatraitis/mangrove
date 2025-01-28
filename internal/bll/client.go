package bll

import (
	"context"

	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/typeconv"
	"github.com/asatraitis/mangrove/internal/utils"
)

type ClientBLL interface {
	GetUserClients() (dto.UserClientsResponse, error)
}
type clientBLL struct {
	ctx context.Context
	*BaseBLL
}

func NewClientBLL(ctx context.Context, baseBLL *BaseBLL) ClientBLL {
	cBll := &clientBLL{
		ctx:     ctx,
		BaseBLL: baseBLL,
	}
	cBll.logger = baseBLL.logger.With().Str("subcomponent", "ClientBLL").Logger()
	return cBll
}

func (b *clientBLL) GetUserClients() (dto.UserClientsResponse, error) {
	const funcName = "GetUserClients"
	userID, err := utils.GetUserIdFromCtx(b.ctx)
	if err != nil {
		b.logger.Err(err).Str("func", funcName).Msg("failed to retrieve userID from context")
		return nil, err
	}

	clients, err := b.dal.Client(b.ctx).GetAllByUserID(userID)
	if err != nil {
		b.logger.Err(err).Str("func", funcName).Msg("failed to retrieve clients from db")
		return nil, err
	}

	userClients, err := typeconv.ConvertClientsToUserClientsResponse(clients)
	if err != nil {
		b.logger.Err(err).Str("func", funcName).Msg("failed to typeconv from clients to userclients dto")
		return nil, err
	}

	return userClients, nil
}
