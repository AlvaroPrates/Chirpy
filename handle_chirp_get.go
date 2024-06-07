package main

import (
	"errors"
	"net/http"
	"sort"
	"strconv"

	"github.com/AlvaroPrates/Chirpy/internal/database"
)

func (cfg *apiConfig) handleRetrieveChirp(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handleRetrieveChirpByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("chirp_id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse id")
	}

	dbChirp, err := cfg.DB.GetChirpByID(id)
	if err != nil {
		if errors.Is(err, database.ErrChirpNotExists) {
			respondWithError(w, http.StatusNotFound, "Chirp does not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp by id")
		return
	}

	chirp := Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
