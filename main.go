package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
	reset          (int)
}

func (cfg *apiConfig) middleWareMetricsIncrement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetMetrics() {
	cfg.fileserverHits = 0
}

func (cfg *apiConfig) handleResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.resetMetrics()
	hitsCount := fmt.Sprintf("Hits count has been reset")
	w.Header().Set("Content-Type", "text/plain charset=utf-8")
	w.Write([]byte(hitsCount))
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	hitsCount := fmt.Sprintf("Hits: %d", cfg.fileserverHits)
	w.Header().Set("Content-Type", "text/plain charset=utf-8")
	w.Write([]byte(hitsCount))

}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)
	apiCfg := apiConfig{}

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", apiCfg.middleWareMetricsIncrement(http.StripPrefix("/app/", fileServer)))
	mux.HandleFunc("/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("/reset", apiCfg.handleResetMetrics)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})
	fmt.Println("Started server on port 8080..")
	log.Fatal(http.ListenAndServe(":8080", corsMux))

}
