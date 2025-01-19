package dal

import (
	"context"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -destination=./mocks/mock_user_tokens.go -package=mocks github.com/asatraitis/mangrove/internal/dal UserTokensDAL
type UserTokensDAL interface {
	Create(pgx.Tx, *models.UserToken) error
	GetByID(uuid.UUID) (*models.UserToken, error)
	GetByIdWithUser(uuid.UUID) (*models.UserToken, error)
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

func (ut *userTokensDAL) GetByID(ID uuid.UUID) (*models.UserToken, error) {
	const funcName = "GetByID"

	row := ut.db.QueryRow(ut.ctx, "SELECT id, user_id, expires FROM user_tokens WHERE id = $1", ID)
	token := &models.UserToken{}

	err := row.Scan(&token.ID, &token.UserID, &token.Expires)
	if err != nil {
		ut.logger.Err(err).Str("func", funcName).Msg("failed to get a token")
		return nil, err
	}

	return token, err
}

func (ut *userTokensDAL) GetByIdWithUser(ID uuid.UUID) (*models.UserToken, error) {
	const funcName = "GetByIdWithUser"

	row := ut.db.QueryRow(ut.ctx, "SELECT ut.id, ut.user_id, ut.expires, u.id, u.username, u.display_name, u.status, u.role  FROM user_tokens ut JOIN users u ON ut.user_id = u.id WHERE ut.id = $1", ID)

	user := &models.User{}
	token := &models.UserToken{}

	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.Expires,
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Status,
		&user.Role,
	)
	if err != nil {
		ut.logger.Err(err).Str("func", funcName).Msg("failed to get a token")
		return nil, err
	}
	token.User = user

	return token, err
}
