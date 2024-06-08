package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AlvaroPrates/Chirpy/internal/auth"
)

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type returningVals struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	expiresAt := time.Now().Add(time.Duration(params.ExpiresInSeconds))

	jwt, err := auth.CreateJWT(cfg.jwtSecret, expiresAt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create jwt")
	}

	respondWithJSON(w, http.StatusOK, returningVals{
		ID:    user.ID,
		Email: user.Email,
		Token: jwt,
	})
}
