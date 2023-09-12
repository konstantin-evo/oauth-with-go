package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"learn.oauth.billing/model"
)

// Billing struct represents the response structure for billing services.
type Billing struct {
	Services []string `json:"services"`
}

// BillingError struct represents an error response.
type BillingError struct {
	Error string `json:"error"`
}

// servicesHandler handles the request to retrieve billing services.
func servicesHandler(w http.ResponseWriter, r *http.Request, app *config) {
	token, err := extractToken(r)
	if err != nil {
		makeErrorMessage(w, err)
		return
	}

	if !validateToken(token, app) {
		makeErrorMessage(w, errors.New("invalid token"))
		return
	}

	claimBytes, err := getClaim(token)
	if err != nil {
		makeErrorMessage(w, err)
		return
	}

	tokenClaim := &model.TokenClaim{}
	err = json.Unmarshal(claimBytes, tokenClaim)
	if err != nil {
		makeErrorMessage(w, err)
		return
	}

	if !isValidAudience(tokenClaim.Aud) {
		makeErrorMessage(w, errors.New("invalid audience"))
		return
	}

	if !strings.Contains(tokenClaim.Scope, "getBillingService") {
		makeErrorMessage(w, errors.New("invalid scope. Required scope is getBillingService"))
		return
	}

	// Prepare the response data
	s := Billing{
		Services: []string{
			"electric",
			"phone",
			"internet",
			"water",
		},
	}

	// Encode and send the response
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	encoder.Encode(s)
}

// extractToken extracts the access token from various sources
func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	formToken := r.FormValue("access_token")
	if formToken != "" {
		return formToken, nil
	}

	queryToken := r.URL.Query().Get("access_token")
	if queryToken != "" {
		return queryToken, nil
	}

	return "", errors.New("token not found")
}

// validateToken sends an introspection request to validate the token.
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
	var tokenIntrospect model.TokenIntrospect
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&tokenIntrospect); err != nil {
		log.Println("Error decoding JSON response:", err)
		return false
	}

	// Check the type of Aud and handle it accordingly
	switch aud := tokenIntrospect.Aud.(type) {
	case string:
		// If Aud is a string, convert it to a single-element slice
		tokenIntrospect.Aud = []string{aud}
	case []interface{}:
		// If Aud is an array, convert it to a slice of strings
		audStrings := make([]string, len(aud))
		for i, v := range aud {
			if s, ok := v.(string); ok {
				audStrings[i] = s
			} else {
				log.Println("Error converting aud value to string:", v)
				return false
			}
		}
		tokenIntrospect.Aud = audStrings
	default:
		log.Println("Unexpected type for aud:", aud)
		return false
	}

	return tokenIntrospect.Active
}

// isValidAudience checks if the 'aud' field in the token claim is valid.
func isValidAudience(aud interface{}) bool {
	switch aud.(type) {
	case string:
		return aud.(string) == "http://localhost:8080/"
	case []interface{}:
		audiences := aud.([]interface{})
		for _, a := range audiences {
			if a == "http://localhost:8080/" {
				return true
			}
		}
	}
	return false
}

// getClaim decodes the JWT and returns its claims.
func getClaim(token string) ([]byte, error) {
	tokenParts := strings.Split(token, ".")
	claim, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		log.Println("Error decoding JWT:", err)
		return []byte{}, err
	}
	return claim, nil
}

// makeErrorMessage creates and sends an error response.
func makeErrorMessage(w http.ResponseWriter, err error) {
	log.Println(err)
	s := &BillingError{Error: err.Error()}
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	encoder.Encode(s)
}
