package main

import (
	"context"
	"log"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// initTracer inicializa um TracerProvider com um exporter OTLP/HTTP.
// Retorna o TracerProvider (para que o chamador possa chamar Shutdown) e um error caso falhe.
//
// - Lê OTEL_EXPORTER_OTLP_ENDPOINT (p.ex. "http://otel-collector:4318/v1/traces").
// - Se não definido, usa defaultOTLP.
// - Suporta endpoints com ou sem esquema (se passar um URL com esquema, usa WithURL).
func initTracer() func(context.Context) error {
	ctx := context.Background()

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://otel-collector:4318"
	}

	var exporter *otlptrace.Exporter
	var err error

	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		// ✅ Se tiver esquema, usa WithEndpointURL corretamente
		exporter, err = otlptracehttp.New(
			ctx,
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpointURL(endpoint+"/v1/traces"),
		)
	} else {
		// ✅ Caso contrário, só host:porta
		exporter, err = otlptracehttp.New(
			ctx,
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithURLPath("/v1/traces"),
		)
	}

	if err != nil {
		log.Fatalf("❌ failed to create OTLP exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.Default()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	log.Printf("✅ OpenTelemetry tracer initialized (endpoint=%s)", endpoint)
	return tp.Shutdown
}
