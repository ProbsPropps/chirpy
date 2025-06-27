package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/ProbsPropps/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	
	dbQueries := database.New(db)

	apiCfg := apiConfig {
		fileserverHits: atomic.Int32{},
		queries: 		dbQueries,
		platform: platform,
	}

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
	mu.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	server := http.Server{
		Handler: 	mu,
		Addr: 		":8080",
	}
	server.ListenAndServe()
}

type apiConfig struct {
	fileserverHits 	atomic.Int32
	queries 		*database.Queries
	platform 		string
}
