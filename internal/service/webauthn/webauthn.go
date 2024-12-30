package webauthn

import (
	"encoding/base64"
	"errors"

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
}
type webAuthN struct {
	logger zerolog.Logger
	wa     *webauthn.WebAuthn
	cache  utils.Cache[string, *webauthn.SessionData]
}

// TODO: Need to add cache expiration logic
func NewWebAuthN(logger zerolog.Logger) (WebAuthN, error) {
	logger = logger.With().Str("component", "WebAuthN").Logger()

	// TODO: remove hard coded confing and add env vars; look at more settings
	waConfig := &webauthn.Config{
		RPDisplayName: "Mangrove",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:3030"},
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
