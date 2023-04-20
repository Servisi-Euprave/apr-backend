package controllers

import (
	"apr-backend/client"
	"apr-backend/internal/auth"
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// NewAuthController creates a new AuthController
func NewAuthController(authServ services.AuthService, gen auth.JwtGenerator) AuthController {
	return AuthController{authServ: authServ, jwtGenerator: gen}
}

type AuthController struct {
	authServ     services.AuthService
	jwtGenerator auth.JwtGenerator
}

// swagger:route POST /api/auth/login/ auth LoginUser
// Used for user authorization.
//
// Parameters:
// +name: credentials
// in: body
// type: credentials
// description: credentials with which to login
//
// Responses:
// 201: jwtResponse
// 400:
// 500:
func (controller AuthController) Login(c *gin.Context) {
	var creds model.CredentialsDto
	if err := c.ShouldBindBodyWith(&creds, binding.JSON); err != nil {
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
	aud := creds.Service
	if aud == "" {
		aud = client.Apr
	}

	token, err := controller.jwtGenerator.GenerateAndSignJWT(creds.Username, aud)
	if err != nil {
		log.Printf("error: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, JwtResponse{Jwt: token})
}
