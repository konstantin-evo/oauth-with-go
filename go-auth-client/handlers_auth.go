package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"io"
	"learn.oauth.client/data/model"
	"log"
	"net/http"
	"net/url"
)

const (
	SessionStateKey  = "SessionState"
	TokenResponseKey = "TokenResponse"
)

type HandlerConfig struct {
	AppVar *config
	Store  *sessions.CookieStore
}

func LoginHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	redirectURL := buildAuthURL(config.AppVar)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	session, _ := config.Store.Get(r, "session")
	delete(session.Values, SessionStateKey)
	delete(session.Values, TokenResponseKey)

	err := session.Save(r, w)
	if err != nil {
		log.Println("Error saving session:", err)
	}

	redirectURL := buildLogoutURL(config.AppVar)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func AuthCodeRedirectHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	authCode := r.URL.Query().Get("code")
	sessionState := r.URL.Query().Get("session_state")

	// Exchange auth code for token
	tokenBytes, err := exchangeAuthCodeForToken(authCode, config.AppVar)
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

	// Save session state in the session
	session, _ := config.Store.Get(r, "session")
	session.Values[SessionStateKey] = sessionState
	err = session.Save(r, w)
	if err != nil {
		log.Println("Error saving session:", err)
	}

	setCookies(w, tokenResponse, sessionState)
	http.Redirect(w, r, config.AppVar.FrontendHost, http.StatusSeeOther)
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	// Get the current Refresh Token from the session
	session, _ := config.Store.Get(r, "session")
	tokenResponse, err := getTokenResponseFromSession(session)
	if err != nil {
		log.Println("Error decoding token response:", err)
		http.Error(w, "Failed to get token response from session", http.StatusInternalServerError)
		return
	}

	// Create a POST request to refresh the token
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", config.AppVar.AppID)
	data.Set("client_secret", "1ANIYGdYJhdeMjXOn6qrSmMU9wiUkXQ2")
	data.Set("refresh_token", tokenResponse.RefreshToken)

	req, err := http.NewRequest("POST", config.AppVar.TokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Println("Error creating a new HTTP request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request to refresh the token
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending HTTP request:", err)
		http.Error(w, "Failed to send token refresh request", http.StatusInternalServerError)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		err := resp.Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
			http.Error(w, "Token refresh request failed", http.StatusInternalServerError)
			return
		}

		http.Error(w, fmt.Sprintf("Token refresh request returned status code %d. Response body: %s", resp.StatusCode, responseBody), resp.StatusCode)
		return
	}

	// Read and parse the updated token
	tokenResponseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, "Failed to read token refresh response", http.StatusInternalServerError)
		return
	}

	// Save the new token in the session
	session.Values[TokenResponseKey] = tokenResponseBytes
	err = session.Save(r, w)
	if err != nil {
		log.Println("Error saving session:", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func exchangeAuthCodeForToken(authCode string, appVar *config) ([]byte, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", appVar.AppID)
	data.Set("client_secret", "1ANIYGdYJhdeMjXOn6qrSmMU9wiUkXQ2")
	data.Set("code", authCode)
	data.Set("redirect_uri", appVar.AuthCodeCallback)

	req, err := http.NewRequest("POST", appVar.TokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		err := resp.Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
			return nil, err
		}

		return nil, fmt.Errorf("token request returned status code %d. Response body: %s", resp.StatusCode, responseBody)
	}

	// Read the response body into a byte slice
	tokenResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return tokenResponse, nil
}

func buildAuthURL(appVar *config) string {
	u, err := url.Parse(appVar.AuthURL)
	if err != nil {
		log.Println(err)
		return ""
	}

	qs := u.Query()
	qs.Add("state", "test_state")
	qs.Add("client_id", appVar.AppID)
	qs.Add("response_type", "code")
	qs.Add("redirect_uri", appVar.AuthCodeCallback)
	u.RawQuery = qs.Encode()

	return u.String()
}

func buildLogoutURL(appVar *config) string {
	u, err := url.Parse(appVar.LogoutURL)
	if err != nil {
		log.Println(err)
		return ""
	}

	q := u.Query()
	q.Add("redirect_uri", appVar.LogoutRedirect)
	u.RawQuery = q.Encode()

	return u.String()
}

func getTokenResponseFromSession(session *sessions.Session) (*model.TokenResponseData, error) {
	tokenResponseStr := getSessionValue(session, TokenResponseKey)

	if tokenResponseStr == "" {
		return nil, fmt.Errorf("token response not found in session")
	}

	tokenResponseBytes := []byte(tokenResponseStr)

	var tokenResponse model.TokenResponseData
	err := json.Unmarshal(tokenResponseBytes, &tokenResponse)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}
