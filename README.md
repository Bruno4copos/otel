# 🌤️ Serviço de Clima por CEP com Observabilidade (OpenTelemetry + Zipkin)

Este projeto demonstra **observabilidade distribuída** em uma arquitetura de **microsserviços Go**, que recebem um CEP, localizam a cidade e retornam o clima atual com métricas e tracing via **OpenTelemetry** e **Zipkin**.


   - **Serviço A (API CEP)** → recebe um CEP (8 dígitos) via `POST /cep`, valida e repassa ao Serviço B.
   - **Serviço B (API Clima)** → consulta a [ViaCEP](https://viacep.com.br) para obter a cidade e usa a API de clima (OpenWeather) para buscar as temperaturas em Celsius, Fahrenheit e Kelvin.
   - Ambos enviam **traces para o OpenTelemetry Collector**, que exporta os dados para o **Zipkin UI**, permitindo visualizar a cadeia completa: `servico-a → servico-b → APIs externas`.

---

## 🧰 Pré-requisitos

   - Chave de API de clima configurada na variável `WEATHER_API_KEY` (definida no `docker-compose.yaml`)  
   - Acesso à internet para chamadas às APIs externas (ViaCEP e OpenWeather)   

---

## 🧩 Estrutura do Projeto (visão geral)

   ```
   servico-a/
   ├── main.go # Inicializa o servidor e o tracer
   ├── handler.go # Roteamento e lógica da API
   ├── tracer.go # Configuração OpenTelemetry
   ├── Dockerfile
   ├── go.mod
   └── go.sum

   servico-b/
   ├── main.go
   ├── handler.go
   ├── cep.go # Consulta API ViaCEP
   ├── weather.go # Consulta API de clima
   ├── tracer.go
   ├── Dockerfile
   ├── go.mod
   └── go.sum

   otel-collector/
   ├── config.yaml # Configuração de receivers/exporters
   └── Dockerfile

   docker-compose.yaml # Orquestra tudo (A + B + OTEL + Zipkin)
   Makefile # Facilita build e execução
   README.md
   ```

---

## 🚀 Como Executar com Docker

1. **Configure as variáveis de ambiente**  
   Defina as variáveis no `docker-compose.yaml` (não é necessário `.env` separado).

2. **Construa e suba os containers**
   ```bash
   docker-compose up --build -d
   ```
3. **Verifique se os serviços estão ativos**

   ```bash
   docker ps
   ```
   - Serviço A: http://localhost:8080
   - Serviço B: http://localhost:8081
   - Zipkin UI: http://localhost:9411

## 🧪 Testando os Serviços (curl / REST Client)

   Após iniciar o ambiente com make up ou docker-compose up -d, você pode testar das seguintes formas:

   🔹 Teste via VS Code REST Client / Postman

      1. Abra o arquivo requests.http na raiz.
      2. Clique em Send Request em cada bloco.
      3. Teste:
         - ✅ CEP válido → retorna cidade + temperaturas
         - ❌ CEP inválido (menos de 8 dígitos) → erro 422
         - ❌ CEP inexistente → erro 404


🔹 Teste via curl

      ```bash
      curl -X POST http://localhost:8080/cep \
           -H "Content-Type: application/json" \
           -d '{"cep":"01001000"}'
      ```
## 🔍 **Visualizando Traces no Zipkin**

   Após uma requisição bem-sucedida:

      1. Acesse http://localhost:9411
      2. Clique em “Run Query”
      3. Você deverá ver spans encadeados:

          ```bash
          servico-a.handleCEP → HTTP POST servico-b → servico-b.handleWeather → HTTP GET ViaCEP
          ```
      
   Isso indica que o tracing distribuído está funcionando corretamente.

## 🧠 **Troubleshooting**

   - invalid URL escape "%2F" → corrige-se garantindo que o endpoint OTEL não tenha barras duplas.
   - Traces não aparecem → verifique logs do otel-collector e se o zipkin está saudável (docker ps → status “healthy”).]

## **Licença**

   MIT — livre para uso, modificação e distribuição.