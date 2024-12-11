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
type dal struct {
	logger zerolog.Logger
	db     *pgxpool.Pool
}

func NewDAL(logger zerolog.Logger, db *pgxpool.Pool) DAL {
	logger = logger.With().Str("component", "DAL").Logger()
	return &dal{
		logger: logger,
		db:     db,
	}
}

func (d *dal) Config(ctx context.Context) ConfigDAL {
	return NewConfigDAL(ctx, d.logger, d.db)
}
