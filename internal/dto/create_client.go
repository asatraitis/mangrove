package dto

type UserClientKeyAlgo string

const (
	CLIENT_KEY_ALGO_EDDSA UserClientKeyAlgo = "EdDSA"
)

// TODO: create common fields struct and embed insetad of repeating?
type CreateClientRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	RedirectURI string            `json:"redirectURI"`
	Status      UserClientStatus  `json:"status"`
	PublicKey   []byte            `json:"publicKey"`
	KeyAlgo     UserClientKeyAlgo `json:"keyAlgo"`
}

type CreateClientResponse UserClient
