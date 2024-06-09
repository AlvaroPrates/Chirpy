package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/AlvaroPrates/Chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env variables: ", err)
	}

	const filepathRoot = "."
	const port = "8080"

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		if err := os.Remove("database.json"); err != nil {
			log.Print("Failed to ", err)
		}
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal("failed to create database: ", err)
	}

	apiCfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.resetMetrics)

	mux.HandleFunc("POST /api/login", apiCfg.handleUserLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUserUpdate)
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
