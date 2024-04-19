package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func chirpHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "couldn't read request")
		return
	}
	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err == nil && len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	}
	if err != nil {
		respondWithError(w, 500, "couldn't unmarshal parameters")
		return
	}

	respondWithJSON(w, 200, returnVals{
		Valid: true,
	})
}
