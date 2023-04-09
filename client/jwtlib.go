package client

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const Principal = "principal"
const Apr = "apr"

// Middleware for gin, check if user is logged in
func CheckAuth(verifier JwtVerifier, serviceName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			return
		}

		//Check header format
		bearer, tokenStr, found := strings.Cut(bearerToken, " ")
		if !found || bearer != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid format for bearer token"})
			return
		}

		//Check that jwt is valid
		claims, err := verifier.ParseJwt(tokenStr)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if !claims.VerifyAudience(serviceName, true) {
			errMsg := fmt.Sprintf("invalid audience: expected %q, got %v", serviceName, claims.Audience)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			return
		}

		if !claims.VerifyIssuer(Apr, true) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid issuer"})
			return
		}

		ctx.Set(Principal, claims.Subject)
	}
}

func ReadRSAPublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	var key *rsa.PublicKey
	keyFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("No key provided, error generating key: %w", err)
	}
	defer keyFile.Close()
	keyData, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading from key file: %w", err)
	}
	key, err = jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return key, fmt.Errorf("Error parsing PMA encoded key: %w", err)
	}
	return key, nil
}

type JwtVerifier interface {
	ParseJwt(token string) (jwt.RegisteredClaims, error)
}

// NewVerifier() returns a new JwtVerifier instance
// @pubKey public RSA key of the signee
// @serviceName used to verify `aud` claims
func NewVerifier(pubKey crypto.PublicKey) JwtVerifier {
	return defaultJwtVerifier{pubKey: pubKey}
}

type defaultJwtVerifier struct {
	pubKey crypto.PublicKey
}

func (jwtGen defaultJwtVerifier) keyFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
	}
	return jwtGen.pubKey, nil
}

// CheckJwt implements JwtVerifier
func (jwtGen defaultJwtVerifier) ParseJwt(token string) (jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(token, &claims, jwtGen.keyFunc)
	if err != nil {
		return claims, fmt.Errorf("Cannot parse jwt: %w", err)
	}
	return claims, claims.Valid()
}
