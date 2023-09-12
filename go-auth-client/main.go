package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
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
	ServicesURL      string
	WebPort          string
}

func main() {

	// Load configuration from environment variables or command-line arguments
	app, err := loadConfig()
	if err != nil {
		log.Panic(err)
	}

	// Load handler config
	handlerConfig, err := loadHandlerConfig(app)
	if err != nil {
		log.Panic(err)
	}

	// Start the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.WebPort),
		Handler: routes(handlerConfig),
	}

	log.Printf("Server listening on port %s\n", app.WebPort)
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
		ServicesURL:      "http://localhost:8082/billing/v1/services",
		WebPort:          port,
	}

	return &app, nil
}

func loadHandlerConfig(app *config) (*HandlerConfig, error) {
	store := sessions.NewCookieStore([]byte("your-secret-key"))

	// Загрузка шаблона из файла
	t := template.Must(template.ParseFiles("src/template/index.html"))

	return &HandlerConfig{
		AppVar:   app,
		Store:    store,
		Template: t,
	}, nil
}
