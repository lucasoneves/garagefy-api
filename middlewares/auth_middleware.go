package middlewares

import (
	"net/http"
	"strings"

	"garagefy-api/services" // Ajuste para o caminho do seu projeto

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Busca o cabeçalho 'Authorization'
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Cabeçalho de autorização é obrigatório"})
			c.Abort() // Interrompe a requisição aqui mesmo
			return
		}

		// 2. O formato esperado é 'Bearer <TOKEN>', então dividimos a string
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido. Use 'Bearer <TOKEN>'"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Valida o token chamando o serviço que criamos antes
		userID, err := services.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado: " + err.Error()})
			c.Abort()
			return
		}

		// 4. Injeta o userID validado dentro do contexto da requisição.
		// Qualquer controller que rodar após esse middleware poderá ler esse valor.
		c.Set("userID", userID)

		// 5. Continua o fluxo para o próximo handler/controller
		c.Next()
	}
}
