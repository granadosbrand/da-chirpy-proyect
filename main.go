package main

import (
	"database/sql"
	"log"
	"net/http"

	"os"
	"sync/atomic"

	"github.com/granadosbrand/da-chirpy-proyect/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	FileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	polka_key      string
}

func main() {

	const port = "8080"
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretKey := os.Getenv("SECRET")
	polka_key := os.Getenv("POLKA_KEY")

	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error opening database %s", err)
	}

	dbQueries := database.New(dbConn)

	apiCfg := &apiConfig{
		FileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		secret:         secretKey,
		polka_key:      polka_key,
	}

	mux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// Users
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgradeUser)

	// Auth
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)

	// mux.HandleFunc("POST /api/validate_chirp", )

	// Chirpy CRUD

	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.getChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirp)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
