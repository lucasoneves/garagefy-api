# Teste de Recuperação de Senha no Postman

## Pré-requisitos

- API rodando em `http://localhost:8080`
- Postman instalado
- Um usuário já cadastrado via `POST /api/auth/register`

---

## 1. Solicitar recuperação de senha

**Endpoint:**
```
POST http://localhost:8080/api/auth/forgot-password
```

**Headers:**
| Key | Value |
|-----|-------|
| Content-Type | application/json |

**Body (raw JSON):**
```json
{
  "email": "usuario@email.com"
}
```

**Resposta esperada (sempre a mesma, por segurança):**
```json
{
  "message": "Se o e-mail existir, você receberá um link de recuperação."
}
```

> ⚠️ Em ambiente local sem SMTP configurado, o e-mail **não é enviado**. O token fica salvo no banco. Você precisa obtê-lo manualmente (próximo passo).

---

## 2. Obter o token de reset (ambiente local)

Conecte no banco PostgreSQL e execute:

```sql
SELECT email, reset_token, reset_token_expiry 
FROM users 
WHERE email = 'usuario@email.com';
```

Copie o valor da coluna `reset_token` (string hexadecimal de 64 caracteres).

---

## 3. Redefinir a senha

**Endpoint:**
```
POST http://localhost:8080/api/auth/reset-password
```

**Headers:**
| Key | Value |
|-----|-------|
| Content-Type | application/json |

**Body (raw JSON):**
```json
{
  "token": "6a8b...token copiado do banco...",
  "password": "NovaSenha@123"
}
```

**Resposta sucesso:**
```json
{
  "message": "Senha redefinida com sucesso!"
}
```

**Resposta erro (token inválido/expirado):**
```json
{
  "error": "Token inválido ou expirado"
}
```

---

## 4. Verificar a nova senha

Faça login com a nova senha:

**Endpoint:**
```
POST http://localhost:8080/api/auth/login
```

**Body (raw JSON):**
```json
{
  "email": "usuario@email.com",
  "password": "NovaSenha@123"
}
```

**Resposta esperada:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid-do-usuario",
    "name": "Nome do Usuário",
    "email": "usuario@email.com"
  }
}
```

---

## 5. Testar pelo navegador (alternativa ao Postman)

Abra no navegador:

```
http://localhost:8080/api/auth/reset-password?token=6a8b...token...
```

Um formulário HTML será renderizado para digitar a nova senha.
