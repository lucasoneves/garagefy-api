# Integração Front-end — Recuperação de Senha

## Endpoints

### 1. Solicitar link de recuperação

```
POST /api/auth/forgot-password
Content-Type: application/json

{ "email": "usuario@email.com" }
```

**Resposta (sempre 200 — não revela se o e-mail existe):**
```json
{ "message": "Se o e-mail existir, você receberá um link de recuperação." }
```

---

### 2. Redefinir a senha

```
POST /api/auth/reset-password
Content-Type: application/json

{ "token": "6a8b...64 caracteres hex...", "password": "NovaSenha@123" }
```

**Resposta sucesso (200):**
```json
{ "message": "Senha redefinida com sucesso!" }
```

**Resposta erro (401):**
```json
{ "error": "Token inválido ou expirado" }
```

---

### 3. Página de redefinição (opcional)

```
GET /api/auth/reset-password?token=6a8b...
```

Retorna um formulário HTML pronto. Útil se você quiser redirecionar o usuário direto para a API em vez de criar uma página no front-end.

---

## Fluxo recomendado no front-end

**Tela 1 — "Esqueci minha senha"**
- Input: e-mail
- Botão "Enviar link"
- `POST /api/auth/forgot-password`
- Mostrar mensagem genérica de sucesso

**E-mail**
- O link enviado deve apontar para uma rota do **seu front-end** (ex: `https://meusite.com/reset-password?token=XXX`)
- Configure a URL base em `RESET_PASSWORD_URL` no `.env` da API

**Tela 2 — "Redefinir senha"**
- Ler `token` da query string (`?token=XXX`)
- Inputs: nova senha + confirmar senha
- Validar que as senhas conferem e têm >= 6 caracteres
- `POST /api/auth/reset-password` com `{ token, password }`
- Redirecionar para o login em caso de sucesso
- Mostrar erro caso o token tenha expirado

---

## Exemplo em React

```tsx
// Tela 1: Solicitar link
async function solicitarReset(email: string) {
  const res = await fetch('/api/auth/forgot-password', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email }),
  });
  return res.json();
}

// Tela 2: Redefinir senha
async function redefinirSenha(token: string, password: string) {
  const res = await fetch('/api/auth/reset-password', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ token, password }),
  });
  return res.json();
}
```

## Configuração de CORS

A API já libera `*` para `Origin`. Para produção, restrinja no `main.go`:

```go
c.Writer.Header().Set("Access-Control-Allow-Origin", "https://seudominio.com")
```
