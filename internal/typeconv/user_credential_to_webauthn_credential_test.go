package typeconv

import (
	"testing"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConvertUserCredentialToWebauthnCredential(t *testing.T) {
	userID := uuid.New()
	userCred := &models.UserCredential{
		ID:                            []byte(userID.String()),
		UserID:                        userID,
		PublicKey:                     []byte("test-public-key"),
		AttestationType:               "basic",
		Transport:                     []protocol.AuthenticatorTransport{protocol.USB, protocol.NFC},
		FlagUserPresent:               true,
		FlagVerified:                  true,
		FlagBackupEligible:            true,
		FlagBackupState:               true,
		AuthAaguid:                    []byte("test-aaguid"),
		AuthSignCount:                 uint32(1),
		AuthCloneWarning:              true,
		AuthAttachment:                protocol.CrossPlatform,
		AttestationClientDataJson:     []byte("test-client-data-json"),
		AttestationDataHash:           []byte("test-data-hash"),
		AttestationAuthenticatorData:  []byte("test-authenticator-data"),
		AttestationPublicKeyAlgorithm: int64(1),
		AttestationObject:             []byte("test-attestation-object"),
	}
	waCred, err := ConvertUserCredentialToWebauthnCredential(userCred)
	assert.NoError(t, err)

	assert.Equal(t, userCred.ID, waCred.ID)
	assert.Equal(t, userCred.PublicKey, waCred.PublicKey)
	assert.Equal(t, userCred.AttestationType, waCred.AttestationType)

	assert.Equal(t, userCred.Transport, waCred.Transport)

	assert.Equal(t, userCred.FlagUserPresent, waCred.Flags.UserPresent)
	assert.Equal(t, userCred.FlagVerified, waCred.Flags.UserVerified)
	assert.Equal(t, userCred.FlagBackupEligible, waCred.Flags.BackupEligible)
	assert.Equal(t, userCred.FlagBackupState, waCred.Flags.BackupState)

	assert.Equal(t, userCred.AuthAaguid, waCred.Authenticator.AAGUID)
	assert.Equal(t, userCred.AuthSignCount, waCred.Authenticator.SignCount)
	assert.Equal(t, userCred.AuthCloneWarning, waCred.Authenticator.CloneWarning)
	assert.Equal(t, userCred.AuthAttachment, waCred.Authenticator.Attachment)

	assert.Equal(t, userCred.AttestationClientDataJson, waCred.Attestation.ClientDataJSON)
	assert.Equal(t, userCred.AttestationDataHash, waCred.Attestation.ClientDataHash)
	assert.Equal(t, userCred.AttestationAuthenticatorData, waCred.Attestation.AuthenticatorData)
	assert.Equal(t, userCred.AttestationPublicKeyAlgorithm, waCred.Attestation.PublicKeyAlgorithm)
	assert.Equal(t, userCred.AttestationObject, waCred.Attestation.Object)
}
