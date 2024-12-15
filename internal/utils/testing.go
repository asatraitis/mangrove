package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/asatraitis/mangrove/configs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type MockResponseWriter struct {
	Body io.ReadCloser
	Code int
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{}
}

func (w *MockResponseWriter) Header() http.Header {
	return http.Header{}
}
func (w *MockResponseWriter) Write(b []byte) (int, error) {
	return 0, nil
}
func (w *MockResponseWriter) WriteHeader(code int) {
	w.Code = code
}

// TODO: consolidate w/ getConnection() in migrator.go and initDbPool in main
func InitDbPool(ctx context.Context) (*pgxpool.Pool, error) {
	conf := configs.NewConf(zerolog.Logger{})
	vars := conf.GetEnvironmentVars()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		vars.MangrovePostgresUser,
		vars.MangrovePostgresPassword,
		vars.MangrovePostgresAddress,
		vars.MangrovePostgresPort,
		vars.MangrovePostgresDBName,
	)
	return pgxpool.New(ctx, connStr)
}
