package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"garagefy-api/config"
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

// POST /api/auth/forgot-password
func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailNormalized := strings.ToLower(strings.TrimSpace(input.Email))

	var user models.User
	if err := config.DB.Where("email = ?", emailNormalized).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Se o e-mail existir, você receberá um link de recuperação."})
		return
	}

	token, err := services.GenerateResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token de recuperação"})
		return
	}

	expiry := time.Now().Add(services.ResetTokenDuration)
	user.ResetToken = &token
	user.ResetTokenExpiry = &expiry

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar token de recuperação"})
		return
	}

	if err := services.SendResetPasswordEmail(user.Email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao enviar e-mail de recuperação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Se o e-mail existir, você receberá um link de recuperação."})
}

// POST /api/auth/reset-password
func ResetPassword(c *gin.Context) {
	var input struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("reset_token = ?", input.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		return
	}

	if user.ResetTokenExpiry == nil || time.Now().After(*user.ResetTokenExpiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		return
	}

	hashedPassword, err := services.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar a senha"})
		return
	}

	user.Password = hashedPassword
	user.ResetToken = nil
	user.ResetTokenExpiry = nil

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao redefinir a senha"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Senha redefinida com sucesso!"})
}

// GET /api/auth/reset-password (renderiza formulário HTML do link do e-mail)
func ResetPasswordForm(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token não fornecido"})
		return
	}

	// Verifica se o token é válido antes de mostrar o formulário
	var user models.User
	if err := config.DB.Where("reset_token = ?", token).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		return
	}

	if user.ResetTokenExpiry == nil || time.Now().After(*user.ResetTokenExpiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Redefinir Senha - Garagefy</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: Arial, sans-serif; background: #f5f5f5; display: flex; justify-content: center; align-items: center; min-height: 100vh; }
		.card { background: white; padding: 40px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); width: 100%%; max-width: 400px; }
		h1 { font-size: 24px; margin-bottom: 8px; color: #333; }
		p { color: #666; margin-bottom: 24px; }
		input { width: 100%%; padding: 12px; border: 1px solid #ddd; border-radius: 4px; font-size: 16px; margin-bottom: 16px; }
		button { width: 100%%; padding: 12px; background: #4CAF50; color: white; border: none; border-radius: 4px; font-size: 16px; cursor: pointer; }
		button:hover { background: #45a049; }
		.error { color: #f44336; margin-bottom: 16px; display: none; }
	</style>
</head>
<body>
	<div class="card">
		<h1>Redefinir Senha</h1>
		<p>Digite sua nova senha.</p>
		<div class="error" id="error"></div>
		<input type="password" id="password" placeholder="Nova senha (mínimo 6 caracteres)" minlength="6" />
		<input type="password" id="confirmPassword" placeholder="Confirmar nova senha" />
		<button onclick="resetPassword()">Redefinir Senha</button>
	</div>
	<script>
		async function resetPassword() {
			const password = document.getElementById('password').value;
			const confirm = document.getElementById('confirmPassword').value;
			const error = document.getElementById('error');

			if (password.length < 6) {
				error.style.display = 'block';
				error.textContent = 'A senha deve ter no mínimo 6 caracteres.';
				return;
			}
			if (password !== confirm) {
				error.style.display = 'block';
				error.textContent = 'As senhas não conferem.';
				return;
			}

			error.style.display = 'none';
			const token = new URLSearchParams(window.location.search).get('token');

			try {
				const res = await fetch('/api/auth/reset-password', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ token, password })
				});
				const data = await res.json();
				if (res.ok) {
					document.querySelector('.card').innerHTML = '<h1>Senha Redefinida</h1><p>Sua senha foi redefinida com sucesso! Você já pode fazer login.</p>';
				} else {
					error.style.display = 'block';
					error.textContent = data.error || 'Erro ao redefinir senha.';
				}
			} catch (e) {
				error.style.display = 'block';
				error.textContent = 'Erro de conexão. Tente novamente.';
			}
		}
	</script>
</body>
</html>`)))
}
