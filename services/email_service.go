package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendResetPasswordEmail(to, token string) error {
	from := os.Getenv("SMTP_FROM")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	if from == "" || host == "" || port == "" {
		return nil
	}

	resetURL := fmt.Sprintf("%s?token=%s", os.Getenv("RESET_PASSWORD_URL"), token)
	if resetURL == "?token=" {
		resetURL = fmt.Sprintf("http://localhost:8080/api/auth/reset-password?token=%s", token)
	}

	subject := "Subject: Recuperação de Senha - Garagefy\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; padding: 20px;">
	<h2>Recuperação de Senha</h2>
	<p>Você solicitou a redefinição da sua senha no Garagefy.</p>
	<p>Clique no link abaixo para criar uma nova senha:</p>
	<p><a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px;">Redefinir Senha</a></p>
	<p>Ou copie e cole este link no seu navegador:</p>
	<p>%s</p>
	<p>Este link expira em 1 hora.</p>
	<p>Se você não solicitou esta redefinição, ignore este e-mail.</p>
</body>
</html>`, resetURL, resetURL)

	msg := subject + mime + body
	auth := smtp.PlainAuth("", user, pass, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}
