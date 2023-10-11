package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	"learn.oauth.client/data/repository"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type config struct {
	AppID            string
	AuthURL          string
	TokenURL         string
	LogoutURL        string
	LogoutRedirect   string
	AuthCodeCallback string
	ServicesURL      string
	FrontendHost     string
	WebPort          string
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

	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=password dbname=oauth sslmode=disable timezone=UTC connect_timeout=5"
	}

	frontendHost := os.Getenv("FRONTEND_HOST")
	if frontendHost == "" {
		frontendHost = "http://localhost:3000"
	}

	conn := connectToDB(dsn)
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	app := config{
		AppID:            "billingApp",
		AuthURL:          "http://localhost:8081/realms/customRealm/protocol/openid-connect/auth",
		LogoutURL:        "http://localhost:8081/realms/customRealm/protocol/openid-connect/logout",
		TokenURL:         "http://localhost:8081/realms/customRealm/protocol/openid-connect/token",
		LogoutRedirect:   "http://localhost:8080/",
		AuthCodeCallback: "http://localhost:8080/authCodeRedirect",
		ServicesURL:      "http://localhost:8082/billing/v1/services",
		FrontendHost:     frontendHost,
		WebPort:          port,
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
