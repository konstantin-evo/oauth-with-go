package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	"learn.oauth.client/data/repository"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type config struct {
	AppID            string
	AppSecret        string
	AuthURL          string
	TokenURL         string
	LogoutURL        string
	LogoutRedirect   string
	AuthCodeCallback string
	ServicesURL      string
	FrontendHost     string
	WebPort          string
	WebHost          string
	Repo             repository.Repository
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

	host := removeHTTPPrefix(app.WebHost)
	// Start the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, app.WebPort),
		Handler: routes(handlerConfig),
	}

	log.Printf("Server listening on %s:%s\n", host, app.WebPort)
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func loadConfig() (*config, error) {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "http://localhost"
	}

	frontendHost := os.Getenv("FRONTEND_HOST")
	if frontendHost == "" {
		frontendHost = "http://localhost:3000"
	}

	protectedResourceHost := os.Getenv("PROTECTED_RESOURCE_HOST")
	if protectedResourceHost == "" {
		protectedResourceHost = "http://localhost:8082"
	}

	keycloakHost := os.Getenv("KEYCLOAK_HOST")
	if keycloakHost == "" {
		keycloakHost = "http://localhost:8080"
	}

	clientID := os.Getenv("CLIENT_ID")
	if clientID == "" {
		clientID = "billingApp"
	}

	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		clientSecret = "1ANIYGdYJhdeMjXOn6qrSmMU9wiUkXQ2"
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=password dbname=oauth sslmode=disable timezone=UTC connect_timeout=5"
	}

	conn := connectToDB(dsn)
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	app := config{
		AppID:            clientID,
		AppSecret:        clientSecret,
		AuthURL:          keycloakHost + "/realms/customRealm/protocol/openid-connect/auth",
		LogoutURL:        keycloakHost + "/realms/customRealm/protocol/openid-connect/logout",
		TokenURL:         keycloakHost + "/realms/customRealm/protocol/openid-connect/token",
		LogoutRedirect:   host + ":" + port + "/logoutRedirect",
		AuthCodeCallback: host + ":" + port + "/authCodeRedirect",
		ServicesURL:      protectedResourceHost + "/billing/v1/services",
		FrontendHost:     frontendHost,
		WebPort:          port,
		WebHost:          host,
		Repo:             repository.NewPostgresRepository(conn),
	}

	return &app, nil
}

func loadHandlerConfig(app *config) (*HandlerConfig, error) {
	store := sessions.NewCookieStore([]byte("your-secret-key"))

	return &HandlerConfig{
		AppVar: app,
		Store:  store,
	}, nil
}

func connectToDB(dsn string) *sql.DB {
	counts := 0

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Printf("Postgres not yet ready (attempt %d): %s", counts+1, err.Error())
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Printf("Giving up after %d attempts: %s", counts, err.Error())
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func removeHTTPPrefix(input string) string {
	if strings.HasPrefix(input, "http://") {
		return input[7:]
	}
	if strings.HasPrefix(input, "https://") {
		return input[8:]
	}
	return input
}
