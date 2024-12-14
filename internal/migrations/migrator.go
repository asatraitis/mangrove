package migrations

import (
	"context"
	"fmt"
	"sort"

	"github.com/asatraitis/mangrove/configs"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

const migrationsTable = "schema_history"

type Migration interface {
	Up(pgx.Tx) error
	Down(pgx.Tx) error
	Version() int
}
type Migrator struct {
	logger zerolog.Logger

	db *pgx.Conn

	currentVersion int
	migrations     []Migration
}

func NewMigrator(vars *configs.EnvVariables, logger zerolog.Logger) (*Migrator, error) {
	l := logger.With().Str("component", "Migrator").Logger()
	ctx := context.Background()

	// Get connection to the database
	db, err := getConnection(ctx, vars)
	if err != nil {
		l.Error().Err(err).Msg("could not connect to the database")
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}
	m := &Migrator{
		db:     db,
		logger: l,
	}

	// Check if migrations table exists
	if err := m.initMigrationTable(); err != nil {
		l.Error().Err(err).Msg("could not initialize migration table")
		return nil, fmt.Errorf("could not initialize migration table: %w", err)
	}

	// Get current version
	currentVersion, err := m.getCurrentVersion()
	if err != nil {
		l.Error().Err(err).Msg("could not get current version")
		return nil, fmt.Errorf("could not get current version: %w", err)
	}
	m.currentVersion = currentVersion

	// Set migrations
	m.setMigrations()

	return m, nil
}

func (m *Migrator) Run() error {
	defer m.db.Close(context.Background())
	if len(m.migrations) == 0 {
		m.logger.Info().Msg("no migrations to apply")
		return nil
	}
	m.logger.Info().Msg("running migrations")
	// find index of the current version
	var currentVersionIdx int
	for i, migration := range m.migrations {
		if migration.Version() == m.currentVersion {
			currentVersionIdx = i
			break
		}
	}

	if m.currentVersion != 0 && currentVersionIdx == len(m.migrations)-1 {
		m.logger.Info().Msg("Migrations up to date")
		return nil
	}

	var startIdx int
	if currentVersionIdx != 0 {
		startIdx += 1
	}

	// iterate over migrations from versionIdx
	for _, migration := range m.migrations[startIdx:] {
		m.logger.Info().Msgf("applying migration %d", migration.Version())
		// start transaction
		tx, err := m.db.Begin(context.Background())
		if err != nil {
			m.logger.Error().Err(err).Msg("could not start transaction")
			return fmt.Errorf("could not start transaction: %w", err)
		}

		// apply migration
		if err := migration.Up(tx); err != nil {
			m.logger.Error().Err(err).Msg("could not apply migration")
			_ = tx.Rollback(context.Background())
			return fmt.Errorf("could not apply migration: %w", err)
		}

		// commit transaction
		if err := tx.Commit(context.Background()); err != nil {
			return fmt.Errorf("could not commit transaction: %w", err)
		}

		// update current version
		err = m.setCurrentVersion(migration.Version())
		if err != nil {
			m.logger.Error().Err(err).Msg("could not update current version")
			return fmt.Errorf("could not update current version: %w", err)
		}
		m.currentVersion = migration.Version()
		m.logger.Info().Msgf("migration %d applied", migration.Version())
	}

	return nil
}

// TODO: Consolidate w/ initDbPool() in main.go
func getConnection(ctx context.Context, vars *configs.EnvVariables) (*pgx.Conn, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", vars.MangrovePostgresUser, vars.MangrovePostgresPassword, vars.MangrovePostgresAddress, vars.MangrovePostgresPort, vars.MangrovePostgresDBName)
	return pgx.Connect(ctx, connectionString)
}

func (m *Migrator) initMigrationTable() error {
	_, err := m.db.Exec(context.Background(), fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (version BIGINT NOT NULL PRIMARY KEY)", migrationsTable))
	return err
}

// get current version by getting largest integer from the migrations table
func (m *Migrator) getCurrentVersion() (int, error) {
	var version int
	err := m.db.QueryRow(context.Background(), fmt.Sprintf("SELECT COALESCE(MAX(version), 0) FROM %s", migrationsTable)).Scan(&version)
	return version, err
}

// sort migrations by version
func (m *Migrator) sortMigrations() {
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version() < m.migrations[j].Version()
	})
}

// inserts version into the migrations table
func (m *Migrator) setCurrentVersion(version int) error {
	_, err := m.db.Exec(context.Background(), fmt.Sprintf("INSERT INTO %s (version) VALUES ($1)", migrationsTable), version)
	return err
}
