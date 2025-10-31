package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
)

func main() {
	// inicializa OTEL/Zipkin tracer
	tp, err := InitTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("error shutting down tracer provider: %v", err)
		}
	}()

	// cria mux e registra rotas
	mux := http.NewServeMux()
	RegisterHandlers(mux)

	// cria servidor HTTP
	srv := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// canal para encerrar servidor com Ctrl+C
	go func() {
		log.Println("Service B (weather orchestrator) listening on :8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v\n", err)
		}
	}()

	// aguarda sinal de término (SIGINT/SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down servico-b...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
	log.Println("servico-b gracefully stopped")
}

// processHandler é um exemplo de endpoint simples que
// encaminha para o handler principal HandleWeatherByCEP.
func processHandler(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("servico-b/main")
	ctx, span := tr.Start(r.Context(), "processHandler")
	defer span.End()

	// Apenas delega ao handler principal (para manter compatibilidade)
	HandleWeatherByCEP(w, r.WithContext(ctx))
}
