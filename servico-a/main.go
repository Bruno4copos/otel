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

func main() {
	shutdown := initTracer()
	defer shutdown(context.Background())

	handler := otelhttp.NewHandler(http.HandlerFunc(handleCEP), "/cep")
	http.Handle("/cep", handler)

	log.Println("🚀 servico-a listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
