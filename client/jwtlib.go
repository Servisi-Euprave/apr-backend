package client

import (
	"crypto/rsa"
	"fmt"
	"net/http"
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
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if !claims.VerifyAudience(serviceName, true) {
			errMsg, _ := fmt.Printf("invalid audience: expected %q, got %v", serviceName, claims.Audience)
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

type JwtVerifier interface {
	ParseJwt(token string) (jwt.RegisteredClaims, error)
}

// NewVerifier() returns a new JwtVerifier instance
// @pubKey public RSA key of the signee
// @serviceName used to verify `aud` claims
func NewVerifier(pubKey rsa.PublicKey) JwtVerifier {
	return defaultJwtVerifier{pubKey: pubKey}
}

type defaultJwtVerifier struct {
	pubKey rsa.PublicKey
}

func (jwtGen defaultJwtVerifier) keyFunc(t *jwt.Token) (interface{}, error) {
	return jwtGen.pubKey, nil
}

// CheckJwt implements JwtVerifier
func (jwtGen defaultJwtVerifier) ParseJwt(token string) (jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	parsed, err := jwt.ParseWithClaims(token, &claims, jwtGen.keyFunc)
	if err != nil {
		return claims, fmt.Errorf("Cannot parse jwt: %w", err)
	}

	// Verify method
	err = parsed.Method.Verify(jwt.SigningMethodRS512.Alg(), parsed.Signature, jwtGen.pubKey)
	if err != nil {
		return claims, fmt.Errorf("Wrong signing method: %w", err)
	}
	return claims, claims.Valid()
}
