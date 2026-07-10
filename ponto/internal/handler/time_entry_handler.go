package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *TimeEntryHandler) ListMine(c *gin.Context) {
	entries, err := h.service.ListByUser(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *TimeEntryHandler) ListAll(c *gin.Context) {
	entries, err := h.service.ListAll(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, entries)
}
