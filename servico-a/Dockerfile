# Etapa 1: build do binário
FROM golang:latest AS builder

WORKDIR /app

# Copia apenas os arquivos de dependências primeiro (melhor cache)
COPY go.mod go.sum ./

COPY ./go-mod-cache/mod /go/pkg/mod
COPY go.mod go.sum ./
RUN go mod verify

# Copia o código-fonte e compila
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /out/service-a .

# Etapa 2: imagem final (mínima)
FROM scratch
WORKDIR /app

# Copia apenas o binário compilado
COPY --from=builder /out/service-a /app/service-a

# Define porta e comando de entrada
EXPOSE 8080
ENTRYPOINT ["/app/service-a"]
