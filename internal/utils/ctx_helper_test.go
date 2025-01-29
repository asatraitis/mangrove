package utils

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

func TestGetUserStatusFromCtx_OK(t *testing.T) {
	ctx := context.WithValue(context.Background(), types.REQ_CTX_KEY_USER_STATUS, "active")

	status, err := GetUserStatusFromCtx(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "active", status)
}

func TestGetUserStatusFromCtx_FAIL_NoValue(t *testing.T) {
	status, err := GetUserStatusFromCtx(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "", status)
}

func TestGetUserStatusFromCtx_FAIL_BadType(t *testing.T) {
	ctx := context.WithValue(context.Background(), types.REQ_CTX_KEY_USER_STATUS, int32(1))

	status, err := GetUserStatusFromCtx(ctx)
	assert.Error(t, err)
	assert.Equal(t, "", status)
}
