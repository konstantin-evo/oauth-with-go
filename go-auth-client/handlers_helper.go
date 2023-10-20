package main

import (
	"encoding/json"
	"learn.oauth.client/data/model"
	"log"
	"net/http"
	"time"
)

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		log.Println("Error writing JSON response:", err)
	}
}

func setCookies(w http.ResponseWriter, tokenResponse model.TokenResponseData, session string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: tokenResponse.AccessToken,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: session,
	})
}

func clearCookie(w http.ResponseWriter, cookieName string) {
	expiration := time.Now().AddDate(0, 0, -1)
	deletedCookie := http.Cookie{
		Name:    cookieName,
		Value:   "",
		Expires: expiration,
		MaxAge:  -1,
		Path:    "/",
	}
	http.SetCookie(w, &deletedCookie)
}

func setCORSHeaders(w http.ResponseWriter, config *HandlerConfig) {
	w.Header().Set("Access-Control-Allow-Origin", config.AppVar.FrontendHost)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
