package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luskation/ponto/internal/apperr"
)

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperr.AppError); ok {
		if appErr.Code == http.StatusInternalServerError {
			log.Printf("erro interno em %s %s: %v", c.Request.Method, c.Request.URL.Path, appErr.Cause)
		}
		c.JSON(appErr.Code, gin.H{"message": appErr.Message})
		return
	}
	log.Printf("erro não tratado em %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
	c.JSON(http.StatusInternalServerError, gin.H{"message": "erro interno"})
}
