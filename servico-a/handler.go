package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer
var cepRe = regexp.MustCompile(`^[0-9]{8}$`)

func init() {
	tracer = otel.Tracer("servico-a")
}

type cepRequest struct {
	Cep string `json:"cep"`
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	var body cepRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	if !cepRe.MatchString(body.Cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// start span for forwarding
	ctx, span := tracer.Start(ctx, "ForwardToServiceB")
	defer span.End()

	payload, _ := json.Marshal(body)
	url := os.Getenv("SERVICE_B_URL")
	if url == "" {
		url = "http://servico-b:8081/process"
	}

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport), Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		log.Printf("error creating request: %v", err)
		http.Error(w, "error creating request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("error forwarding to servico-b: %v", err)
		http.Error(w, "error forwarding request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}
