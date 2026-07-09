package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luskation/ponto/internal/handler"
	"github.com/luskation/ponto/internal/middleware"
	"github.com/luskation/ponto/internal/repository"
	"github.com/luskation/ponto/internal/service"
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

	companyRepo := repository.NewCompanyRepository(pool)
	companyService := service.NewCompanyService(companyRepo)
	companyHandler := handler.NewCompanyHandler(companyService)

	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/forgot-password", authHandler.ForgotPassword)
		api.POST("/auth/reset-password", authHandler.ResetPassword)

		api.POST("/companies", middleware.AuthRequired(), middleware.AdminOnly(), companyHandler.Create)
		api.GET("/companies", middleware.AuthRequired(), middleware.AdminOnly(), companyHandler.List)
		api.GET("/companies/:id", middleware.AuthRequired(), middleware.AdminOnly(), companyHandler.GetByID)
		api.PUT("/companies/:id", middleware.AuthRequired(), middleware.AdminOnly(), companyHandler.Update)
		api.DELETE("/companies/:id", middleware.AuthRequired(), middleware.AdminOnly(), companyHandler.Delete)

		api.POST("/users", middleware.AuthRequired(), middleware.AdminOnly(), userHandler.Create)
		api.GET("/users", middleware.AuthRequired(), middleware.AdminOnly(), userHandler.List)
		api.GET("/users/:id", middleware.AuthRequired(), userHandler.GetByID)
		api.PUT("/users/:id", middleware.AuthRequired(), userHandler.Update)
		api.DELETE("/users/:id", middleware.AuthRequired(), middleware.AdminOnly(), userHandler.Delete)
	}

	log.Println("servidor rodando em :8080")
	r.Run(":8080")
}
