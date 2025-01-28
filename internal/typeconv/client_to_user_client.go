package typeconv

import (
	"errors"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
)

func ConvertClientsToUserClientsResponse(clients []*models.Client) ([]dto.UserClient, error) {
	if clients == nil {
		return nil, errors.New("client is nil")
	}

	var userClients dto.UserClientsResponse
	for _, client := range clients {
		userClients = append(userClients, dto.UserClient{
			ID:          client.ID.String(),
			UserID:      client.UserID.String(),
			Name:        client.Name,
			Description: client.Description,
			Type:        client.Type,
			RedirectURI: client.RedirectURI,
			Status:      dto.UserClientStatus(client.Status),
		})
	}

	return userClients, nil
}
