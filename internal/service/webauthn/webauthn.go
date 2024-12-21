package webauthn

import (
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
		RPOrigins:     []string{"http://localhost"},
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
