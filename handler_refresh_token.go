package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/thomassifflet/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	strUserID, err := auth.ValidateJWT(refreshToken, cfg.jwtSecret, "chirpy-refresh")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	userID, convErr := strconv.Atoi(strUserID)
	if convErr != nil {
		respondWithError(w, http.StatusInternalServerError, "user id conversion error")
		return
	}

	newAccessToken, err := auth.MakeJWT(userID, cfg.jwtSecret, "chirpy-access", time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create jwt")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newAccessToken,
	})

}
