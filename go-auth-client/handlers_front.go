package main

import (
	"encoding/json"
	"learn.oauth.client/data/model"
	"log"
	"net/http"
	"strings"
)

func GetTokenDataHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	// Add CORS headers to allow access from specific origins
	setCORSHeaders(w, config)

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

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	// Add CORS headers to allow access from specific origins
	setCORSHeaders(w, config)

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

	tokenBytes, err := SendRefreshTokenRequest(w, config, tokenData.RefreshToken)
	if err != nil {
		log.Println("Error exchanging auth code for token:", err)
		http.Error(w, "Failed to exchange authorization code for token", http.StatusInternalServerError)
		return
	}

	// Save token to DB
	var tokenResponse model.TokenResponseData
	err = json.Unmarshal(tokenBytes, &tokenResponse)
	if err != nil {
		log.Println("Error decoding token response:", err)
		http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
		return
	}

	_, err = config.AppVar.Repo.Insert(tokenResponse)
	if err != nil {
		log.Println("Error saving token to the database:", err)
		http.Error(w, "Failed to save token to the database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tokenResponse)
	if err != nil {
		log.Println("Error encoding JSON error response:", err)
		return
	}
}
