package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrInvalidCEP  = errors.New("invalid cep format")
	ErrCEPNotFound = errors.New("cep not found")
)

var cepRegexp = regexp.MustCompile(`^[0-9]{8}$`)

// GetCityByCEP valida o CEP e consulta a API ViaCEP para retornar a cidade.
// Retorna ErrInvalidCEP quando o formato estiver incorreto.
// Retorna ErrCEPNotFound quando o viacep indicar que não existe.
// Gera um span OTEL chamado "ViaCEP Lookup".
func GetCityByCEP(ctx context.Context, cep string) (string, error) {
	// valida formato
	if !cepRegexp.MatchString(cep) {
		return "", ErrInvalidCEP
	}

	// cria span
	ctx, span := tracer.Start(ctx, "ViaCEP Lookup")
	defer span.End()
	span.SetAttributes(attribute.String("cep", cep))

	// chamada ViaCEP
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// ler corpo
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// parsing genérico — ViaCEP devolve {"erro": true} quando não encontrado
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	if _, hasErro := data["erro"]; hasErro {
		return "", ErrCEPNotFound
	}

	// campo de interesse é "localidade"
	if loc, ok := data["localidade"].(string); ok && loc != "" {
		span.SetAttributes(attribute.String("city", loc))
		return loc, nil
	}
	// fallback para "cidade" caso algo mude
	if loc, ok := data["cidade"].(string); ok && loc != "" {
		span.SetAttributes(attribute.String("city", loc))
		return loc, nil
	}

	// se não encontrou cidade, considera não encontrado
	return "", ErrCEPNotFound
}
