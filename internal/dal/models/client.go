package models

import (
	"time"

	"github.com/google/uuid"
)

type ClientStatus string

const (
	CLIENT_STATUS_ACTIVE    ClientStatus = "active"
	CLIENT_STATUS_PAUSED    ClientStatus = "paused"
	CLIENT_STATUS_SUSPENDED ClientStatus = "suspended"
)

type ClientKeyAlgo string

const (
	CLIENT_KEY_ALGO_EDDSA ClientKeyAlgo = "EdDSA"
)

type Client struct {
	ID           uuid.UUID     `json:"id"`
	UserID       uuid.UUID     `json:"userId"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	RedirectURI  string        `json:"redirectURI"`
	PublicKey    []byte        `json:"publicKey"`
	KeyExpiresAt time.Time     `json:"keyExpiresAt"`
	KeyAlgo      ClientKeyAlgo `json:"keyAlgo"`
	Status       ClientStatus  `json:"status"`
}
