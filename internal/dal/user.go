package dal

import (
	"context"
	"errors"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -destination=./mocks/mock_user.go -package=mocks github.com/asatraitis/mangrove/internal/dal UserDAL
type UserDAL interface {
	Create(pgx.Tx, *models.User) error
	GetByID(uuid.UUID) (*models.User, error)
	GetByUsernameWithCredentials(string) (*models.User, error)
	GetByIdWithCredentials(uuid.UUID) (*models.User, error)
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

func (ud *userDAL) GetByID(ID uuid.UUID) (*models.User, error) {
	const funcName = "GetByID"

	row := ud.db.QueryRow(ud.ctx, "SELECT id, username, display_name, email, status, role FROM users WHERE id = $1", ID)
	user := &models.User{}

	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.Status, &user.Role)
	if err != nil {
		ud.logger.Err(err).Str("func", funcName).Msg("failed to get a user")
		return nil, err
	}

	return user, err
}

func (ud *userDAL) GetByUsernameWithCredentials(username string) (*models.User, error) {
	const funcName = "GetByUsernameWithCredentials"

	// TODO: would it be more efficient to JOIN these queries?
	const userQuery = "SELECT id, display_name, status, role FROM users WHERE username = $1"
	var user models.User
	row := ud.db.QueryRow(ud.ctx, userQuery, username)
	err := row.Scan(&user.ID, &user.DisplayName, &user.Status, &user.Role)
	if err != nil {
		ud.logger.Err(err).Str("func", funcName).Str("username", username).Msg("failed to get a user")
		return nil, err
	}
	user.Username = username

	credentials, err := ud.dal.UserCredentials(ud.ctx).GetByUserID(user.ID)
	if err != nil {
		ud.logger.Err(err).Str("func", funcName).Msg("failed to get a user credentials")
		return nil, err
	}

	user.Credentials = credentials
	return &user, err
}

func (ud *userDAL) GetByIdWithCredentials(ID uuid.UUID) (*models.User, error) {
	const funcName = "GetByIdWithCredentials"

	const userQuery = "SELECT id, username, display_name, status, role FROM users WHERE id = $1"
	var user models.User
	row := ud.db.QueryRow(ud.ctx, userQuery, ID)
	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Status, &user.Role)
	if err != nil {
		ud.logger.Err(err).Str("func", funcName).Str("ID", ID.String()).Msg("failed to get a user")
		return nil, err
	}

	credentials, err := ud.dal.UserCredentials(ud.ctx).GetByUserID(user.ID)
	if err != nil {
		ud.logger.Err(err).Str("func", funcName).Msg("failed to get a user credentials")
		return nil, err
	}

	user.Credentials = credentials
	return &user, err
}
