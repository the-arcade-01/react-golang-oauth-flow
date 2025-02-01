package models

type UserDB struct {
	UserID   int    `json:"user_id"`
	GoogleID string `json:"google_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

type LoginResponse struct {
	Status  int     `json:"status"`
	Success bool    `json:"success"`
	Data    *UserDB `json:"data"`
}
