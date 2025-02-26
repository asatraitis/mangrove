// Migration generated by tools/migration_gen.js
package migrations

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type initial_20241203101104 struct {
	version int
}

func Newinitial_20241203101104() Migration {
	return &initial_20241203101104{
		version: 20241203101104,
	}
}

func (m *initial_20241203101104) Version() int {
	return m.version
}

func (m *initial_20241203101104) Up(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			display_name TEXT NOT NULL,
			email TEXT,
			status TEXT NOT NULL,
			role TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS user_credentials (
			id bytea NOT NULL,
			user_id uuid NOT NULL REFERENCES users(id),
			public_key bytea NOT NULL,
			attestation_type text NOT NULL DEFAULT 'none',
			transport text[] NOT NULL DEFAULT '{}',
			flag_user_present boolean NOT NULL DEFAULT FALSE,
			flag_verified boolean NOT NULL DEFAULT FALSE,
			flag_backup_eligible boolean NOT NULL DEFAULT FALSE,
			flag_backup_state boolean NOT NULL DEFAULT FALSE,
			auth_aaguid bytea,
			auth_sign_count integer NOT NULL DEFAULT 0,
			auth_clone_warning boolean NOT NULL DEFAULT FALSE,
			auth_attachment text,
			attestation_client_data_json bytea,
			attestation_data_hash bytea,
			attestation_authenticator_data bytea,
			attestation_public_key_algorithm bigint NOT NULL,
			attestation_object bytea
		);
		CREATE TABLE IF NOT EXISTS user_tokens (
			id UUID PRIMARY KEY,
			user_id uuid NOT NULL REFERENCES users(id),
			expires timestamp NOT NULL
		);
		CREATE TABLE IF NOT EXISTS clients (
			id UUID PRIMARY KEY,
			user_id uuid NOT NULL REFERENCES users(id),
			name TEXT NOT NULL,
			description TEXT,
			redirect_uri TEXT NOT NULL,
			public_key bytea NOT NULL,
			key_algo TEXT NOT NULl,
			key_expires_at timestamp NOT NULL,
			status TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS config (
			key TEXT PRIMARY KEY,
			label TEXT NOT NULL,
			value TEXT,
			type TEXT NOT NULL DEFAULT 'string',
			description TEXT
		);
		INSERT INTO config (key, label, value, type, description) VALUES 
			('initSACode', 'Super admin initialization code', '', 'string', 'Code to initialize superadmin user'),
			('instanceReady', 'Instance initialized', 'false', 'boolean', 'Flag to indicate if the instance has been initialized and ready'),
			('initAttempts', 'Super admin initialization attempts', '0', 'int', 'Application allows 3 attempts to provide a correct registration code and finish registration. Resets on application start.')
		ON CONFLICT DO NOTHING;
	`)
	return err
}
func (m *initial_20241203101104) Down(tx pgx.Tx) error {
	_, err := tx.Exec(context.Background(), `
        -- SQL migration revert
		DROP TABLE IF EXISTS clients
		DROP TABLE IF EXISTS user_tokens;
		DROP TABLE IF EXISTS user_credentials;
		DROP TABLE IF EXISTS users;
		DROP TABLE IF EXISTS config;
    `)
	return err
}
