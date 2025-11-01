# Etapa 1: build do binário
FROM golang:latest AS builder

WORKDIR /app

# Copia apenas os arquivos de dependências primeiro (melhor cache)
COPY go.mod go.sum ./

COPY go.mod go.sum ./
RUN go mod download

# Copia o código-fonte e compila
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /out/servico-a .

# Etapa 2: imagem final (mínima)
FROM debian:bookworm-slim
WORKDIR /app

# Copia apenas o binário compilado
COPY --from=builder /out/servico-a /app/servico-a
RUN apt-get update && apt-get install -y grep coreutils wget iputils-ping curl ca-certificates && update-ca-certificates

# Define porta e comando de entrada
EXPOSE 8080
ENTRYPOINT ["/app/servico-a"]
