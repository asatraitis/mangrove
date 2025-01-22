package dto

import "github.com/go-webauthn/webauthn/protocol"

type InitLoginResponse struct {
	PublicKey  protocol.PublicKeyCredentialRequestOptions `json:"publicKey"`
	SessionKey string                                     `json:"sessionKey"`
}
