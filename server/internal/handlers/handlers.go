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

// Greet godoc
// @Summary Greet the user
// @Description Responds with a greeting message
// @Tags greet
// @Produce json
// @Success 200 {string} string "Hello, World"
// @Router /greet [get]
func (h *Handlers) Greet(w http.ResponseWriter, r *http.Request) {
	models.ResponseWithJSON(w, http.StatusOK, "Hello, World")
}

// Login godoc
// @Summary Login a user
// @Description Authenticates a user and returns tokens. Sets a secure HTTP-only cookie for the refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.LoginUserReqBody true "Login request body"
// @Success 200 {object} models.LoginUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var body *models.LoginUserReqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		models.ResponseWithJSON(w, http.StatusBadRequest, &models.ErrorResponse{Status: http.StatusBadRequest, Error: "Please provide correct details"})
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

// Register godoc
// @Summary Register a new user
// @Description Registers a new user and returns tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.RegisterUserReqBody true "Register request body"
// @Success 200 {object} models.RegisterUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var body *models.RegisterUserReqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		models.ResponseWithJSON(w, http.StatusBadRequest, &models.ErrorResponse{Status: http.StatusBadRequest, Error: "Please provide correct details"})
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

// GenerateAuthTokens godoc
// @Summary Generate new auth tokens
// @Description Generates new access and refresh tokens using the existing refresh token in cookies. Sets a new secure HTTP-only cookie for the refresh token.
// @Tags auth
// @Produce json
// @Success 200 {object} models.LoginUserResponse
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Refresh token is missing or invalid"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /refresh-token [post]
func (h *Handlers) GenerateAuthTokens(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		models.ResponseWithJSON(w, http.StatusUnauthorized, &models.ErrorResponse{Status: http.StatusBadRequest, Error: "Please login again"})
		return
	}

	refreshToken := cookie.Value
	result, er := h.api.GenerateAuthTokens(w, r.Context(), refreshToken)
	if er != nil {
		models.ResponseWithJSON(w, er.Status, er)
		return
	}
	models.ResponseWithJSON(w, result.Status, result)
}

// Logout godoc
// @Summary Logout a user
// @Description Logs out a user, invalidates the refresh token, and clears the refresh token cookie
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.LogoutResponse
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /logout [post]
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	userId, err := validateAuthToken(r)
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}
	result, err := h.api.Logout(r.Context(), userId)
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
		Secure:   false, // TODO: change for production
		SameSite: http.SameSiteNoneMode,
	})
	models.ResponseWithJSON(w, result.Status, result)
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Retrieves the current authenticated user's information
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse "Unauthorized: Missing or invalid token"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /users/me [get]
func (h *Handlers) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userId, err := validateAuthToken(r)
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}
	result, err := h.api.GetUserById(r.Context(), userId)
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}
	models.ResponseWithJSON(w, result.Status, result)
}

// No needed but okay for parsing userId from claims
func validateAuthToken(r *http.Request) (int, *models.ErrorResponse) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return -1, &models.ErrorResponse{Status: http.StatusUnauthorized, Error: "User unauthorized, please login again"}
	}
	if _, ok := claims["userId"]; !ok {
		return -1, &models.ErrorResponse{Status: http.StatusUnauthorized, Error: "User unauthorized, please login again"}
	}
	if _, ok := claims["exp"]; !ok {
		return -1, &models.ErrorResponse{Status: http.StatusUnauthorized, Error: "User unauthorized, please login again"}
	}
	userId := int(claims["userId"].(float64))
	expTime := claims["exp"].(time.Time).Unix()

	if expTime < time.Now().Unix() {
		return -1, &models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "token has expired",
		}
	}

	return userId, nil
}
