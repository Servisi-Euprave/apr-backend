package model

type Credentials struct {
	Username     string
	PasswordHash []byte
}

type CredentialsDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Service  string `json:"service,omitempty"`
}
