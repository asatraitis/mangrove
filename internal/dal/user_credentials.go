package dal

import (
	"context"
	"errors"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -destination=./mocks/mock_user_credential.go -package=mocks github.com/asatraitis/mangrove/internal/dal UserCredentialsDAL
type UserCredentialsDAL interface {
	Create(tx pgx.Tx, credential *models.UserCredential) error
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

	return nil
}
