package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type AppVar struct {
	AuthCode     string
	SessionState string
}

var oauth = struct {
	authURL   string
	logoutURL string
}{
	authURL:   "http://localhost:8081/realms/customRealm/protocol/openid-connect/auth",
	logoutURL: "http://localhost:8081/realms/customRealm/protocol/openid-connect/logout",
}

var t = template.Must(template.ParseFiles("template/index.html"))
var appVar = AppVar{}

func main() {
	fmt.Println("Hello!")
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/authCodeRedirect", authCodeRedirect)
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	t.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, oauth.authURL, nil)
	if err != nil {
		log.Println(err)
		return
	}

	qs := url.Values{}
	qs.Add("state", "test_state")
	qs.Add("client_id", "billingApp")
	qs.Add("response_type", "code")
	qs.Add("redirect_uri", "http://localhost:8080/authCodeRedirect")

	req.URL.RawQuery = qs.Encode()
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	q := url.Values{}
	q.Add("redirect_uri", "http://localhost:8080")

	logoutURL, err := url.Parse(oauth.logoutURL)
	if err != nil {
		log.Println(err)
		return
	}

	logoutURL.RawQuery = q.Encode()
	http.Redirect(w, r, logoutURL.String(), http.StatusFound)
}

func authCodeRedirect(w http.ResponseWriter, r *http.Request) {
	appVar.AuthCode = r.URL.Query().Get("code")
	appVar.SessionState = r.URL.Query().Get("session_state")
	r.URL.RawQuery = ""
	fmt.Printf("Request queries: %+v", appVar)
	http.Redirect(w, r, "http://localhost:8080", http.StatusAccepted)
}
