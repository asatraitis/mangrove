package bll

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/typeconv"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/google/uuid"
)

type ClientBLL interface {
	GetUserClients() (dto.UserClientsResponse, error)
	Create(dto.CreateClientRequest) (*dto.CreateClientResponse, error)
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

func (b *clientBLL) Create(newUserClient dto.CreateClientRequest) (*dto.CreateClientResponse, error) {
	const funcName = "Create"
	if err := validateCreateReq(newUserClient); err != nil {
		b.logger.Err(err).Str("func", funcName).Msg("failed to validate CreateClientRequest")
		return nil, err
	}

	userID, err := utils.GetUserIdFromCtx(b.ctx)
	if err != nil {
		b.logger.Err(err).Str("func", funcName).Msg("failed to retrieve userID from context")
		return nil, err
	}

	client, err := typeconv.ConvertCreateClientRequestToClient(&newUserClient)
	if err != nil {
		b.logger.Err(err).Str("func", funcName).Msg("failed to typeconv dto client to model")
		return nil, err
	}

	ID, err := uuid.NewV7()
	if err != nil {
		b.logger.Err(err).Str("func", funcName).Msg("failed to create UUID")
		return nil, err
	}

	client.ID = ID
	client.UserID = userID

	// TODO: add config to allow max 3 months for key expiration
	client.KeyExpiresAt = time.Now().Add(time.Hour * 24 * 30)

	err = b.dal.Client(b.ctx).Create(nil, client)
	if err != nil {
		b.logger.Err(err).Str("func", funcName).Msg("failed to create clien in db")
		return nil, err
	}

	clientRes, err := typeconv.ConvertClientToCreateClientResponse(client)
	if err != nil {
		b.logger.Err(err).Str("func", funcName).Msg("failed to typeconv model client to dto")
		return nil, err
	}

	return clientRes, nil
}

func validateCreateReq(req dto.CreateClientRequest) error {
	var err error
	if strings.TrimSpace(req.Name) == "" {
		err = errors.Join(err, errors.New("missing name"))
	}
	if strings.TrimSpace(req.Type) == "" {
		err = errors.Join(err, errors.New("missing type"))
	}
	if strings.TrimSpace(req.RedirectURI) == "" {
		err = errors.Join(err, errors.New("missing redirectUri"))
	}
	if !slices.Contains([]dto.UserClientStatus{dto.CLIENT_STATUS_ACTIVE, dto.CLIENT_STATUS_PAUSED, dto.CLIENT_STATUS_SUSPENDED}, req.Status) {
		err = errors.Join(err, errors.New("missing or wrong status"))
	}
	if req.PublicKey == nil {
		err = errors.Join(err, errors.New("missing publicKey"))
	}
	if !slices.Contains([]dto.UserClientKeyAlgo{dto.CLIENT_KEY_ALGO_EDDSA}, req.KeyAlgo) {
		err = errors.Join(err, errors.New("missing or wrong keyAlgo"))
	}

	return err
}
