package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/edipretoro/boot.dev/go_web_server/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareReturnMetrics(w http.ResponseWriter, req *http.Request) {
	template := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(([]byte)(fmt.Sprintf(template, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareResetMetrics(w http.ResponseWriter, req *http.Request) {
	if os.Getenv("PLATFORM") == "dev" {
		cfg.dbQueries.DeleteAllUsers(req.Context())
	}
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(([]byte)("OK"))
}
