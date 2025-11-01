# Etapa 1: build do binário
FROM golang:latest AS builder

WORKDIR /app

# Copia apenas os arquivos de dependências primeiro (para cache eficiente)
COPY go.mod go.sum ./

COPY go.mod go.sum ./
RUN go mod download

# Copia o código-fonte e compila o serviço
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /out/servico-b .

# Etapa 2: imagem final
FROM debian:bookworm-slim
WORKDIR /app

# Copia apenas o binário compilado
COPY --from=builder /out/servico-b /app/servico-b
RUN apt-get update && apt-get install -y grep coreutils wget iputils-ping curl ca-certificates && update-ca-certificates

# Expõe porta e define comando de inicialização
EXPOSE 8081
ENTRYPOINT ["/app/servico-b"]
