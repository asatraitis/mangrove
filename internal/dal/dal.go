package dal

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

//go:generate mockgen -destination=./mocks/mock_dal.go -package=mocks github.com/asatraitis/mangrove/internal/dal DAL
type DAL interface {
	Config(ctx context.Context) ConfigDAL
}
type BaseDAL struct {
	logger zerolog.Logger
	db     *pgxpool.Pool
}
type dal struct {
	*BaseDAL
}

func NewDAL(logger zerolog.Logger, db *pgxpool.Pool) DAL {
	logger = logger.With().Str("component", "DAL").Logger()
	return &dal{
		BaseDAL: &BaseDAL{
			logger: logger,
			db:     db,
		},
	}
}

func (d *dal) Config(ctx context.Context) ConfigDAL {
	return NewConfigDAL(ctx, d.BaseDAL)
}
