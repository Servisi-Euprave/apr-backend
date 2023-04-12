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
	"github.com/go-playground/validator/v10"
)

const ServiceID = "apr"

func NewUserController(userServ services.UserService, jwtGen auth.JwtGenerator) UserController {
	return UserController{
		userService: userServ,
		jwtGen:      jwtGen,
	}
}

type UserController struct {
	userService services.UserService
	jwtGen      auth.JwtGenerator
}

func (usrCtr UserController) RegisterUser(c *gin.Context) {
	var usr model.User

	if err := c.ShouldBindBodyWith(&usr, binding.JSON); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Must provide user as JSON")
			return
		}
		errMsg := make(map[string]string)
		for _, e := range errs {
			errMsg[e.Field()] = model.UserErrors[e.Field()]
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errMsg)
		return
	}

	err := usrCtr.userService.SaveUser(usr)
	if err == services.DatabaseError {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	usrCtr.jwtGen.GenerateAndSignJWT(usr.Username, client.Apr)
}

func (usrCtr UserController) GetUserByUsername(c *gin.Context) {
	if principal, ok := c.Get(client.Principal); ok {
		log.Printf("Principal: %s\n", principal)
	}
	c.JSON(http.StatusOK, gin.H{"username": c.Param("username")})
}
