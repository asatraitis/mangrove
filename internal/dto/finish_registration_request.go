package dto

import (
	"github.com/go-webauthn/webauthn/protocol"
)

type FinishRegistrationRequest struct {
	Credential protocol.CredentialCreationResponse `json:"credential"`
	UserID     string                              `json:"userId"`
}
