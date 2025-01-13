package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/jwtauth/v5"
	"github.com/the-arcade-01/auth-flow/server/internal/config"
	"github.com/the-arcade-01/auth-flow/server/internal/models"
	"github.com/the-arcade-01/auth-flow/server/internal/utils"
)

type Repository struct {
	db   *sql.DB
	auth *jwtauth.JWTAuth
}

func NewRepository() *Repository {
	config := config.NewAppConfig()
	return &Repository{
		db:   config.DbClient,
		auth: config.AuthToken,
	}
}

func (r *Repository) RegisterUser(ctx context.Context, user *models.User) (int, error) {
	query := `select count(email) from users where email = ?`
	var row int
	err := r.db.QueryRowContext(ctx, query, user.Email).Scan(row)
	if err != nil {
		utils.Log.ErrorContext(ctx, "error fetching user", "function", "RegisterUser", "error", err)
		return http.StatusInternalServerError, fmt.Errorf("please try again later")
	}

	if row != 0 {
		return http.StatusBadRequest, fmt.Errorf("email already taken, please use different email")
	}

	hashPassword, err := getHashPassword(user.Password)
	if err != nil {
		utils.Log.ErrorContext(ctx, "error on generating hash password", "function", "RegisterUser", "error", err)
		return http.StatusInternalServerError, fmt.Errorf("please try again later")
	}
	query = `insert into users (email, password) values (?,?)`
	_, err = r.db.ExecContext(ctx, query, user.Email, hashPassword)
	if err != nil {
		utils.Log.ErrorContext(ctx, "error on saving user in db", "function", "RegisterUser", "error", err)
		return http.StatusInternalServerError, fmt.Errorf("please try again later")
	}

	return http.StatusOK, nil
}

func (r *Repository) LoginUser(ctx context.Context, user *models.User) (string, string, int, error) {
	existUser := &models.User{}
	query := `select user_id, email, password from users where email = ?`
	err := r.db.QueryRowContext(ctx, query, user.Email).Scan(&existUser.UserId, &existUser.Email, &existUser.Password)
	if err != nil {
		utils.Log.ErrorContext(ctx, "error on fetching user", "function", "LoginUser", "error", err)
		return "", "", http.StatusInternalServerError, fmt.Errorf("please try again later")
	}

	isValid := checkPassword(existUser.Password, user.Password)
	if !isValid {
		return "", "", http.StatusUnauthorized, fmt.Errorf("incorrect password, please try again")
	}

	accessToken, err := generateAuthToken(existUser.UserId, r.auth, time.Now().Add(15*time.Minute))
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("please try again later")
	}
	refreshToken, err := generateAuthToken(existUser.UserId, r.auth, time.Now().Add(24*time.Hour))
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("please try again later")
	}

	query = `insert into refresh_tokens_table (user_id, refresh_token) values (?,?) ON DUPLICATE KEY UPDATE refresh_token = values(refresh_token), created_at = CURRENT_TIMESTAMP`
	_, err = r.db.ExecContext(ctx, query, existUser.UserId, refreshToken)
	if err != nil {
		utils.Log.ErrorContext(ctx, "error on saving refresh token in db", "function", "LoginUser", "error", err)
		return "", "", http.StatusInternalServerError, fmt.Errorf("please try again later")
	}

	return accessToken, refreshToken, http.StatusOK, nil
}

func (r *Repository) LogoutUser(ctx context.Context, userId int) (int, error) {
	query := `delete from refresh_tokens_table where userId = ?`
	_, err := r.db.ExecContext(ctx, query, userId)
	if err != nil {
		utils.Log.ErrorContext(ctx, "error on deleting from refresh_tokens_table", "function", "LogoutUser", "error", err)
		return http.StatusInternalServerError, nil
	}
	return http.StatusOK, nil
}

func (r *Repository) GetUserById(ctx context.Context, userId int) (*models.User, int, error) {
	var user *models.User
	query := `select (user_id, email, created_at) from users where user_id = ?`
	err := r.db.QueryRowContext(ctx, query, userId).Scan(&user)
	if err != nil {
		utils.Log.ErrorContext(ctx, "error on fetching user", "function", "GetUserById", "error", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("please try again later")
	}
	return user, http.StatusOK, nil
}

func getHashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func generateAuthToken(userId int, auth *jwtauth.JWTAuth, expireTime time.Time) (string, error) {
	claims := map[string]interface{}{
		"userId": userId,
		"exp":    expireTime,
	}
	_, token, err := auth.Encode(claims)
	if err != nil {
		utils.Log.Error("error on generating auth token", "function", "generateAuthToken", "error", err)
		return "", err
	}
	return token, nil
}
