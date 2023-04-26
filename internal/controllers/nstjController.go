package controllers

import (
	"apr-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NstjController struct {
	nstjServ services.NstjService
}

func NewNstjController(nstjService services.NstjService) NstjController {
	return NstjController{nstjServ: nstjService}
}

// swagger:route GET /api/nstj/ nstj FindAll
// Gets all available NSTJ codes
// Responses:
// 200: []nstj
// 500: errRes
func (nstjCtr NstjController) FindAll(c *gin.Context) {
	services, err := nstjCtr.nstjServ.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, services)
}
