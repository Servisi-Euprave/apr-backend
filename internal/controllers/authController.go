package controllers

import (
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewAuthController creates a new AuthController
func NewAuthController(authServ services.AuthService) AuthController {
	return AuthController{authServ: authServ}
}

type AuthController struct {
	authServ services.AuthService
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
	c.SetCookie("apr_session", "test", 3600, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}
