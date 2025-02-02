package models

import "github.com/google/uuid"

type UserStatus string

const (
	USER_STATUS_ACTIVE    UserStatus = "active"
	USER_STATUS_INACTIVE  UserStatus = "inactive"
	USER_STATUS_SUSPENDED UserStatus = "suspended"
	USER_STATUS_PENDING   UserStatus = "pending"
)

type UserRole string

const (
	USER_ROLE_USER      UserRole = "user"
	USER_ROLE_ADMIN     UserRole = "admin"
	USER_ROLE_SUPERUSER UserRole = "superadmin"
)

type User struct {
	ID          uuid.UUID         `json:"id"`
	Username    string            `json:"username"`
	DisplayName string            `json:"displayName"`
	Email       *string           `json:"email"`
	Status      UserStatus        `json:"status"`
	Role        UserRole          `json:"role"`
	Token       *UserToken        `json:"token,omitempty"`
	Credentials []*UserCredential `json:"credentials,omitempty"`
}
