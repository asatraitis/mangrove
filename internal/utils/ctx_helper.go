package utils

import (
	"context"
	"errors"

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
