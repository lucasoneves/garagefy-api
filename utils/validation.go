package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		field := ve[0].Field()
		tag := ve[0].Tag()
		param := ve[0].Param()

		switch tag {
		case "required":
			return "O campo " + field + " é obrigatório"
		case "email":
			return "O campo " + field + " deve ser um e-mail válido"
		case "min":
			return "O campo " + field + " deve ter pelo menos " + param + " caracteres"
		case "max":
			return "O campo " + field + " deve ter no máximo " + param + " caracteres"
		case "oneof":
			return "O campo " + field + " deve ser um dos valores: " + strings.ReplaceAll(param, " ", ", ")
		default:
			return "O campo " + field + " é inválido"
		}
	}
	return "Dados inválidos"
}
