package models

import (
	"encoding/json"
	"net/http"
	"time"
)

// TODO: request body validation remaining

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type User struct {
	UserId    int       `json:"user_id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterUserReqBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type LoginUserReqBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Status int `json:"status"`
	Data   struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
}

type UserResponse struct {
	Status int   `json:"status"`
	Data   *User `json:"data"`
}

type LogoutResponse struct {
	Status int `json:"status"`
}

func ResponseWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
