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

// swagger:response
// Response on successful login or registration, returns a valid JWT used for
// authentication.
type JwtResponse struct {
	// The JWT
	Jwt string `json:"jwt"`
}

// swagger:response
// Error is used to specify what kind of error occured when processing request.
type ErrorResponse struct {
	Error string `json:"error"`
}

// swagger:response
// Error is used to specify what kind of error occured when processing request.
type InvalidBodyResponse struct {
	ValidationErrors map[string]string
}

// swagger:route POST /api/user/ users RegisterUser
// Registers a new user
// Parameters:
// +name: user
// in: body
// type: user
// description: User which to register
//
// Responses:
// 201: jwtResponse
// 400: InvalidStruct
// 500: errorResponse
func (usrCtr UserController) RegisterUser(c *gin.Context) {
	var usr model.User

	if err := c.ShouldBindBodyWith(&usr, binding.JSON); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "Must provide user as JSON."})
			return
		}
		errMsg := make(map[string]string)
		for _, e := range errs {
			errMsg[e.Field()] = model.UserErrors[e.Field()]
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, InvalidBodyResponse{ValidationErrors: errMsg})
		return
	}

	err := usrCtr.userService.SaveUser(usr)
	if err == services.DatabaseError {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{"Error when connecting to database."})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{"Internal server error."})
		return
	}

	jwt, err := usrCtr.jwtGen.GenerateAndSignJWT(usr.Username, client.Apr)
	if err != nil {
		log.Printf("Error creating token: %s", err.Error())
	}
	c.JSON(http.StatusCreated, JwtResponse{Jwt: jwt})
	return

}

func (usrCtr UserController) GetUserByUsername(c *gin.Context) {
	if principal, ok := c.Get(client.Principal); ok {
		log.Printf("Principal: %s\n", principal)
	}
	c.JSON(http.StatusOK, gin.H{"username": c.Param("username")})
}
