package typeconv

import (
	"testing"
	"time"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConvertClientsToUserClientsResponse_OK(t *testing.T) {
	now := time.Now()
	clients := []*models.Client{
		{
			ID:           uuid.MustParse("0bdd05ec-8008-4869-b6ec-6d812ce95508"),
			UserID:       uuid.MustParse("0bdd05ec-8008-4869-b6ec-6d812ce95507"),
			Name:         "test-client-name-0",
			Description:  "test-client-description-0",
			RedirectURI:  "http://localhost:3030",
			PublicKey:    []byte("test-public-key-0"),
			KeyExpiresAt: now,
			KeyAlgo:      models.CLIENT_KEY_ALGO_EDDSA,
			Status:       models.CLIENT_STATUS_ACTIVE,
		},
		{
			ID:           uuid.MustParse("0bdd05ec-8008-4869-b6ec-6d812ce95509"),
			UserID:       uuid.MustParse("0bdd05ec-8008-4869-b6ec-6d812ce95507"),
			Name:         "test-client-name-1",
			Description:  "test-client-description-1",
			RedirectURI:  "http://localhost:3031",
			PublicKey:    []byte("test-public-key-1"),
			KeyExpiresAt: now,
			KeyAlgo:      models.CLIENT_KEY_ALGO_EDDSA,
			Status:       models.CLIENT_STATUS_ACTIVE,
		},
	}
	userClients, err := ConvertClientsToUserClientsResponse(clients)
	assert.NoError(t, err)
	assert.Len(t, userClients, 2)

	assert.Equal(t, "0bdd05ec-8008-4869-b6ec-6d812ce95508", userClients[0].ID)
	assert.Equal(t, "0bdd05ec-8008-4869-b6ec-6d812ce95509", userClients[1].ID)
	assert.Equal(t, "0bdd05ec-8008-4869-b6ec-6d812ce95507", userClients[0].UserID)
	assert.Equal(t, "0bdd05ec-8008-4869-b6ec-6d812ce95507", userClients[1].UserID)
	assert.Equal(t, "0bdd05ec-8008-4869-b6ec-6d812ce95507", userClients[0].UserID)
	assert.Equal(t, "test-client-name-0", userClients[0].Name)
	assert.Equal(t, "test-client-name-1", userClients[1].Name)
	assert.Equal(t, "test-client-description-0", userClients[0].Description)
	assert.Equal(t, "test-client-description-1", userClients[1].Description)
	assert.Equal(t, "http://localhost:3030", userClients[0].RedirectURI)
	assert.Equal(t, "http://localhost:3031", userClients[1].RedirectURI)
	assert.Equal(t, dto.UserClientStatus("active"), userClients[0].Status)
	assert.Equal(t, dto.UserClientStatus("active"), userClients[1].Status)
}

func TestConvertClientsToUserClientsResponse_FAIL_NilUsers(t *testing.T) {
	_, err := ConvertClientsToUserClientsResponse(nil)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "client is nil")
}
