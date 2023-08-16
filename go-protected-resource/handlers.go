package main

import (
	"encoding/json"
	"net/http"
)

type Billing struct {
	Services []string `json:"services"`
}

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	s := Billing{
		Services: []string{
			"electir",
			"phone",
			"internet",
			"water",
		},
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(s)
	w.Header().Add("Content-Type", "application/json")
}
