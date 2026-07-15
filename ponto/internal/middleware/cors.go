package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS libera o frontend (origem diferente da API) a chamar a API a partir do navegador.
// Não usamos cookies/credentials, só o header Authorization, então liberar "*" é seguro aqui.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
