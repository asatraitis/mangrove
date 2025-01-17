package dal

import (
	"context"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserTokensDAL interface {
	Create(pgx.Tx, *models.UserToken) error
	Get(uuid.UUID) (*models.UserToken, error)
}
type userTokensDAL struct {
	ctx context.Context
	*BaseDAL
}

func NewUserTokensDAL(ctx context.Context, baseDAL *BaseDAL) UserTokensDAL {
	utDAL := userTokensDAL{
		ctx:     ctx,
		BaseDAL: baseDAL,
	}
	utDAL.logger = baseDAL.logger.With().Str("subcomponent", "UserTokensDAL").Logger()
	return &utDAL
}

func (ut *userTokensDAL) Create(tx pgx.Tx, token *models.UserToken) error {
	const funcName string = "Create"
	const query = "INSERT INTO user_tokens VALUES ($1, $2, $3);"

	if tx == nil {
		_, err := ut.db.Exec(
			ut.ctx,
			query,
			token.ID,
			token.UserID,
			token.Expires,
		)
		if err != nil {
			ut.logger.Err(err).Str("func", funcName).Msg("failed to insert user token")
		}
		return err
	}

	_, err := tx.Exec(
		ut.ctx,
		query,
		token.ID,
		token.UserID,
		token.Expires,
	)
	if err != nil {
		ut.logger.Err(err).Str("func", funcName).Msg("failed to insert user token")
	}
	return err
}

func (ut *userTokensDAL) Get(ID uuid.UUID) (*models.UserToken, error) {
	const funcName = "Get"

	row := ut.db.QueryRow(ut.ctx, "SELECT id, user_id, expires FROM user_tokens WHERE id = $1", ID)
	token := &models.UserToken{}

	err := row.Scan(&token.ID, &token.UserID, &token.Expires)
	if err != nil {
		ut.logger.Err(err).Str("func", funcName).Msg("failed to get a token")
		return nil, err
	}

	return token, err
}
