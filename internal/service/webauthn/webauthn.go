package webauthn

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/asatraitis/mangrove/internal/typeconv"
	"github.com/asatraitis/mangrove/internal/utils"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type WebAuthNUser struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
}

func (wau *WebAuthNUser) WebAuthnID() []byte {
	return []byte(wau.ID.String())
}
func (wau *WebAuthNUser) WebAuthnName() string {
	return wau.Name
}
func (wau *WebAuthNUser) WebAuthnDisplayName() string {
	return wau.DisplayName
}
func (wau *WebAuthNUser) WebAuthnCredentials() []webauthn.Credential {
	return wau.Credentials
}

type WebAuthN interface {
	BeginRegistration() (*protocol.CredentialCreation, error)
	FinishRegistration(string, *protocol.CredentialCreationResponse) (*webauthn.Credential, error)
	BeginLogin(*models.User, []webauthn.Credential) (*protocol.CredentialAssertion, string, error)
	FinishLogin(*dto.FinishLoginRequest, *models.User) (*webauthn.Credential, error)
	GetSession(key string) *webauthn.SessionData
}
type webAuthN struct {
	logger zerolog.Logger
	wa     *webauthn.WebAuthn
	cache  utils.Cache[string, *webauthn.SessionData]
}

// TODO: Need to add cache expiration logic
// technically, this will only be used for superadmin registration once - not critical to empty cache
// but likely need to add expiration logic to cache package
func NewWebAuthN(logger zerolog.Logger) (WebAuthN, error) {
	logger = logger.With().Str("component", "WebAuthN").Logger()

	// TODO: remove hard coded confing and add env vars; look at more settings
	waConfig := &webauthn.Config{
		RPDisplayName: "Mangrove",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:3030"},
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,             // Require the response from the client comes before the end of the timeout.
				Timeout:    time.Second * 60, // Standard timeout for login sessions.
				TimeoutUVD: time.Second * 60, // Timeout for login sessions which have user verification set to discouraged.
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,             // Require the response from the client comes before the end of the timeout.
				Timeout:    time.Second * 60, // Standard timeout for registration sessions.
				TimeoutUVD: time.Second * 60, // Timeout for login sessions which have user verification set to discouraged.
			},
		},
	}
	wa, err := webauthn.New(waConfig)
	if err != nil {
		logger.Err(err).Msg("Bad webauthn config")
		return nil, err
	}

	return &webAuthN{
		logger: logger,
		wa:     wa,
		cache:  utils.NewCache[string, *webauthn.SessionData](),
	}, nil
}

func (w *webAuthN) BeginRegistration() (*protocol.CredentialCreation, error) {
	id, err := uuid.NewV7()
	if err != nil {
		// TODO: add logging
		return nil, err
	}

	newUser := &WebAuthNUser{
		ID: id,
	}

	opts, session, err := w.wa.BeginRegistration(newUser)
	if err != nil {
		// TODO: add logging
		return nil, err
	}

	w.cache.SetValue(newUser.ID.String(), session)

	return opts, nil
}

func (w *webAuthN) FinishRegistration(userID string, credResp *protocol.CredentialCreationResponse) (*webauthn.Credential, error) {
	const funcName string = "FinishRegistration"

	if userID == "" {
		err := errors.New("missing userID")
		w.logger.Err(err).Str("func", funcName).Msg("failed to parse registration credential; missing userID")
		return nil, err
	}

	bUserID, err := base64.StdEncoding.DecodeString(userID)
	if err != nil {
		w.logger.Err(err).Str("func", funcName).Msg("failed to parse registration credential; failed to decode userID")
		return nil, err
	}

	userSession := w.cache.GetValue(string(bUserID))
	if userSession == nil {
		err = errors.New("failed to get user session")
		w.logger.Err(err).Str("func", funcName).Msg("failed to get user session from cache")
		return nil, err
	}

	parsedUserID, err := uuid.Parse(string(bUserID))
	if err != nil {
		w.logger.Err(err).Str("func", funcName).Msg("failed to parse user UUID")
		return nil, err
	}

	user := WebAuthNUser{ID: parsedUserID}

	parsedCred, err := credResp.Parse()
	if err != nil {
		w.logger.Err(err).Str("func", funcName).Msg("failed to parse registration credential")
		return nil, errors.New("failed to parse user credential")
	}

	credential, err := w.wa.CreateCredential(&user, *userSession, parsedCred)
	if err != nil {
		w.logger.Err(err).Str("func", funcName).Msg("failed to create webauthn credential")
		return nil, errors.New("failed to create user credential")
	}

	return credential, nil
}

func (w *webAuthN) BeginLogin(user *models.User, creds []webauthn.Credential) (*protocol.CredentialAssertion, string, error) {
	cacheKey, err := uuid.NewV7()
	if err != nil {
		// TODO: add logging
		return nil, "", err
	}
	placeholderUser := &WebAuthNUser{
		ID:          user.ID,
		Name:        user.Username,
		DisplayName: user.DisplayName,
		Credentials: creds,
	}
	opts, session, err := w.wa.BeginLogin(placeholderUser)
	if err != nil {
		return nil, "", err
	}
	w.cache.SetValue(cacheKey.String(), session)
	return opts, cacheKey.String(), nil
}

func (w *webAuthN) GetSession(key string) *webauthn.SessionData {
	return w.cache.GetValue(key)
}

func (w *webAuthN) FinishLogin(req *dto.FinishLoginRequest, user *models.User) (*webauthn.Credential, error) {
	const funcName = "FinishLogin"

	if req.SessionKey == "" {
		err := errors.New("could not find session in cache")
		w.logger.Err(err).Str("func", funcName).Msg("failed to get user session")
		return nil, err
	}

	userSession := w.cache.GetValue(req.SessionKey)
	if userSession == nil {
		err := errors.New("could not find session in cache")
		w.logger.Err(err).Str("func", funcName).Msg("failed to get user session")
		return nil, err
	}

	var credentials []webauthn.Credential
	for _, credential := range user.Credentials {
		waCredential, err := typeconv.ConvertUserCredentialToWebauthnCredential(credential)
		if err != nil {
			w.logger.Err(err).Str("func", funcName).Msg("failed to typeconv model.credential to webauthn.credential")
			return nil, err
		}
		credentials = append(credentials, *waCredential)
	}

	waUser := WebAuthNUser{
		ID:          user.ID,
		Name:        user.Username,
		DisplayName: user.DisplayName,
		Credentials: credentials,
	}

	parsedCred, err := req.Credential.Parse()
	if err != nil {
		w.logger.Err(err).Str("func", funcName).Msg("failed to parse credentials")
		return nil, err
	}

	credential, err := w.wa.ValidateLogin(&waUser, *userSession, parsedCred)
	if err != nil {
		w.logger.Err(err).Str("func", funcName).Msg("failed to validate login")
		return nil, err
	}

	if credential.Authenticator.CloneWarning {
		err := errors.New("cloned key error")
		w.logger.Err(err).Str("func", funcName).Msg("authenticator clone warning")
		return nil, err
	}

	return credential, nil
}
