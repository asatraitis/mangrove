package dal

import (
	"context"
	"errors"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -destination=./mocks/mock_user_credential.go -package=mocks github.com/asatraitis/mangrove/internal/dal UserCredentialsDAL
type UserCredentialsDAL interface {
	Create(tx pgx.Tx, credential *models.UserCredential) error
	GetByUserID(uuid.UUID) ([]*models.UserCredential, error)
	UpdateSignCount(pgx.Tx, []byte, uint32) error
}
type userCredentialsDAL struct {
	ctx context.Context
	*BaseDAL
}

func NewUserCredentialsDAL(ctx context.Context, baseDAL *BaseDAL) UserCredentialsDAL {
	ucDAL := userCredentialsDAL{
		ctx:     ctx,
		BaseDAL: baseDAL,
	}
	ucDAL.logger = baseDAL.logger.With().Str("subcomponent", "UserCredentialsDAL").Logger()
	return &ucDAL
}

func (uc *userCredentialsDAL) Create(tx pgx.Tx, credential *models.UserCredential) error {
	const funcName string = "Create"
	const query string = "INSERT INTO user_credentials VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);"

	if credential == nil {
		uc.logger.Error().Str("func", funcName).Msg("nil credential")
		return errors.New("failed to create user credential; nil credential")
	}

	if tx == nil {
		_, err := uc.db.Exec(
			uc.ctx,
			query,
			credential.ID,
			credential.UserID,
			credential.PublicKey,
			credential.AttestationType,
			credential.Transport,
			credential.FlagUserPresent,
			credential.FlagVerified,
			credential.FlagBackupEligible,
			credential.FlagBackupState,
			credential.AuthAaguid,
			credential.AuthSignCount,
			credential.AuthCloneWarning,
			credential.AuthAttachment,
			credential.AttestationClientDataJson,
			credential.AttestationDataHash,
			credential.AttestationAuthenticatorData,
			credential.AttestationPublicKeyAlgorithm,
			credential.AttestationObject,
		)
		if err != nil {
			uc.logger.Err(err).Str("func", funcName).Msg("failed to insert user credential")
		}
		return err
	}

	_, err := tx.Exec(
		uc.ctx,
		query,
		credential.ID,
		credential.UserID,
		credential.PublicKey,
		credential.AttestationType,
		credential.Transport,
		credential.FlagUserPresent,
		credential.FlagVerified,
		credential.FlagBackupEligible,
		credential.FlagBackupState,
		credential.AuthAaguid,
		credential.AuthSignCount,
		credential.AuthCloneWarning,
		credential.AuthAttachment,
		credential.AttestationClientDataJson,
		credential.AttestationDataHash,
		credential.AttestationAuthenticatorData,
		credential.AttestationPublicKeyAlgorithm,
		credential.AttestationObject,
	)
	if err != nil {
		uc.logger.Err(err).Str("func", funcName).Msg("failed to insert user credential")
	}

	return err
}

func (uc *userCredentialsDAL) GetByUserID(userID uuid.UUID) ([]*models.UserCredential, error) {
	const funcName = "GetByUserID"

	const userCredentialsQuery = "SELECT id, user_id, public_key, attestation_type, transport, flag_user_present, flag_verified, flag_backup_eligible, flag_backup_state, auth_aaguid, auth_sign_count, auth_clone_warning, auth_attachment, attestation_client_data_json, attestation_data_hash, attestation_authenticator_data, attestation_public_key_algorithm, attestation_object FROM user_credentials WHERE user_id = $1"
	var credentials []*models.UserCredential
	err := pgxscan.Select(uc.ctx, uc.db, &credentials, userCredentialsQuery, userID)
	if err != nil {
		uc.logger.Err(err).Str("func", funcName).Msg("failed to get a user credentials")
		return nil, err
	}

	return credentials, nil
}

func (uc *userCredentialsDAL) UpdateSignCount(tx pgx.Tx, ID []byte, signCount uint32) error {
	const funcName = "UpdateSignCount"
	const query = "UPDATE user_credentials SET auth_sign_count=$1 WHERE id=$2"

	if tx == nil {
		_, err := uc.db.Exec(uc.ctx, query, signCount, ID)
		if err != nil {
			uc.logger.Err(err).Str("func", funcName).Msg("failed to update user credential sign count")
		}
		return err
	}
	_, err := tx.Exec(uc.ctx, query, signCount, ID)
	if err != nil {
		uc.logger.Err(err).Str("func", funcName).Msg("failed to update user credential sign count")
	}
	return err
}
