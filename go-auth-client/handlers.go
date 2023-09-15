package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"learn.oauth.client/model"
)

const (
	SessionStateKey  = "SessionState"
	TokenResponseKey = "TokenResponse"
)

type HandlerConfig struct {
	AppVar   *config
	Store    *sessions.CookieStore
	Template *template.Template
}

func homeHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	// Retrieve SessionState and token response from the session
	// to display on front-end
	session, _ := config.Store.Get(r, "session-name")
	sessionState := getSessionValue(session, SessionStateKey)
	tokenResponse, err := getTokenResponseFromSession(session)
	if err != nil {
		// Create an empty token response object to avoid nil
		log.Println("Error decoding token response:", err)
		tokenResponse = &model.TokenResponseData{}
	}

	// Decode access token (JWT)
	decodedToken, err := decodeAccessToken(tokenResponse.AccessToken)
	if err != nil {
		log.Println("Error decoding access token:", err)
		// Handle the error and send an error message to the front-end
	}

	data := model.FrontData{
		SessionState: sessionState,
		Token:        tokenResponseToMap(*tokenResponse),
		DecodedToken: decodedToken,
	}

	err = config.Template.Execute(w, data)
	if err != nil {
		log.Println("Template execution error:", err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	redirectURL := buildAuthURL(config.AppVar)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	session, _ := config.Store.Get(r, "session-name")
	delete(session.Values, SessionStateKey)
	delete(session.Values, TokenResponseKey)

	err := session.Save(r, w)
	if err != nil {
		log.Println("Error saving session:", err)
	}

	redirectURL := buildLogoutURL(config.AppVar)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func authCodeRedirectHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	authCode := r.URL.Query().Get("code")
	sessionState := r.URL.Query().Get("session_state")

	// Exchange auth code for token
	tokenBytes, err := exchangeAuthCodeForToken(authCode, config.AppVar)
	if err != nil {
		log.Println("Error exchanging auth code for token:", err)
		http.Error(w, "Failed to exchange authorization code for token", http.StatusInternalServerError)
		return
	}

	// Save token bytes and session state in the session
	session, _ := config.Store.Get(r, "session-name")
	session.Values[SessionStateKey] = sessionState
	session.Values[TokenResponseKey] = tokenBytes
	err = session.Save(r, w)
	if err != nil {
		log.Println("Error saving session:", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func servicesHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	// Create a request to the protected resource endpoint
	req, err := http.NewRequest("GET", config.AppVar.ServicesURL, nil)
	if err != nil {
		log.Println("Error creating a new HTTP request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session, _ := config.Store.Get(r, "session-name")
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

func refreshTokenHandler(w http.ResponseWriter, r *http.Request, config *HandlerConfig) {
	// Get the current Refresh Token from the session
	session, _ := config.Store.Get(r, "session-name")
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

func tokenResponseToMap(response model.TokenResponseData) map[string]interface{} {
	data := make(map[string]interface{})
	data["AccessToken"] = response.AccessToken
	data["TokenType"] = response.TokenType
	data["ExpiresIn"] = response.ExpiresIn
	data["RefreshToken"] = response.RefreshToken
	data["Scope"] = response.Scope
	return data
}

func decodeAccessToken(accessToken string) (map[string]interface{}, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
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
