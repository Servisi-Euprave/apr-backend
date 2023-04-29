package controllers

import (
	"apr-backend/client"
	"apr-backend/internal/auth"
	"apr-backend/internal/db"
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"log"
	"net/http"
	"strconv"

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
// 400: errRes
// 500: errRes
func (controller AuthController) Login(c *gin.Context) {
	var creds model.CredentialsDto
	if err := c.ShouldBindBodyWith(&creds, binding.JSON); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := controller.authServ.CheckCredentials(creds); err != nil {
		switch err {
		case db.DatabaseError:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		default:
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
			return
		}
	}

	token, err := controller.jwtGenerator.GenerateAndSignJWT(creds.PIB, client.Apr)
	if err != nil {
		log.Printf("error: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, JwtResponse{Jwt: token})
}

// swagger:route POST /api/auth/login/:service auth SSOLogin
// Generate SSO token for service.
//
// Parameters:
//   - name: service
//     in: path
//     description: name of service for which to generate JWT
//     required: false
//     type: integer
//     format: string
//
// Security:
//   - bearerAuth:
//
// Responses:
// 201: jwtResponse
// 400: errRes
// 500: errRes
func (controller AuthController) SSOLogin(c *gin.Context) {
	serviceName := c.Param("service")
	principal := c.GetString("principal")
	pib, err := strconv.Atoi(principal)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JWT"})
		return
	}
	ssoToken, err := controller.jwtGenerator.GenerateAndSignJWT(pib, serviceName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "Couldn't generate JWT"})
		return
	}
	c.JSON(http.StatusOK, JwtResponse{Jwt: ssoToken})
}
