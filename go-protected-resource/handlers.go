package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Billing struct {
	Services []string `json:"services"`
}

type BillingError struct {
	Error string `json:"error"`
}

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	token, err := extractToken(r)
	if err != nil {
		log.Println("Error extracting token:", err)
		s := &BillingError{Error: err.Error()}
		encoder := json.NewEncoder(w)
		encoder.Encode(s)
		w.Header().Add("Content-Type", "application/json")
		return
	}

	if !validateToken(token) {
		s := &BillingError{Error: "Invalid token"}
		encoder := json.NewEncoder(w)
		encoder.Encode(s)
		w.Header().Add("Content-Type", "application/json")
		return
	}

	s := Billing{
		Services: []string{
			"electric",
			"phone",
			"internet",
			"water",
		},
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(s)
	w.Header().Add("Content-Type", "application/json")
}

func extractToken(r *http.Request) (string, error) {
	// Try to extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		return authHeader, nil
	}

	// Try to extract token from form data
	formToken := r.FormValue("access_token")
	if formToken != "" {
		return formToken, nil
	}

	// Try to extract token from query parameters
	queryToken := r.URL.Query().Get("access_token")
	if queryToken != "" {
		return queryToken, nil
	}

	return "", errors.New("token not found")
}

func validateToken(token string) bool {
	return false
}
