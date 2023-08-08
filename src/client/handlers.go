package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
)

var t = template.Must(template.ParseFiles("template/index.html"))
var store = sessions.NewCookieStore([]byte("your-secret-key"))

type authSession struct {
	AuthCode     string
	SessionState string
}

// Handle the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, _ := store.Get(r, "session-name")

	// Retrieve AuthCode and SessionState from session
	authCodeValue := session.Values["AuthCode"]
	sessionStateValue := session.Values["SessionState"]

	// Check if authCodeValue and sessionStateValue are not nil before type asserting
	var authCode, sessionState string
	if authCodeValue != nil {
		authCode = authCodeValue.(string)
	}
	if sessionStateValue != nil {
		sessionState = sessionStateValue.(string)
	}

	// Populate the data structure for template rendering
	data := authSession{
		AuthCode:     authCode,
		SessionState: sessionState,
	}

	// Execute the template with the data
	err := t.Execute(w, data)
	if err != nil {
		log.Println("Template execution error:", err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request, appVar *config) {
	redirectURL := buildAuthURL(appVar)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// Handle the logout process
func logoutHandler(w http.ResponseWriter, r *http.Request, appVar *config) {
	// Get the session
	session, _ := store.Get(r, "session-name")

	// Clear AuthCode and SessionState from the session
	delete(session.Values, "AuthCode")
	delete(session.Values, "SessionState")

	// Save the session
	err := session.Save(r, w)
	if err != nil {
		log.Println("Error saving session:", err)
	}

	// Redirect to the logout URL
	redirectURL := buildLogoutURL(appVar)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// Handle the callback after authorization code is received
func authCodeRedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, _ := store.Get(r, "session-name")

	// Store AuthCode and SessionState in the session
	session.Values["AuthCode"] = r.URL.Query().Get("code")
	session.Values["SessionState"] = r.URL.Query().Get("session_state")
	session.Save(r, w)

	// Clear the query parameters and redirect to the home page
	r.URL.RawQuery = ""
	log.Printf("Request queries: %+v", session.Values)
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
