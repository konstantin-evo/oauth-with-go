package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct {
	tokenIntrospectURL string
	clientIntrospect   string
	clientSecret       string
	webPort            string
}

func main() {

	// Load configuration from environment variables or command-line arguments
	app, err := loadConfig()
	if err != nil {
		log.Panic(err)
	}

	// Start the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.webPort),
		Handler: routes(app),
	}

	log.Printf("Server listening on port %s\n", app.webPort)
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func loadConfig() (*config, error) {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	app := config{
		tokenIntrospectURL: "http://localhost:8081/realms/customRealm/protocol/openid-connect/token/introspect",
		clientIntrospect:   "introspectClient",
		clientSecret:       "Q7pFiY2pmkRl3Eq6eBjvmKWcZJTwmZSo",
		webPort:            port,
	}

	return &app, nil
}
