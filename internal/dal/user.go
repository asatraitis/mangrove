package dal

import (
	"context"
	"errors"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -destination=./mocks/mock_user.go -package=mocks github.com/asatraitis/mangrove/internal/dal UserDAL
type UserDAL interface {
	Create(pgx.Tx, *models.User) error
}

type userDAL struct {
	ctx context.Context
	*BaseDAL
}

func NewUserDAL(ctx context.Context, baseDAL *BaseDAL) UserDAL {
	udal := &userDAL{
		ctx:     ctx,
		BaseDAL: baseDAL,
	}
	udal.logger = baseDAL.logger.With().Str("subcomponent", "UserDAL").Logger()
	return udal
}

func (ud *userDAL) Create(tx pgx.Tx, user *models.User) error {
	const funcName string = "Create"
	const query string = "INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6);"

	if user == nil {
		ud.logger.Error().Str("func", funcName).Msg("user is nil")
		return errors.New("failed to insert user; nil user")
	}

	if tx == nil {
		_, err := ud.db.Exec(
			ud.ctx,
			query,
			user.ID,
			user.Username,
			user.DisplayName,
			user.Email,
			user.Status,
			user.Role,
		)
		if err != nil {
			ud.logger.Err(err).Str("func", funcName).Msg("failed to insert user")
		}
		return err
	}

	_, err := tx.Exec(
		ud.ctx,
		query,
		user.ID,
		user.Username,
		user.DisplayName,
		user.Email,
		user.Status,
		user.Role,
	)
	if err != nil {
		ud.logger.Err(err).Str("func", funcName).Msg("failed to insert user")
	}
	return err
}
