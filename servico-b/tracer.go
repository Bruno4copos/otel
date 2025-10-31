package main

import (
	"context"
	"log"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	otlptrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	otlptracehttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const defaultOTLP = "http://otel-collector:4318/v1/traces"

// InitTracer inicializa um TracerProvider com um exporter OTLP/HTTP.
// Retorna o TracerProvider (para que o chamador possa chamar Shutdown) e um error caso falhe.
//
// - Lê OTEL_EXPORTER_OTLP_ENDPOINT (p.ex. "http://otel-collector:4318/v1/traces").
// - Se não definido, usa defaultOTLP.
// - Suporta endpoints com ou sem esquema (se passar um URL com esquema, usa WithURL).
func InitTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = defaultOTLP
	}

	var exporter *otlptrace.Exporter
	var err error

	// Se endpoint contém esquema (http:// ou https://) usamos WithURL,
	// caso contrário usamos WithEndpoint (host:port).
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithURLPath(endpoint),
			// Se estiver usando http (não https), WithInsecure é necessário.
			// Com URL https://... a opção WithInsecure é ignorada.
			otlptracehttp.WithInsecure(),
		)
	} else {
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithInsecure(),
		)
	}
	if err != nil {
		log.Printf("failed creating OTLP HTTP exporter: %v", err)
		return nil, err
	}

	// Resource (pode adicionar atributos se quiser)
	res, err := resource.New(ctx)
	if err != nil {
		log.Printf("failed creating resource: %v", err)
		// continue mesmo assim (não crítico)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
