package controllers

import (
	"apr-backend/client"
	"apr-backend/internal/auth"
	"apr-backend/internal/db"
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Response on successful login or registration, returns a valid JWT used for
// authentication.
// swagger:response jwtRes
type JwtResponse struct {
	// The JWT
	Jwt string `json:"jwt"`
}

// Error is used to specify what kind of error occured when processing request.
// swagger:response errRes
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is used to specify field errors on struct.
// swagger:response invalidBodyRes
type InvalidBodyResponse struct {
	ValidationErrors map[string]string
}

type CompanyController struct {
	comServ services.CompanyService
	jwtGen  auth.JwtGenerator
}

func NewCompanyController(comServ services.CompanyService, jwtGen auth.JwtGenerator) CompanyController {
	return CompanyController{
		comServ: comServ,
		jwtGen:  jwtGen,
	}
}

// swagger:route POST /api/company/ company CreateCompany
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

	err := companyCtr.comServ.SaveCompany(&company)
	if err != nil {
		log.Printf("Couldn't save company: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "Error when saving to database."})
		return
	}

	jwt, err := companyCtr.jwtGen.GenerateAndSignJWT(company.PIB, client.Apr)
	if err != nil {
		log.Printf("Error creating token: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, JwtResponse{Jwt: jwt})
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

// swagger:route GET /api/company/ company FindCompanies
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

// swagger:route GET /api/company/:pib company FindOne
// Finds one company by its pib
// Responses:
// 200: company
// 500: errRes
func (comCtr CompanyController) FindOne(c *gin.Context) {
	pibParam := c.Param("pib")
	pib, err := strconv.Atoi(pibParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Provided pib %s is invalid", pibParam)})
		return
	}

	company, err := comCtr.comServ.FindOne(pib)
	if errors.Is(err, db.NoSuchPibError) {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("Couldnt' find company with pib %s", pibParam)})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, company)
}
