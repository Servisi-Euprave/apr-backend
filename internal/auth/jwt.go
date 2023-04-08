package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type JwtGenerator interface {
	SignJwt(claims jwt.RegisteredClaims) (string, error)
	ParseJwt(token string) (jwt.RegisteredClaims, error)
}

func NewJwtGenerator(privateKeyFile string) (JwtGenerator, error) {
	var generator jwtGeneratorRsa
	keyFile, err := os.Open(privateKeyFile)
	if err != nil {
		generator.key, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return generator, fmt.Errorf("No key provided, error generating key: %w", err)
		}
		return generator, nil
	}
	defer keyFile.Close()
	keyData, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading from key file: %w", err)
	}
	generator.key, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return generator, fmt.Errorf("Error parsing PMA encoded key: %w", err)
	}
	return generator, nil
}

type jwtGeneratorRsa struct {
	key         *rsa.PrivateKey
	serviceName string
}

func (jwtGen jwtGeneratorRsa) SignJwt(claims jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	signed, err := token.SignedString(jwtGen.key)
	if err != nil {
		return "", fmt.Errorf("Error signing token: %w", err)
	}
	return signed, nil
}

func (jwtGen jwtGeneratorRsa) keyFunc(t *jwt.Token) (interface{}, error) {
	return jwtGen.key.Public(), nil
}

func (jwtGen jwtGeneratorRsa) ParseJwt(token string) (jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	parsed, err := jwt.ParseWithClaims(token, &claims, jwtGen.keyFunc)
	if err != nil {
		return claims, fmt.Errorf("Cannot parse jwt: %w", err)
	}

	// Verify method
	err = parsed.Method.Verify(jwt.SigningMethodRS512.Alg(), parsed.Signature, jwtGen.key.Public())
	if err != nil {
		return claims, fmt.Errorf("Wrong signing method: %w", err)
	}

	return claims, claims.Valid()
}
