package dal

import (
	"context"
	"errors"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -destination=./mocks/mock_clientsDAL.go -package=mocks github.com/asatraitis/mangrove/internal/dal ClientsDAL
type ClientsDAL interface {
	Create(pgx.Tx, *models.Client) error
	GetAllByUserID(uuid.UUID) ([]*models.Client, error)
}
type clientsDAL struct {
	ctx context.Context
	*BaseDAL
}

func NewClientsDAL(ctx context.Context, baseDAL *BaseDAL) ClientsDAL {
	cDal := clientsDAL{
		ctx:     ctx,
		BaseDAL: baseDAL,
	}
	cDal.logger = baseDAL.logger.With().Str("subcomponent", "ClientsDAL").Logger()
	return &cDal
}

func (c *clientsDAL) Create(tx pgx.Tx, client *models.Client) error {
	const funcName = "Create"
	const query = "INSERT INTO clients VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"

	if client == nil {
		c.logger.Error().Str("func", funcName).Msg("nil client")
		return errors.New("failed to create a client; nil")
	}

	var args = []interface{}{
		client.ID,
		client.UserID,
		client.Name,
		client.Description,
		client.Type,
		client.RedirectURI,
		client.PublicKey,
		client.KeyAlgo,
		client.KeyExpiresAt,
		client.Status,
	}

	if tx == nil {
		_, err := c.db.Exec(
			c.ctx,
			query,
			args...,
		)
		if err != nil {
			c.logger.Err(err).Str("func", funcName).Msg("failed to insert client")
		}
		return err
	}

	_, err := tx.Exec(
		c.ctx,
		query,
		args...,
	)
	if err != nil {
		c.logger.Err(err).Str("func", funcName).Msg("failed to insert client")
	}
	return err
}

func (c *clientsDAL) GetAllByUserID(userID uuid.UUID) ([]*models.Client, error) {
	const funcName = "GetAllByUserID"
	const query = "SELECT id, user_id, name, description, type, redirect_uri, status FROM clients WHERE user_id=$1"

	var clients []*models.Client
	err := pgxscan.Select(c.ctx, c.db, &clients, query, userID)
	if err != nil {
		c.logger.Err(err).Str("func", funcName).Msg("failed to get a clients")
		return nil, err
	}

	return clients, nil
}
