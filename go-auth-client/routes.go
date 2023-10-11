package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"

	"github.com/gorilla/mux"
)

func routes(handlerConfig *HandlerConfig) *mux.Router {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(w, r, handlerConfig)
	})
	router.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		servicesHandler(w, r, handlerConfig)
	})
	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logoutHandler(w, r, handlerConfig)
	})
	router.HandleFunc("/authCodeRedirect", func(w http.ResponseWriter, r *http.Request) {
		authCodeRedirectHandler(w, r, handlerConfig)
	})
	router.HandleFunc("/refreshToken", func(w http.ResponseWriter, r *http.Request) {
		refreshTokenHandler(w, r, handlerConfig)
	})

	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println("Error dumping request:", err)
		} else {
			log.Printf("Request:\n%s\n", requestDump)
		}

		var body []byte
		if r.Header.Get("Content-Type") == "application/json" {
			body, err = io.ReadAll(r.Body)
			if err != nil {
				log.Println("Error reading request body:", err)
			} else {
				log.Printf("Request Body:\n%s\n", body)
			}
			r.Body = io.NopCloser(bytes.NewReader(body)) // Restore the body for further processing
		}

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		// Log response status and headers
		log.Printf("Response Status: %d\n", recorder.Code)
		for key, values := range recorder.Header() {
			for _, value := range values {
				log.Printf("Response Header - %s: %s\n", key, value)
				w.Header().Add(key, value)
			}
		}

		// Write response status
		w.WriteHeader(recorder.Code)

		// Write the response body to the actual response writer
		responseBody := recorder.Body.Bytes()

		if len(responseBody) > 0 && strings.Contains(r.Header.Get("Accept"), "application/json") {
			log.Printf("Response:\n%s\n", responseBody)
			_, err := w.Write(responseBody)
			if err != nil {
				log.Println("Error writing response:", err)
			}
		} else {
			_, err := w.Write(responseBody)
			if err != nil {
				log.Println("Error writing response:", err)
			}
		}
	})
}
