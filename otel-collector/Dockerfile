# Etapa 1: usar a imagem oficial do collector como base para copiar o binário
FROM otel/opentelemetry-collector-contrib:0.74.0 AS collector

# Etapa 2: imagem Debian completa (com shell e utilitários)
FROM debian:bookworm-slim

# Atualiza pacotes e instala utilitários úteis
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates curl iputils-ping bash grep coreutils procps wget netcat-traditional vim less && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copia o binário do OpenTelemetry Collector da imagem base
COPY --from=collector /otelcol-contrib /otelcol-contrib

# Copia o arquivo de configuração para o local correto
COPY config.yaml /etc/otel/config.yaml

# (Opcional) Exibe o conteúdo do config para depuração durante build
RUN echo "=========== /etc/otel/config.yaml ===========" && cat /etc/otel/config.yaml && echo "============================================"

# Define as portas expostas
EXPOSE 4317 4318 8888

# Define o entrypoint
ENTRYPOINT ["/otelcol-contrib", "--config", "/etc/otel/config.yaml"]
