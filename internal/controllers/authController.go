package controllers

import (
	"apr-backend/client"
	"apr-backend/internal/auth"
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// NewAuthController creates a new AuthController
func NewAuthController(authServ services.AuthService, gen auth.JwtGenerator) AuthController {
	return AuthController{authServ: authServ, jwtGenerator: gen}
}

type AuthController struct {
	authServ     services.AuthService
	jwtGenerator auth.JwtGenerator
}

// Login is a handler for user login
func (controller AuthController) Login(c *gin.Context) {
	var creds model.CredentialsDto
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := controller.authServ.CheckCredentials(creds); err != nil {
		switch err {
		case services.DatabaseError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		default:
			log.Println(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
	}
	aud := client.Apr
	if creds.Service != "" {
		aud = creds.Service
	}

	claims := jwt.RegisteredClaims{
		Issuer:    client.Apr,
		Subject:   creds.Username,
		Audience:  []string{aud},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token, err := controller.jwtGenerator.SignJwt(claims)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"message": "User registered, error occured when creating jwt"})
		return
	}

	c.SetCookie("apr_session_jwt", token, 60*60*24, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}
