package utils

import (
	"context"
	"errors"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/handler/types"
	"github.com/google/uuid"
)

func GetUserIdFromCtx(ctx context.Context) (uuid.UUID, error) {
	value := ctx.Value(types.REQ_CTX_KEY_USER_ID)
	s, ok := value.(string)
	if !ok {
		return uuid.Nil, errors.New("failed type assertion")
	}
	userID, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errors.New("failed to parse user uuid")
	}

	return userID, nil
}

func GetUserStatusFromCtx(ctx context.Context) (models.UserStatus, error) {
	value := ctx.Value(types.REQ_CTX_KEY_USER_STATUS)
	s, ok := value.(models.UserStatus)
	if !ok {
		return "", errors.New("failed type assertion - UserStatus")
	}
	return s, nil
}

func GetUserRoleFromCtx(ctx context.Context) (models.UserRole, error) {
	value := ctx.Value(types.REQ_CTX_KEY_USER_ROLE)
	s, ok := value.(models.UserRole)
	if !ok {
		return "", errors.New("failed type assertion - UserRole")
	}
	return s, nil
}
