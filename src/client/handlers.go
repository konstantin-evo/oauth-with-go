package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var t = template.Must(template.ParseFiles("template/index.html"))

func homeHandler(w http.ResponseWriter, r *http.Request, appVar *config) {
	t.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request, appVar *config) {
	redirectURL := buildAuthURL(appVar)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request, appVar *config) {
	redirectURL := buildLogoutURL(appVar)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func authCodeRedirectHandler(w http.ResponseWriter, r *http.Request, authSession *authSession) {
	authSession.AuthCode = r.URL.Query().Get("code")
	authSession.SessionState = r.URL.Query().Get("session_state")
	r.URL.RawQuery = ""
	log.Printf("Request queries: %+v", authSession)
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
