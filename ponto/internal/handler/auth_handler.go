package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luskation/ponto/internal/auth"
	"github.com/luskation/ponto/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(s *service.UserService) *AuthHandler {
	return &AuthHandler{userService: s}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email e senha são obrigatórios"})
		return
	}

	user, err := h.userService.GetByEmail(c.Request.Context(), input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "credenciais inválidas"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "credenciais inválidas"})
		return
	}

	token, err := auth.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "erro ao gerar token"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var input ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email é obrigatório"})
		return
	}

	user, err := h.userService.GetByEmail(c.Request.Context(), input.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "usuário não encontrado"})
		return
	}

	token, err := auth.GenerateResetToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "erro ao gerar token"})
		return
	}

	// Nota: em produção o token deveria ser enviado por e-mail, nunca no response.
	c.JSON(http.StatusOK, gin.H{"reset_token": token})
}

type ResetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "token e new_password são obrigatórios"})
		return
	}

	userID, err := auth.ParseResetToken(input.Token)
	if err != nil || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "token inválido ou expirado"})
		return
	}

	if err := h.userService.UpdatePassword(c.Request.Context(), userID, input.NewPassword); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "senha atualizada com sucesso"})
}
