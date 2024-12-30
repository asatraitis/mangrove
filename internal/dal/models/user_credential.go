package models

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
)

type UserCredential struct {
	ID                            []byte                            `json:"id"`
	UserID                        uuid.UUID                         `json:"userId"`
	PublicKey                     []byte                            `json:"publicKey"`
	AttestationType               string                            `json:"attestationType"`
	Transport                     []protocol.AuthenticatorTransport `json:"transport"`
	FlagUserPresent               bool                              `json:"flagUserPresent"`
	FlagVerified                  bool                              `json:"flagVerified"`
	FlagBackupEligible            bool                              `json:"flagBackupEligible"`
	FlagBackupState               bool                              `json:"flagBackupState"`
	AuthAaguid                    []byte                            `json:"authAaguid"`
	AuthSignCount                 uint32                            `json:"authSignCount"`
	AuthCloneWarning              bool                              `json:"authCloneWarning"`
	AuthAttachment                protocol.AuthenticatorAttachment  `json:"authAttachment"`
	AttestationClientDataJson     []byte                            `json:"attestationClientDataJson"`
	AttestationDataHash           []byte                            `json:"attestationDataHash"`
	AttestationAuthenticatorData  []byte                            `json:"attestationAuthenticatorData"`
	AttestationPublicKeyAlgorithm int64                             `json:"attestationPublicKeyAlgorithm"`
	AttestationObject             []byte                            `json:"attestationObject"`
}
