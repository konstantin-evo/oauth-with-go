package main

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"learn.oauth.client/data/model"
	"log"
	"net/http"
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

func getSessionValue(session *sessions.Session, key string) string {
	value := session.Values[key]
	if value != nil {
		if strValue, ok := value.(string); ok {
			return strValue // Value is a string, return it as is
		} else if byteSliceValue, ok := value.([]uint8); ok {
			return string(byteSliceValue) // Convert byte slice to string
		}
	}
	return ""
}

func getCookieValue(r *http.Request, cookieName string) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Printf("Cookie not found by name: %s\n", cookieName)
		return ""
	}
	return cookie.Value
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

func setCORSHeaders(w http.ResponseWriter, config *HandlerConfig) {
	w.Header().Set("Access-Control-Allow-Origin", config.AppVar.FrontendHost)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
