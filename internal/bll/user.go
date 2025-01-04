package bll

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
)

type UserBLL interface {
	CreateUserSession() (*protocol.CredentialCreation, error)
	RegisterSuperAdmin(*dto.FinishRegistrationRequest) error
}
type userBLL struct {
	ctx    context.Context
	hasher utils.Crypto
	*BaseBLL
}

func NewUserBLL(ctx context.Context, baseBLL *BaseBLL) UserBLL {
	uBll := &userBLL{
		ctx:     ctx,
		hasher:  utils.NewStandardCrypto([]byte(baseBLL.vars.MangroveSalt)),
		BaseBLL: baseBLL,
	}
	uBll.logger = baseBLL.logger.With().Str("subcomponent", "UserBLL").Logger()
	return uBll
}

func (u *userBLL) CreateUserSession() (*protocol.CredentialCreation, error) {
	const funcName string = "CreateUserSession"
	creds, err := u.webauthn.BeginRegistration()
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to generate user registration credentials")
		return nil, err
	}

	return creds, nil
}

func (u *userBLL) RegisterSuperAdmin(req *dto.FinishRegistrationRequest) error {
	const funcName string = "RegisterSuperAdmin"

	cred, err := u.webauthn.FinishRegistration(req.UserID, &req.Credential)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to validate user credential")
		return errors.New("failed to validate registration request")
	}

	bUserID, err := base64.StdEncoding.DecodeString(req.UserID)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to parse decode userID")
		return errors.New("failed to create superadmin")
	}

	userUUID, err := uuid.Parse(string(bUserID))
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to parse userID UUID")
		return errors.New("failed to create superadmin")
	}

	tx, err := u.dal.BeginTx(u.ctx)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to start DB transaction")
		return errors.New("failed to create superadmin")
	}
	defer func() {
		if err != nil {
			tx.Rollback(u.ctx)
		}
	}()

	err = u.dal.User(u.ctx).Create(tx, &models.User{
		ID:          userUUID,
		Username:    "Superadmin",
		DisplayName: "Superadmin",
		Email:       nil,
		Status:      models.USER_STATUS_ACTIVE,
		Role:        "superadmin",
	})
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to create superadmin in DB")
		return errors.New("failed to create superadmin")
	}

	// TODO: add typeconv
	ucred := &models.UserCredential{
		ID:                            cred.ID,
		UserID:                        userUUID,
		PublicKey:                     cred.PublicKey,
		AttestationType:               cred.AttestationType,
		Transport:                     cred.Transport,
		FlagUserPresent:               cred.Flags.UserPresent,
		FlagVerified:                  cred.Flags.UserVerified,
		FlagBackupEligible:            cred.Flags.BackupEligible,
		FlagBackupState:               cred.Flags.BackupState,
		AuthAaguid:                    cred.Authenticator.AAGUID,
		AuthSignCount:                 cred.Authenticator.SignCount,
		AuthCloneWarning:              cred.Authenticator.CloneWarning,
		AuthAttachment:                cred.Authenticator.Attachment,
		AttestationClientDataJson:     cred.Attestation.ClientDataJSON,
		AttestationDataHash:           cred.Attestation.ClientDataHash,
		AttestationAuthenticatorData:  cred.Attestation.AuthenticatorData,
		AttestationPublicKeyAlgorithm: cred.Attestation.PublicKeyAlgorithm,
		AttestationObject:             cred.Attestation.Object,
	}
	if len(cred.Transport) == 0 {
		ucred.Transport = []protocol.AuthenticatorTransport{}
	}

	err = u.dal.UserCredentials(u.ctx).Create(tx, ucred)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to create superadmin credential in DB")
		return errors.New("failed to create superadmin")
	}

	err = u.dal.Config(u.ctx).Set(dal.CONFIG_INSTANCE_READY, "true")
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to update instanceReady config in DB")
		return errors.New("failed to update instanceReady config")
	}
	err = u.appConfig.UpdateConfig(dal.CONFIG_INSTANCE_READY, "true")
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to update instanceReady in appConfig in DB")
		return errors.New("failed to update instanceReady config")
	}

	err = tx.Commit(u.ctx)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to commit DB transaction")
		return errors.New("failed to create superadmin")
	}
	return nil
}
