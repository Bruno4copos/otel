package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer = otel.Tracer("servico-b/handler")

// request representa o corpo recebido via POST (enviado pelo Service A)
type request struct {
	CEP string `json:"cep"`
}

// response representa o corpo retornado pelo Service B
type response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// HandleWeatherByCEP é o handler principal do Service B
func HandleWeatherByCEP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "HandleWeatherByCEP")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// valida formato do CEP
	if len(req.CEP) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	span.SetAttributes(attribute.String("cep", req.CEP))

	// busca cidade via ViaCEP (modulo cep.go)
	city, err := GetCityByCEP(ctx, req.CEP)
	if err != nil {
		if err == ErrInvalidCEP {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		if err == ErrCEPNotFound {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("internal error: %v", err), http.StatusInternalServerError)
		return
	}

	// busca temperatura via WeatherAPI (modulo weather.go)
	tempC, err := GetTemperatureByCity(ctx, city)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching weather: %v", err), http.StatusInternalServerError)
		return
	}

	resp := response{
		City:  city,
		TempC: tempC,
		TempF: CelsiusToFahrenheit(tempC),
		TempK: CelsiusToKelvin(tempC),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CelsiusToFahrenheit converte °C → °F
func CelsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

// CelsiusToKelvin converte °C → K
func CelsiusToKelvin(c float64) float64 {
	return c + 273
}

// HealthHandler apenas responde se o serviço está OK
func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// RegisterHandlers registra todas as rotas HTTP do serviço
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/weather", HandleWeatherByCEP)
	mux.HandleFunc("/healthz", HealthHandler)
}
