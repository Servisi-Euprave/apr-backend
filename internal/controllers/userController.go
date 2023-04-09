package controllers

import (
	"apr-backend/client"
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const ServiceID = "apr"

func NewUserController(userServ services.UserService) UserController {
	return UserController{userService: userServ}
}

type UserController struct {
	userService services.UserService
}

func (usrCtr UserController) RegisterUser(c *gin.Context) {
	var usr model.User

	if err := c.ShouldBindBodyWith(&usr, binding.JSON); err != nil {
		errs := err.(validator.ValidationErrors)
		errMsg := make(map[string]string)
		for _, e := range errs {
			errMsg[e.Field()] = e.Error()
		}
		c.AbortWithError(http.StatusInternalServerError, err)
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
}

func (usrCtr UserController) GetUserByUsername(c *gin.Context) {
	if principal, ok := c.Get(client.Principal); ok {
		log.Printf("Principal: %s\n", principal)
	}
	c.JSON(http.StatusOK, gin.H{"username": c.Param("username")})
}
