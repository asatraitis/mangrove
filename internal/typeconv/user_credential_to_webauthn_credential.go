package typeconv

import (
	"errors"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/go-webauthn/webauthn/webauthn"
)

func ConvertUserCredentialToWebauthnCredential(ucred *models.UserCredential) (*webauthn.Credential, error) {
	if ucred == nil {
		return nil, errors.New("UserCredential is nil")
	}
	wacred := &webauthn.Credential{
		ID:              ucred.ID,
		PublicKey:       ucred.PublicKey,
		AttestationType: ucred.AttestationType,
	}

	wacred.Transport = ucred.Transport

	wacred.Flags.BackupEligible = ucred.FlagBackupEligible
	wacred.Flags.BackupState = ucred.FlagBackupState
	wacred.Flags.UserPresent = ucred.FlagUserPresent
	wacred.Flags.UserVerified = ucred.FlagVerified

	wacred.Authenticator.AAGUID = ucred.AuthAaguid
	wacred.Authenticator.Attachment = ucred.AuthAttachment
	wacred.Authenticator.CloneWarning = ucred.AuthCloneWarning
	wacred.Authenticator.SignCount = ucred.AuthSignCount

	wacred.Attestation.AuthenticatorData = ucred.AttestationAuthenticatorData
	wacred.Attestation.ClientDataHash = ucred.AttestationDataHash
	wacred.Attestation.ClientDataJSON = ucred.AttestationClientDataJson
	wacred.Attestation.Object = ucred.AttestationObject
	wacred.Attestation.PublicKeyAlgorithm = ucred.AttestationPublicKeyAlgorithm

	return wacred, nil
}
