package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func parsePagination(c *gin.Context) (page, limit int) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}

type PaginatedResponse struct {
	Data  any `json:"data"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}
