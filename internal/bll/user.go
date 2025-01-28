package bll

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/asatraitis/mangrove/internal/dal"
	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/typeconv"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

type UserBLL interface {
	CreateUserSession() (*protocol.CredentialCreation, error)
	RegisterSuperAdmin(*dto.FinishRegistrationRequest) error
	CreateToken(uuid.UUID) (*models.UserToken, error)
	GetUserByID(uuid.UUID) (*models.User, error)
	ValidateTokenAndGetUser(uuid.UUID) (*models.User, error)
	InitLogin(string) (protocol.PublicKeyCredentialRequestOptions, string, error)
	FinishLogin(*dto.FinishLoginRequest) (*dto.MeResponse, error)
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

	// TODO: consolidate app config updates in configDAL
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

func (u *userBLL) CreateToken(userID uuid.UUID) (*models.UserToken, error) {
	const funcName = "CreateToken"

	id, err := uuid.NewV7()
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Str("userID", userID.String()).Msg("failed to generate token ID")
		return nil, errors.New("failed to create user token")
	}

	token := &models.UserToken{
		ID:      id,
		UserID:  userID,
		Expires: time.Now().Add(time.Hour * 24),
	}

	// TODO: maybe use first 12bits of the v7 uuid to extract date time?
	err = u.dal.UserTokens(u.ctx).Create(nil, token)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Str("userID", userID.String()).Msg("failed to create token in db")
		return nil, errors.New("failed to create user token")
	}

	return token, nil
}

func (u *userBLL) GetUserByID(userID uuid.UUID) (*models.User, error) {
	const funcName = "GetUserByID"

	user, err := u.dal.User(u.ctx).GetByID(userID)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Str("userID", userID.String()).Msg("failed to get user from db")
		return nil, errors.New("failed to get user")
	}
	if user == nil {
		u.logger.Err(err).Str("func", funcName).Str("userID", userID.String()).Msg("failed to get user from db")
		return nil, errors.New("failed to get user")
	}

	return user, nil
}

func (u *userBLL) ValidateTokenAndGetUser(tokenID uuid.UUID) (*models.User, error) {
	const funcName = "ValidateTokenAndGetUser"

	userToken, err := u.dal.UserTokens(u.ctx).GetByIdWithUser(tokenID)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Str("tokenID", tokenID.String()).Msg("failed to get user token from db")
		return nil, errors.New("failed to validate token")
	}
	// make sure not nil
	if userToken == nil || userToken.User == nil {
		u.logger.Err(err).Str("func", funcName).Str("tokenID", tokenID.String()).Msg("returned nil token/user")
		return nil, errors.New("failed to retrieve user data")
	}

	// check if expired
	if !time.Now().Before(userToken.Expires) {
		u.logger.Err(err).Str("func", funcName).Str("tokenID", tokenID.String()).Msg("user token expired")
		return nil, errors.New("expired token")
	}

	return userToken.User, nil
}

func (u *userBLL) InitLogin(username string) (protocol.PublicKeyCredentialRequestOptions, string, error) {
	const funcName = "InitLogin"

	user, err := u.dal.User(u.ctx).GetByUsernameWithCredentials(username)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to init login credentials")
		return protocol.PublicKeyCredentialRequestOptions{}, "", err
	}

	var webauthnCreds []webauthn.Credential
	for _, userCred := range user.Credentials {
		webauthnCred, err := typeconv.ConvertUserCredentialToWebauthnCredential(userCred)
		if err != nil {
			u.logger.Err(err).Str("func", funcName).Msg("failed to typeconv credential")
			return protocol.PublicKeyCredentialRequestOptions{}, "", errors.New("failed to init login credentials")
		}
		webauthnCreds = append(webauthnCreds, *webauthnCred)
	}

	creds, sessionKey, err := u.webauthn.BeginLogin(user, webauthnCreds)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to init login credentials")
		return protocol.PublicKeyCredentialRequestOptions{}, "", err
	}

	return creds.Response, sessionKey, nil
}

func (u *userBLL) FinishLogin(login *dto.FinishLoginRequest) (*dto.MeResponse, error) {
	const funcName = "FinishLogin"

	if login == nil {
		err := errors.New("missing req struct")
		u.logger.Err(err).Str("func", funcName).Msg("failed to login")
		return nil, err
	}
	if login.SessionKey == "" {
		err := errors.New("missing session key")
		u.logger.Err(err).Str("func", funcName).Msg("failed to login")
		return nil, err
	}

	session := u.webauthn.GetSession(login.SessionKey)
	if session == nil {
		err := errors.New("missing session for the key")
		u.logger.Err(err).Str("func", funcName).Msg("failed to login")
		return nil, err
	}

	userID, err := uuid.Parse(string(session.UserID))
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to get user")
		return nil, err
	}

	user, err := u.dal.User(u.ctx).GetByIdWithCredentials(userID)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed get user from db")
		return nil, err
	}

	waCredential, err := u.webauthn.FinishLogin(login, user)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed webauthn validation")
		return nil, err
	}

	err = u.dal.UserCredentials(u.ctx).UpdateSignCount(nil, waCredential.ID, waCredential.Authenticator.SignCount)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed get update credential")
		return nil, err
	}

	meResponse, err := typeconv.ConvertUserToMeResponse(user)
	if err != nil {
		u.logger.Err(err).Str("func", funcName).Msg("failed to typeconv user to me")
		return nil, err
	}

	return meResponse, nil
}
