package models

import "github.com/google/uuid"

type UserStatus string

const USER_STATUS_ACTIVE UserStatus = "active"

type User struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"username"`
	DisplayName string     `json:"displayName"`
	Email       *string    `json:"email"`
	Status      UserStatus `json:"status"`
	Role        string     `json:"role"`
}
