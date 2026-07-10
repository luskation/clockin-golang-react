package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luskation/ponto/internal/domain"
	"github.com/luskation/ponto/internal/service"
)

type TimeEntryHandler struct {
	service *service.TimeEntryService
}

func NewTimeEntryHandler(s *service.TimeEntryService) *TimeEntryHandler {
	return &TimeEntryHandler{service: s}
}

func (h *TimeEntryHandler) Register(c *gin.Context) {
	entry, err := h.service.RegisterEntry(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entry)
}

func (h *TimeEntryHandler) Update(c *gin.Context) {
	var input domain.TimeEntry
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload inválido"})
		return
	}
	input.ID = c.Param("id")
	if err := h.service.Update(c.Request.Context(), c.GetString("user_id"), c.GetString("role"), &input); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, input)
}

func (h *TimeEntryHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.GetString("user_id"), c.GetString("role"), c.Param("id")); err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TimeEntryHandler) ListMine(c *gin.Context) {
	page, limit := parsePagination(c)
	entries, total, err := h.service.ListByUser(c.Request.Context(), c.GetString("user_id"), page, limit)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, PaginatedResponse{Data: entries, Page: page, Limit: limit, Total: total})
}

func (h *TimeEntryHandler) ListAll(c *gin.Context) {
	page, limit := parsePagination(c)
	entries, total, err := h.service.ListAll(c.Request.Context(), page, limit)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, PaginatedResponse{Data: entries, Page: page, Limit: limit, Total: total})
}
