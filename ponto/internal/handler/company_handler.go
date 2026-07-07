package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luskation/ponto/internal/domain"
	"github.com/luskation/ponto/internal/service"
)

type CompanyHandler struct {
	service *service.CompanyService
}

func NewCompanyHandler(s *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: s}
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var input domain.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload inválido"})
		return
	}
	if err := h.service.Create(c.Request.Context(), &input); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, input)
}

func (h *CompanyHandler) GetByID(c *gin.Context) {
	company, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) List(c *gin.Context) {
	companies, err := h.service.List(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, companies)
}

func (h *CompanyHandler) Update(c *gin.Context) {
	var input domain.Company
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload inválido"})
		return
	}
	input.ID = c.Param("id")
	if err := h.service.Update(c.Request.Context(), &input); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, input)
}

func (h *CompanyHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
