# 🌤️ Serviço de Clima por CEP (com OpenTelemetry + Zipkin)

Este projeto contém **dois microsserviços em Go** que, juntos, recebem um CEP e retornam o clima atual da cidade correspondente, com métricas e tracing distribuído via **OpenTelemetry** e **Zipkin**.

---

## 🧩 Estrutura do Projeto

```
.
├── servico-a/ # Serviço A: recebe o input do usuário
│ ├── main.go
│ ├── handler.go
│ └── tracer.go
│
├── servico-b/ # Serviço B: busca cidade e clima
│ ├── main.go
│ ├── handler.go
│ ├── cep.go
│ ├── weather.go
│ └── tracer.go
│
├── docker-compose.yml # Orquestra tudo (serviços + OTEL Collector + Zipkin)
├── Makefile # Facilita build e execução
└── README.md
```

---

## ⚙️ Requisitos

- Go 1.22+
- Docker + Docker Compose
- Conta gratuita no [WeatherAPI](https://www.weatherapi.com/) (necessário `API_KEY`)
- Internet (para consumir ViaCEP e WeatherAPI)

---

## 🚀 Execução com Docker

1. **Configure as variáveis de ambiente**

   Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:

   ```bash
   WEATHER_API_KEY=sua_chave_aqui
   OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4318/v1/traces
   ```

2. **Suba todo o ambiente**

   ```bash
   make up
   ```