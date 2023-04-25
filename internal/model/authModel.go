package model

type Credentials struct {
	Username     string
	PasswordHash []byte
}

// Login credentials
//
// Structure used during authentication.
// swagger:model credentials
type CredentialsDto struct {
	// Required: true
	PIB int `json:"pib"`
	// Required: true
	Password string `json:"password"`
	// Name of service used in aud claim in the JWT.
	// Example: javne_nabavke
	Service string `json:"service,omitempty"`
}
