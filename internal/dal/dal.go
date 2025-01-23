package dal

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

//go:generate mockgen -destination=./mocks/mock_dal.go -package=mocks github.com/asatraitis/mangrove/internal/dal DAL
type DAL interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)

	Config(ctx context.Context) ConfigDAL
	User(ctx context.Context) UserDAL
	UserCredentials(ctx context.Context) UserCredentialsDAL
	UserTokens(ctx context.Context) UserTokensDAL
}
type BaseDAL struct {
	logger zerolog.Logger
	db     *pgxpool.Pool
	dal    DAL
}
type dal struct {
	*BaseDAL
}

func NewDAL(logger zerolog.Logger, db *pgxpool.Pool) DAL {
	logger = logger.With().Str("component", "DAL").Logger()
	dal := &dal{
		BaseDAL: &BaseDAL{
			logger: logger,
			db:     db,
		},
	}

	// TODO: is this an anti-pattern? Want to use different DAL's withing other DAL's
	// would it be better to instantiate specific DALs with New...?
	dal.dal = dal
	return dal
}

func (d *dal) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return d.db.BeginTx(ctx, pgx.TxOptions{})
}

func (d *dal) Config(ctx context.Context) ConfigDAL {
	return NewConfigDAL(ctx, d.BaseDAL)
}
func (d *dal) User(ctx context.Context) UserDAL {
	return NewUserDAL(ctx, d.BaseDAL)
}
func (d *dal) UserCredentials(ctx context.Context) UserCredentialsDAL {
	return NewUserCredentialsDAL(ctx, d.BaseDAL)
}
func (d *dal) UserTokens(ctx context.Context) UserTokensDAL {
	return NewUserTokensDAL(ctx, d.BaseDAL)
}
