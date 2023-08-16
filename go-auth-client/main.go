package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct {
	AppID            string
	AuthURL          string
	TokenURL         string
	LogoutURL        string
	LogoutRedirect   string
	AuthCodeCallback string
	WebPort          string
}

func main() {

	// Load configuration from environment variables or command-line arguments
	app, err := loadConfig()
	if err != nil {
		log.Panic(err)
	}

	// Start the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.WebPort),
		Handler: routes(app),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func loadConfig() (*config, error) {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := config{
		AppID:            "billingApp",
		AuthURL:          "http://localhost:8081/realms/customRealm/protocol/openid-connect/auth",
		LogoutURL:        "http://localhost:8081/realms/customRealm/protocol/openid-connect/logout",
		TokenURL:         "http://localhost:8081/realms/customRealm/protocol/openid-connect/token",
		LogoutRedirect:   "http://localhost:8080/",
		AuthCodeCallback: "http://localhost:8080/authCodeRedirect",
		WebPort:          port,
	}

	return &app, nil
}
