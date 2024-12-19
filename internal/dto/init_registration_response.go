package dto

import "github.com/go-webauthn/webauthn/protocol"

type InitRegistrationResponse struct {
	PublicKey protocol.PublicKeyCredentialCreationOptions `json:"publicKey"`
}
