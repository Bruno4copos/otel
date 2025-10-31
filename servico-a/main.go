package main

import (
	"log"
	"net/http"
)

func main() {
	// initialize tracer
	if err := initTracer(); err != nil {
		log.Printf("otel init error: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cep", cepHandler)
	mux.HandleFunc("/healthz", healthHandler)

	log.Println("servico-a listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
