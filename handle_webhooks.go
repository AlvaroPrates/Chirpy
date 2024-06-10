package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AlvaroPrates/Chirpy/internal/database"
)

func (cfg *apiConfig) handleWebhookPolka(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err = cfg.DB.UpgradeUser(params.Data.UserID); err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "User does not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
