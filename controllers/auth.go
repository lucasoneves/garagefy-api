package controllers

import (
	"net/http"
	"strings"

	"garagefy-api/config" // Ajuste para o caminho real do seu projeto
	"garagefy-api/models"
	"garagefy-api/services"

	"github.com/gin-gonic/gin"
)

// POST /api/auth/register
func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Força o e-mail em minúsculas para evitar duplicidade por erro de digitação
	emailNormalized := strings.ToLower(strings.TrimSpace(input.Email))

	// Verifica se o e-mail já está cadastrado
	var existingUser models.User
	if err := config.DB.Where("email = ?", emailNormalized).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Este e-mail já está em uso"})
		return
	}

	// Encripta a senha com bcrypt antes de salvar
	hashedPassword, err := services.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar a senha"})
		return
	}

	newUser := models.User{
		Name:     input.Name,
		Email:    emailNormalized,
		Password: hashedPassword,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar o usuário"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário criado com sucesso!"})
}

// POST /api/auth/login
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailNormalized := strings.ToLower(strings.TrimSpace(input.Email))

	var user models.User
	// Busca o usuário pelo e-mail
	if err := config.DB.Where("email = ?", emailNormalized).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "E-mail ou senha inválidos"})
		return
	}

	// Valida a senha fornecida contra o hash salvo no banco
	if !services.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "E-mail ou senha inválidos"})
		return
	}

	// Gera o token JWT para o usuário autenticado
	token, err := services.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token de acesso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
