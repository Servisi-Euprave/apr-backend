package controllers

import (
	"apr-backend/client"
	"apr-backend/internal/db"
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"errors"
	"log"
	"net/http"
	"strconv"

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

// swagger:route POST /api/companies/ company CreateCompany
// Registers a new company for logged-in user.
//
// Parameters:
// +name: company
// in: body
// type: company
// description: company to be created
//
// Security:
// bearerAuth:
//
// Responses:
// 201: company Company which was created
// 400:
// 500:
// This text will appear as description of your response body.
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
	c.JSON(http.StatusCreated, company)
}

const (
	pageQuery      = "page"
	sortQuery      = "order"
	delatnostQuery = "delatnost"
	sedisteQuery   = "sediste"
	mestoQuery     = "mesto"
	ascQuery       = "asc"
)

const (
	naziv int = iota
)

// swagger:route GET /api/companies/ company FindCompanies
// Filters, sorts and paginates companies
//
// Parameters:
// +name: page
// in: query
// required: false
// type: integer
// format: int32
// +name: order
// in: query
// required: false
// type: string
// +name: asc
// in: query
// required: false
// type: boolean
// description: whether to sort ascending or descending.
// +name: delatnost
// in: query
// required: false
// type: string
// description: delatnost by which to filter by, must be a valid delatnost
// +name: mesto
// in: query
// required: false
// type: string
// description: mesto by which to filter
//
// Responses:
// 201: []company
// 400:
// 500:
func (companyCtr CompanyController) FindCompanies(c *gin.Context) {
	pageStr, ok := c.GetQuery(pageQuery)
	page := 0
	var err error
	if ok {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			page = 0
		}
	}
	column, ok := c.GetQuery(sortQuery)
	if !ok {
		column = "PIB"
	}
	delatnost := c.Query(delatnostQuery)
	sediste := c.Query(sedisteQuery)
	mesto := c.Query(mestoQuery)
	ascStr, _ := c.GetQuery(ascQuery)
	asc := true
	if ascStr == "false" {
		asc = false
	}

	companies, err := companyCtr.comServ.FindCompanies(model.CompanyFilter{
		OrderBy:   column,
		Asc:       asc,
		Page:      page,
		Mesto:     mesto,
		Sediste:   sediste,
		Delatnost: delatnost,
	})

	if errors.Is(err, db.DatabaseError) {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if errors.Is(err, db.InvalidFilter) {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, companies)
	return
}
