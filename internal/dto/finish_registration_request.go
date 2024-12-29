package dto

import (
	"github.com/go-webauthn/webauthn/protocol"
)

type FinishRegistrationRequest struct {
	protocol.CredentialCreationResponse
}
