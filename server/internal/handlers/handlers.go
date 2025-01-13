package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/the-arcade-01/auth-flow/server/internal/models"
	"github.com/the-arcade-01/auth-flow/server/internal/service"
)

type Handlers struct {
	api *service.ApiService
}

func NewHandlers() *Handlers {
	return &Handlers{
		api: service.NewApiService(),
	}
}

func (h *Handlers) Greet(w http.ResponseWriter, r *http.Request) {
	models.ResponseWithJSON(w, http.StatusOK, "Hello, World")
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var body *models.LoginUserReqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		models.ResponseWithJSON(w, http.StatusBadRequest, "Please provide correct details")
		return
	}
	defer r.Body.Close()
	result, err := h.api.Login(w, r.Context(), body)
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}
	models.ResponseWithJSON(w, result.Status, result)
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var body *models.RegisterUserReqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		models.ResponseWithJSON(w, http.StatusBadRequest, "Please provide correct details")
		return
	}
	defer r.Body.Close()
	result, err := h.api.Register(r.Context(), body)
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}
	models.ResponseWithJSON(w, result.Status, result)
}

func (h *Handlers) GenerateRefreshToken(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) DeleteRefreshToken(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	_, claims, _ := jwtauth.FromContext(r.Context())
	result, err := h.api.Logout(r.Context(), int(claims["userId"].(float64)))
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}
	models.ResponseWithJSON(w, result.Status, result)
}

func (h *Handlers) GetUserById(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	result, err := h.api.GetUserById(r.Context(), int(claims["userId"].(float64)))
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}
	models.ResponseWithJSON(w, result.Status, result)
}
