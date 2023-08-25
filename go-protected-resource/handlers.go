package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Billing struct {
	Services []string `json:"services"`
}

type BillingError struct {
	Error string `json:"error"`
}

type TokenIntrospect struct {
	Exp    int    `json:"exp"`
	Nbf    int    `json:"nbf"`
	Iat    int    `json:"iat"`
	Jti    string `json:"jti"`
	Aud    string `json:"aud"`
	Typ    string `json:"typ"`
	Acr    string `json:"acr"`
	Active bool   `json:"active"`
}

func servicesHandler(w http.ResponseWriter, r *http.Request, app *config) {
	token, err := extractToken(r)
	if err != nil {
		log.Println("Error extracting token:", err)
		s := &BillingError{Error: err.Error()}
		encoder := json.NewEncoder(w)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // Установите статус код 400
		encoder.Encode(s)
		return
	}

	if !validateToken(token, app) {
		s := &BillingError{Error: "Invalid token"}
		encoder := json.NewEncoder(w)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized) // Установите статус код 401
		encoder.Encode(s)
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
		// Check if the Authorization header starts with "Bearer "
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer "), nil
		}
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

func validateToken(token string, app *config) bool {

	form := url.Values{}
	form.Add("token", token)
	form.Add("token_type_hint", "requesting_party_token")

	req, err := http.NewRequest("POST", app.tokenIntrospectURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("Error creating a new HTTP request:", err)
		return false
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(app.clientIntrospect, app.clientSecret)

	// Send the request and handle the response
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Println("Error sending HTTP request:", err)
		return false
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println("Status code is not 200:", err)
		return false
	}

	// Read and parse the response
	byteBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return false
	}

	tokenIntrospect := &TokenIntrospect{}
	err = json.Unmarshal(byteBody, tokenIntrospect)
	if err != nil {
		log.Println("Error unmarshalling JSON response:", err)
		return false
	}

	return tokenIntrospect.Active
}
