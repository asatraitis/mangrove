package typeconv

import (
	"testing"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestConvertCreateClientRequestToClient_OK(t *testing.T) {
	clientReq := dto.CreateClientRequest{
		Name:        "test-name",
		Description: "test-desc",
		RedirectURI: "http://test.com",
		Status:      dto.UserClientStatus("active"),
		PublicKey:   []byte("pub_key"),
		KeyAlgo:     dto.UserClientKeyAlgo("EdDSA"),
	}
	client, err := ConvertCreateClientRequestToClient(&clientReq)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, clientReq.Name, client.Name)
	assert.Equal(t, clientReq.Description, client.Description)
	assert.Equal(t, clientReq.RedirectURI, client.RedirectURI)
	assert.Equal(t, models.ClientStatus(clientReq.Status), client.Status)
	assert.Equal(t, clientReq.PublicKey, client.PublicKey)
	assert.Equal(t, models.ClientKeyAlgo(clientReq.KeyAlgo), client.KeyAlgo)
}
