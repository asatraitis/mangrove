package dto

import "github.com/go-webauthn/webauthn/protocol"

type FinishLoginRequest struct {
	Credential protocol.CredentialAssertionResponse `json:"credential"`
	SessionKey string                               `json:"sessionKey"`
}
