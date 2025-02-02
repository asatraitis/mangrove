package typeconv

import (
	"errors"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
)

func ConvertClientToCreateClientResponse(client *models.Client) (*dto.CreateClientResponse, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	return &dto.CreateClientResponse{
		ID:          client.ID.String(),
		UserID:      client.UserID.String(),
		Name:        client.Name,
		Description: client.Description,
		RedirectURI: client.RedirectURI,
		Status:      dto.UserClientStatus(client.Status),
	}, nil
}
