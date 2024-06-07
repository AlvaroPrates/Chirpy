package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/AlvaroPrates/Chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		if err := os.Remove("database.json"); err != nil {
			log.Printf("Failed to %s", err)
		}
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("failed to create database: %s", err)
	}

	apiCfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.resetMetrics)

	mux.HandleFunc("POST /api/login", apiCfg.handleUserLogin)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleRetrieveChirp)
	mux.HandleFunc("GET /api/chirps/{chirp_id}", apiCfg.handleRetrieveChirpByID)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
