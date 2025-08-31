package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/edipretoro/boot.dev/go_web_server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var apiCfg *apiConfig

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Problem when connecting to the database: %v", err)
	}
	dbQueries := database.New(db)
	mux := http.NewServeMux()
	apiCfg = &apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
	}

	// /app/ part of this app
	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))),
	)

	// /admin/ part of this app
	mux.HandleFunc("GET /admin/metrics", apiCfg.middlewareReturnMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.middlewareResetMetrics)

	// /api/ part of this app
	mux.HandleFunc("GET /api/healthz", apiHealthz)
	// mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("POST /api/users", addUser)
	mux.HandleFunc("GET /api/chirps", getAllChirps)
	mux.HandleFunc("POST /api/chirps", addChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", getChirpByID)
	mux.HandleFunc("POST /api/login", loginUser)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
