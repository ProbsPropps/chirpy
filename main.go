package main

import (
	"net/http"
	"sync/atomic"
)

func main() {
	var apiCfg apiConfig
	mu := http.NewServeMux()
	mu.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mu.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mu.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mu.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mu.HandleFunc("POST /api/validate_chirp", handlerValidate)
	server := http.Server{
		Handler: mu,
		Addr: ":8080",
	}
	server.ListenAndServe()
}

type apiConfig struct {
	fileserverHits atomic.Int32
}
