package model

type Credentials struct {
	Username     string
	PasswordHash []byte
}

// Credentials with which to login
// swagger:model credentials
type CredentialsDto struct {
	// Required: true
	Username string `json:"username"`
	// Required: true
	Password string `json:"password"`
	// Service for which to issue token on succesful login. This name
	// will be written to aud claim in JWT
	// Example: javne_nabavke
	Service string `json:"service,omitempty"`
}
