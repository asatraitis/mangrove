package typeconv

import (
	"errors"
	"slices"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
)

func ConvertUserToMeResponse(user *models.User) (*dto.MeResponse, error) {
	meRole := dto.UserRole(user.Role)
	if !slices.Contains([]dto.UserRole{
		dto.USER_ROLE_ADMIN,
		dto.USER_ROLE_SUPERUSER,
		dto.USER_ROLE_USER,
	}, meRole) {
		return nil, errors.New("user role is not available in the list of roles. Role: " + string(meRole))
	}

	meStatus := dto.UserStatus(user.Status)
	if !slices.Contains([]dto.UserStatus{
		dto.USER_STATUS_ACTIVE,
		dto.USER_STATUS_INACTIVE,
		dto.USER_STATUS_PENDING,
	}, meStatus) {
		return nil, errors.New("user status is not available in the list of statuses. Status: " + string(meStatus))
	}

	return &dto.MeResponse{
		ID:          user.ID.String(),
		DisplayName: user.DisplayName,
		Role:        meRole,
		Status:      dto.UserStatus(user.Status),
	}, nil
}
