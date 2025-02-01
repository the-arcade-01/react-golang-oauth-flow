package models

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ResponseWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
