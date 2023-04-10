package controllers

import (
	"apr-backend/client"
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type CompanyController struct {
	comServ services.CompanyService
}

func NewCompanyController(comServ services.CompanyService) CompanyController {
	return CompanyController{
		comServ: comServ,
	}
}

func (companyCtr CompanyController) CreateCompany(c *gin.Context) {
	principal := c.GetString(client.Principal)
	if principal == "" {
		log.Printf("Principal wasn't set in gin context\n")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var company model.Company
	if err := c.ShouldBindWith(&company, binding.JSON); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			if errors.Is(err, model.ErrInvalidDelatnost) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Delatnost": "Delatnost must be one of predefined values"})
				return
			}
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Must provide valid company as JSON"})
			return
		}
		errMsg := make(map[string]string)
		for _, e := range errs {
			errMsg[e.Field()] = model.CompanyErrors[e.Field()]
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errMsg)
		return
	}

	company.Vlasnik = principal
	err := companyCtr.comServ.SaveCompany(company)
	if err != nil {
		log.Printf("Couldn't save company: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, "Successfully created company")

}
