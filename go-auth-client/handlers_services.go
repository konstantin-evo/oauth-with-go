package main

import (
	"context"
	"encoding/json"
	"io"
	"learn.oauth.client/data/model"
	"log"
	"net/http"
	"time"
)

func ServicesHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	// Create a request to the protected resource endpoint
	req, err := http.NewRequest("GET", config.AppVar.ServicesURL, nil)
	if err != nil {
		log.Println("Error creating a new HTTP request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session, _ := config.Store.Get(r, "session")
	tokenResponse, err := getTokenResponseFromSession(session)
	if err != nil {
		log.Println("Error decoding token response:", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+tokenResponse.AccessToken)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelFunc()

	// Send the request and handle the response
	c := &http.Client{}
	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		log.Println("Error sending HTTP request:", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(res.Body)

	// Read and parse the response
	byteBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		errorResponse := &model.BillingError{}

		err = json.Unmarshal(byteBody, errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.StatusCode)
		err := json.NewEncoder(w).Encode(errorResponse)
		if err != nil {
			log.Println("Error encoding JSON error response:", err)
		}
		return
	}

	billingResponse := &model.Billing{}
	err = json.Unmarshal(byteBody, billingResponse)
	if err != nil {
		log.Println("Error unmarshalling JSON response:", err)
		return
	}

	// Prepare the JSON response with only the services
	jsonResponse := map[string]interface{}{
		"services": billingResponse.Services,
	}

	// Marshal the JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(jsonResponse)
	if err != nil {
		log.Println("Error encoding JSON response:", err)
		return
	}
}
