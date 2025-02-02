package typeconv

import (
	"errors"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
)

func ConvertCreateClientRequestToClient(createUserReq *dto.CreateClientRequest) (*models.Client, error) {
	if createUserReq == nil {
		return nil, errors.New("createUserReq is nil")
	}

	return &models.Client{
		Name:        createUserReq.Name,
		Description: createUserReq.Description,
		RedirectURI: createUserReq.RedirectURI,
		Status:      models.ClientStatus(createUserReq.Status),
		PublicKey:   createUserReq.PublicKey,
		KeyAlgo:     models.ClientKeyAlgo(createUserReq.KeyAlgo),
	}, nil
}
