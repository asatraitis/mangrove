package typeconv

import (
	"testing"
	"time"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConvertClientToCreateClientResponse_OK(t *testing.T) {
	userID := "00000000-1111-2222-3333-000000000000"
	ID := "00000000-1111-2222-3333-000000000001"
	expires := time.Now()
	client := models.Client{
		ID:           uuid.MustParse(ID),
		UserID:       uuid.MustParse(userID),
		Name:         "test-client-name",
		Description:  "test-client-description",
		Type:         "app",
		RedirectURI:  "http://localhost:3030",
		PublicKey:    []byte("test-public-key"),
		KeyExpiresAt: expires,
		KeyAlgo:      models.CLIENT_KEY_ALGO_EDDSA,
		Status:       models.CLIENT_STATUS_ACTIVE,
	}

	createdClient, err := ConvertClientToCreateClientResponse(&client)

	assert.NoError(t, err)
	assert.NotNil(t, createdClient)

	assert.Equal(t, "00000000-1111-2222-3333-000000000001", createdClient.ID)
	assert.Equal(t, "00000000-1111-2222-3333-000000000000", createdClient.UserID)
	assert.Equal(t, "test-client-name", createdClient.Name)
	assert.Equal(t, "test-client-description", createdClient.Description)
	assert.Equal(t, "app", createdClient.Type)
	assert.Equal(t, "http://localhost:3030", createdClient.RedirectURI)
	assert.Equal(t, dto.UserClientStatus("active"), createdClient.Status)
}
