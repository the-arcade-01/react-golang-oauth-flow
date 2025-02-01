package service

import (
	"context"
	"server/internal/models"
	"server/internal/repository"
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo: repository.New(),
	}
}

func (svc *AuthService) Login(ctx context.Context, userInfo map[string]string) (*models.LoginResponse, *models.Response) {
	user := &models.UserDB{
		GoogleID: userInfo["id"],
		Email:    userInfo["email"],
		Name:     userInfo["name"],
		Picture:  userInfo["picture"],
	}
	userID, status, err := svc.repo.Login(ctx, user)
	if err != nil {
		return nil, &models.Response{Status: status, Success: false, Message: err.Error()}
	}
	user.UserID = userID
	return &models.LoginResponse{Status: status, Success: true, Data: user}, nil
}
