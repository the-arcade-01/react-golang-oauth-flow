package service

import (
	"context"
	"net/http"
	"time"

	"github.com/the-arcade-01/auth-flow/server/internal/models"
	"github.com/the-arcade-01/auth-flow/server/internal/repository"
)

type ApiService struct {
	repo *repository.Repository
}

func NewApiService() *ApiService {
	return &ApiService{
		repo: repository.NewRepository(),
	}
}

func (api *ApiService) Login(w http.ResponseWriter, ctx context.Context, body *models.LoginUserReqBody) (*models.LoginUserResponse, *models.ErrorResponse) {
	user := &models.User{Email: body.Email, Password: body.Password}
	accessToken, refreshToken, status, err := api.repo.LoginUser(ctx, user)
	if err != nil {
		return nil, &models.ErrorResponse{Status: status, Message: err.Error()}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",   // TODO: after testing change this to specific path
		Secure:   false, //TODO: need to change after deployment
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return &models.LoginUserResponse{
		Status: status,
		Data: struct {
			AccessToken string `json:"access_token"`
		}{
			AccessToken: accessToken,
		},
	}, nil
}

func (api *ApiService) Register(ctx context.Context, body *models.RegisterUserReqBody) (*models.RegisterUserResponse, *models.ErrorResponse) {
	user := &models.User{Email: body.Email, Password: body.Password}
	status, err := api.repo.RegisterUser(ctx, user)
	if err != nil {
		return nil, &models.ErrorResponse{Status: status, Message: err.Error()}
	}
	return &models.RegisterUserResponse{Status: status, Message: "User registered successfully"}, nil
}

func (api *ApiService) GenerateAuthTokens(w http.ResponseWriter, ctx context.Context, refreshToken string) (*models.LoginUserResponse, *models.ErrorResponse) {
	accessToken, refreshToken, status, err := api.repo.GenerateAuthTokens(ctx, refreshToken)
	if err != nil {
		return nil, &models.ErrorResponse{Status: status, Message: err.Error()}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",   // TODO: after testing change this to specific path
		Secure:   false, //TODO: need to change after deployment
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return &models.LoginUserResponse{
		Status: status,
		Data: struct {
			AccessToken string `json:"access_token"`
		}{
			AccessToken: accessToken,
		},
	}, nil
}

func (api *ApiService) Logout(ctx context.Context, userId int) (*models.LogoutResponse, *models.ErrorResponse) {
	status, err := api.repo.LogoutUser(ctx, userId)
	if err != nil {
		return nil, &models.ErrorResponse{Status: status}
	}
	return &models.LogoutResponse{Status: status}, nil
}

func (api *ApiService) GetUserById(ctx context.Context, userId int) (*models.UserResponse, *models.ErrorResponse) {
	user, status, err := api.repo.GetUserById(ctx, userId)
	if err != nil {
		return nil, &models.ErrorResponse{Status: status, Message: err.Error()}
	}
	return &models.UserResponse{Status: status, Data: user}, nil
}
