package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/thomassifflet/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no authorization to delete chirp : invalid JWT token")
		return
	}
	userID, err := auth.ValidateJWT(userToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not authenticate : cannot validate jwt token")
		return
	}

	userIDConv, err := strconv.Atoi(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not convert user id into integer")
		return
	}
	if dbChirp.AuthorID != userIDConv {
		respondWithError(w, http.StatusForbidden, "cannot delete chip : wrong chirp author")
		return
	}

	fmt.Printf("Chirp author ID : %v, UserID : %v", dbChirp.AuthorID, userID)
	deletedChirp, err := cfg.DB.DeleteChirp(chirpID, dbChirp.AuthorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete chirp")
		return
	}
	fmt.Printf("%v", deletedChirp)

	respondWithJSON(w, http.StatusOK, "{}")
}
