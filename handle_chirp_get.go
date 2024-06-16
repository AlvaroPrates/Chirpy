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

	authorID := -1
	authorIDStr := r.URL.Query().Get("author_id")
	if authorIDStr != "" {
		authorID, err = strconv.Atoi(authorIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
	}

	sortDirection := "asc"
	sortDirectionStr := r.URL.Query().Get("sort")
	if sortDirectionStr == "desc" {
		sortDirection = sortDirectionStr
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if authorID != -1 && dbChirp.AuthorID != authorID {
			continue
		}

		chirps = append(chirps, Chirp{
			ID:       dbChirp.ID,
			AuthorID: dbChirp.AuthorID,
			Body:     dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortDirection == "desc" {
			return chirps[i].ID > chirps[j].ID
		}
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
		if errors.Is(err, database.ErrNotExist) {
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
