# Estágio de Compilação (Build)
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o resto do código fonte
COPY . .

# Compila o binário da aplicação estaticamente
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Estágio de Execução (Runtime leve)
FROM alpine:latest

WORKDIR /app

# Copia o binário compilado do estágio anterior
COPY --from=builder /app/main .

# Cria a pasta de uploads para garantir que ela exista no container
RUN mkdir -p uploads

# Expõe a porta que a API escuta
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./main"]