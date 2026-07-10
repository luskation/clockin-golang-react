package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luskation/ponto/internal/domain"
	"github.com/luskation/ponto/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Create(c *gin.Context) {
	var input domain.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload inválido"})
		return
	}
	if err := h.service.Create(c.Request.Context(), &input); err != nil {
		handleError(c, err)
		return
	}
	input.Password = ""
	c.JSON(http.StatusCreated, input)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	user, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)
	users, total, err := h.service.List(c.Request.Context(), page, limit)
	if err != nil {
		handleError(c, err)
		return
	}
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(http.StatusOK, PaginatedResponse{Data: users, Page: page, Limit: limit, Total: total})
}

func (h *UserHandler) Update(c *gin.Context) {
	var input domain.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload inválido"})
		return
	}
	input.ID = c.Param("id")
	if err := h.service.Update(c.Request.Context(), &input); err != nil {
		handleError(c, err)
		return
	}
	input.Password = ""
	c.JSON(http.StatusOK, input)
}

func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.GetString("user_id"), c.Param("id")); err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
