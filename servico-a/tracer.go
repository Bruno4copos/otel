package main

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	otlptrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var tp *sdktrace.TracerProvider

func initTracer() error {
	ctx := context.Background()
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://otel-collector:4318/v1/traces"
	}

	exporter, err := otlptrace.New(ctx, otlptrace.WithInsecure(), otlptrace.WithEndpoint("otel-collector:4318"), otlptrace.WithURLPath(endpoint))
	if err != nil {
		log.Printf("failed to create OTLP exporter: %v", err)
		return err
	}

	res, _ := resource.New(ctx)
	tp = sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))
	tpCtx := tp
	otel.SetTracerProvider(tpCtx)
	return nil
}

func shutdownTracer(ctx context.Context) error {
	if tp == nil {
		return nil
	}
	return tp.Shutdown(ctx)
}
