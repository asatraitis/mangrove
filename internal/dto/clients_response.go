package dto

type UserClientStatus string

const (
	CLIENT_STATUS_ACTIVE    UserClientStatus = "active"
	CLIENT_STATUS_PAUSED    UserClientStatus = "paused"
	CLIENT_STATUS_SUSPENDED UserClientStatus = "suspended"
)

type UserClient struct {
	ID          string           `json:"id"`
	UserID      string           `json:"userId"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Type        string           `json:"type"`
	RedirectURI string           `json:"redirectURI"`
	Status      UserClientStatus `json:"status"`
}

type UserClientsResponse []UserClient
