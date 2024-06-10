package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AlvaroPrates/Chirpy/internal/auth"
	"github.com/AlvaroPrates/Chirpy/internal/database"
)

func (cfg *apiConfig) handleChirpDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user id")
	}

	chirpID, err := strconv.Atoi(r.PathValue("chirp_id"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse chirp id")
	}

	chirp, err := cfg.DB.GetChirpByID(chirpID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Chirp does not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp by id")
		return
	}

	if chirp.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "Cannot delete unauthored chirp")
		return
	}

	if err = cfg.DB.DeleteChirp(chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
