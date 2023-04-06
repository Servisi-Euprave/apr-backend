package controllers

import (
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func NewUserController(userServ services.UserService) UserController {
	return UserController{userService: userServ}
}

type UserController struct {
	userService services.UserService
}

func (usrCtr UserController) RegisterUser(c *gin.Context) {
	var usr model.User
	if err := c.BindJSON(&usr); err != nil {
		errs := err.(validator.ValidationErrors)
		errMsg := make(map[string]string)
		for _, e := range errs {
			errMsg[e.Field()] = e.Error()
		}
		c.JSON(http.StatusBadRequest, errMsg)
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

	//TODO: Save session and return actual session ID
	c.SetCookie("apr_session", "test", 3600, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}
