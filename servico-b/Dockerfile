# Etapa 1: build do binário
FROM golang:latest AS builder

WORKDIR /app

# Copia apenas os arquivos de dependências primeiro (para cache eficiente)
COPY go.mod go.sum ./

COPY ./go-mod-cache/mod /go/pkg/mod
COPY go.mod go.sum ./
RUN go mod verify

# Copia o código-fonte e compila o serviço
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /out/service-b .

# Etapa 2: imagem final
FROM scratch
WORKDIR /app

# Copia apenas o binário compilado
COPY --from=builder /out/service-b /app/service-b

# Expõe porta e define comando de inicialização
EXPOSE 8081
ENTRYPOINT ["/app/service-b"]
