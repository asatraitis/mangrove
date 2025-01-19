package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asatraitis/mangrove/internal/handler/types"
	"github.com/google/uuid"
)

func TestGetUserIdFromCtx_OK(t *testing.T) {
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), types.REQ_CTX_KEY_USER_ID, userID.String())

	id, err := GetUserIdFromCtx(ctx)

	assert.NoError(t, err)
	assert.Equal(t, userID, id)
}

func TestGetUserIdFromCtx_FAIL_No_value(t *testing.T) {
	_, err := GetUserIdFromCtx(context.Background())

	assert.Error(t, err)
}

func TestGetUserIdFromCtx_FAIL_Not_UUID(t *testing.T) {
	ctx := context.WithValue(context.Background(), types.REQ_CTX_KEY_USER_ID, "123-abc")

	_, err := GetUserIdFromCtx(ctx)

	assert.Error(t, err)
}
