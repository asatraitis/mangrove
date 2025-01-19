package models

import (
	"time"

	"github.com/google/uuid"
)

type UserToken struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"userId"`
	Expires time.Time `json:"expires"`
	User    *User     `json:"user,omitempty"`
}
