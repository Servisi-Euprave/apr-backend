package auth

import (
	"apr-backend/client"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtGenerator interface {
	SignJwt(claims jwt.RegisteredClaims) (string, error)
	GenerateAndSignJWT(username, audience string) (string, error)
	client.JwtVerifier
}

func ReadRSAPrivateKeyFromFile(keyFilePath string) (*rsa.PrivateKey, error) {
	var key *rsa.PrivateKey
	keyFile, err := os.Open(keyFilePath)
	if err != nil {
		key, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return key, fmt.Errorf("No key provided, error generating key: %w", err)
		}
		return key, nil
	}
	defer keyFile.Close()
	keyData, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading from key file: %w", err)
	}
	key, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return key, fmt.Errorf("Error parsing PMA encoded key: %w", err)
	}
	return key, nil
}

func NewJwtGenerator(key *rsa.PrivateKey) JwtGenerator {
	generator := jwtGeneratorRsa{
		key:         key,
		serviceName: client.Apr,
		JwtVerifier: client.NewVerifier(key.Public()),
	}
	return generator
}

type jwtGeneratorRsa struct {
	key         *rsa.PrivateKey
	serviceName string
	client.JwtVerifier
}

// GenerateAndSignJWT implements JwtGenerator
func (jwtGen jwtGeneratorRsa) GenerateAndSignJWT(username string, audience string) (string, error) {

	aud := client.Apr
	if audience != "" {
		aud = audience
	}

	claims := jwt.RegisteredClaims{
		Issuer:    client.Apr,
		Subject:   username,
		Audience:  []string{aud},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	return jwtGen.SignJwt(claims)
}

func (jwtGen jwtGeneratorRsa) SignJwt(claims jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	signed, err := token.SignedString(jwtGen.key)
	if err != nil {
		return "", fmt.Errorf("Error signing token: %w", err)
	}
	return signed, nil
}
