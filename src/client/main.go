package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type config struct {
	AppID            string
	AuthURL          string
	TokenURL         string
	LogoutURL        string
	LogoutRedirect   string
	AuthCodeCallback string
}

func main() {

	app := config{
		AppID:            "billingApp",
		AuthURL:          "http://localhost:8081/realms/customRealm/protocol/openid-connect/auth",
		LogoutURL:        "http://localhost:8081/realms/customRealm/protocol/openid-connect/logout",
		TokenURL:         "http://localhost:8081/realms/customRealm/protocol/openid-connect/token",
		LogoutRedirect:   "http://localhost:8080/",
		AuthCodeCallback: "http://localhost:8080/authCodeRedirect",
	}

	r := mux.NewRouter()

	// Create a file server to serve static files from the "src" directory,
	// Specify a prefix for static files and attach the file server to the route for serving static files
	staticFileServer := http.FileServer(http.Dir("src"))
	staticFilesRoute := "/src/"
	r.PathPrefix(staticFilesRoute).Handler(http.StripPrefix(staticFilesRoute, staticFileServer))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(w, r)
	})
	r.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		tokenHandler(w, r, &app)
	})
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(w, r, &app)
	})
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logoutHandler(w, r, &app)
	})
	r.HandleFunc("/authCodeRedirect", func(w http.ResponseWriter, r *http.Request) {
		authCodeRedirectHandler(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(addr, r))
}
