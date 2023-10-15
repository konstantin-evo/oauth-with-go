package main

import (
	"net/http"
	"strings"
)

func GetTokenData(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	// Add CORS headers to allow access from specific origins
	w.Header().Set("Access-Control-Allow-Origin", config.AppVar.FrontendHost)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// If this is a preflight request, send an empty response
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}

	accessToken := strings.TrimPrefix(authorizationHeader, "Bearer ")
	tokenData, err := config.AppVar.Repo.GetByAccessToken(accessToken)
	if err != nil {
		http.Error(w, "Error fetching token data", http.StatusInternalServerError)
		return
	}

	if tokenData == nil {
		http.Error(w, "Token not found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, tokenData)
}
