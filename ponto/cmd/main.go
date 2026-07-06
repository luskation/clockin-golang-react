package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luskation/ponto/internal/repository"
)

func main() {
	// Carrega variáveis do .env
	if err := godotenv.Load(".env"); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatalf("erro ao carregar .env: %v", err)
		}
	}

	// Conecta ao banco
	ctx := context.Background()
	pool, err := repository.NewPool(ctx)
	if err != nil {
		log.Fatalf("falha na conexão com banco: %v", err)
	}
	defer pool.Close()

	log.Println("banco conectado com sucesso")

	// Servidor HTTP
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("servidor rodando em :8080")
	r.Run(":8080")
}
