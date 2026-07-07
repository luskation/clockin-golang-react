package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luskation/ponto/internal/apperr"
)

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperr.AppError); ok {
		c.JSON(appErr.Code, gin.H{"message": appErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"message": "erro interno"})
}
