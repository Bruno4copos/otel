package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
)

var ErrWeatherAPIKeyMissing = errors.New("weather api key not configured")

// GetTemperatureByCity consulta WeatherAPI (current.json) para obter temperatura em Celsius.
// Requer a variável de ambiente WEATHER_API_KEY configurada.
// Gera um span OTEL chamado "Weather Lookup".
func GetTemperatureByCity(ctx context.Context, city string) (float64, error) {
	ctx, span := tracer.Start(ctx, "Weather Lookup")
	defer span.End()
	span.SetAttributes(attribute.String("city_query", city))

	key := os.Getenv("WEATHER_API_KEY")
	if key == "" {
		return 0, ErrWeatherAPIKeyMissing
	}

	// prepara URL (escape do city)
	q := url.QueryEscape(city)
	apiURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", key, q)

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		// tenta ler corpo para dar contexto ao erro
		var bodyBytes []byte
		bodyBytes, _ = io.ReadAll(res.Body)
		return 0, fmt.Errorf("weatherapi error status %d: %s", res.StatusCode, string(bodyBytes))
	}

	// estrutura simplificada para extrair temp_c
	var wres struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}

	if err := json.NewDecoder(res.Body).Decode(&wres); err != nil {
		return 0, err
	}

	span.SetAttributes(attribute.Float64("temp_c", wres.Current.TempC))
	return wres.Current.TempC, nil
}
