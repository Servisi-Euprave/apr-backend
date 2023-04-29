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
}
