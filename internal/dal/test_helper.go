package dal

import (
	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
)

func getUserCredential(userID uuid.UUID) *models.UserCredential {
	return &models.UserCredential{
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
}
