package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type CepRequest struct {
	CEP string `json:"cep"`
}
type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func main() {
	shutdown := initTracer()
	defer shutdown(context.Background())

	http.Handle("/weather", otelhttp.NewHandler(http.HandlerFunc(HandleWeatherByCEP), "/weather"))
	http.HandleFunc("/healthz", HealthHandler)

	log.Println("🌤️ servico-b listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
