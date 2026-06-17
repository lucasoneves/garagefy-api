package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("garagefy_secret_key_2026_super_secure")

const ResetTokenDuration = 1 * time.Hour
const ResetTokenBytes = 32

// HashPassword transforma a senha em texto limpo em um hash seguro
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compara a senha do login com o hash salvo no banco
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken cria o token JWT assinado com o ID do usuário
func GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Token válido por 3 dias
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateResetToken gera um token aleatório seguro para recuperação de senha
func GenerateResetToken() (string, error) {
	bytes := make([]byte, ResetTokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ValidateToken valida a assinatura do JWT e retorna o ID do usuário se estiver tudo certo
func ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inválido")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		subStr, ok := claims["sub"].(string)
		if !ok {
			return uuid.Nil, errors.New("claim sub inválido")
		}

		userID, err := uuid.Parse(subStr)
		if err != nil {
			return uuid.Nil, errors.New("id do usuário inválido no token")
		}
		return userID, nil
	}

	return uuid.Nil, errors.New("token inválido")
}
