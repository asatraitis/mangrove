package typeconv

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asatraitis/mangrove/internal/dal/models"
	"github.com/asatraitis/mangrove/internal/dto"
	"github.com/google/uuid"
)

func TestConvertUserToMeResponse_OK(t *testing.T) {
	user := &models.User{
		ID:          uuid.New(),
		Username:    "test-user",
		DisplayName: "test-display",
		Status:      models.USER_STATUS_ACTIVE,
		Role:        models.USER_ROLE_USER,
	}

	me, err := ConvertUserToMeResponse(user)

	assert.NoError(t, err)
	assert.NotNil(t, me)
	assert.Equal(t, user.ID.String(), me.ID)
	assert.Equal(t, "test-display", me.DisplayName)
	assert.Equal(t, dto.UserStatus("active"), me.Status)
	assert.Equal(t, dto.UserRole("user"), me.Role)
}

func TestConvertUserToMeResponse_FAIL_BadStatus(t *testing.T) {
	user := &models.User{
		ID:          uuid.New(),
		Username:    "test-user",
		DisplayName: "test-display",
		Status:      models.UserStatus("random"),
		Role:        models.USER_ROLE_USER,
	}

	me, err := ConvertUserToMeResponse(user)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "user status is not available in the list of statuses")
	assert.Nil(t, me)
}

func TestConvertUserToMeResponse_FAIL_BadRole(t *testing.T) {
	user := &models.User{
		ID:          uuid.New(),
		Username:    "test-user",
		DisplayName: "test-display",
		Status:      models.USER_STATUS_ACTIVE,
		Role:        models.UserRole("random"),
	}

	me, err := ConvertUserToMeResponse(user)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "user role is not available in the list of roles")
	assert.Nil(t, me)
}
