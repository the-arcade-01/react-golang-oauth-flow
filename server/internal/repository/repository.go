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
	err := r.db.QueryRowContext(ctx, query, user.Email).Scan(&row)
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
		if err == sql.ErrNoRows {
			return "", "", http.StatusBadRequest, fmt.Errorf("please check credentials")
		}
		utils.Log.ErrorContext(ctx, "error on fetching user", "function", "LoginUser", "error", err)
		return "", "", http.StatusInternalServerError, fmt.Errorf("please try again later")
	}

	isValid := checkPassword(existUser.Password, user.Password)
	if !isValid {
		return "", "", http.StatusUnauthorized, fmt.Errorf("incorrect password, please try again")
	}

	accessToken, refreshToken, status, err := r.getAuthTokens(ctx, existUser.UserId)
	return accessToken, refreshToken, status, err
}

func (r *Repository) LogoutUser(ctx context.Context, userId int) (int, error) {
	query := `delete from refresh_tokens_table where user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userId)
	if err != nil {
		utils.Log.ErrorContext(ctx, "error on deleting from refresh_tokens_table", "function", "LogoutUser", "error", err)
		return http.StatusInternalServerError, nil
	}
	return http.StatusOK, nil
}

func (r *Repository) GetUserById(ctx context.Context, userId int) (*models.User, int, error) {
	user := &models.User{}
	query := `select user_id, email, created_at from users where user_id = ?`
	err := r.db.QueryRowContext(ctx, query, userId).Scan(&user.UserId, &user.Email, &user.CreatedAt)
	if err != nil {
		utils.Log.ErrorContext(ctx, "error on fetching user", "function", "GetUserById", "error", err)
		return nil, http.StatusInternalServerError, fmt.Errorf("please try again later")
	}
	return user, http.StatusOK, nil
}

func (r *Repository) GenerateAuthTokens(ctx context.Context, refreshToken string) (string, string, int, error) {
	token, err := r.auth.Decode(refreshToken)
	if err != nil {
		return "", "", http.StatusUnauthorized, fmt.Errorf("please login again")
	}
	var claims map[string]interface{}

	if token != nil {
		claims, err = token.AsMap(context.Background())
		if err != nil {
			return "", "", http.StatusUnauthorized, fmt.Errorf("please login again")
		}
	} else {
		claims = map[string]interface{}{}
	}
	if _, ok := claims["userId"]; !ok {
		return "", "", http.StatusUnauthorized, fmt.Errorf("please login again")
	}

	userId := int(claims["userId"].(float64))

	var dbRefreshToken string
	query := `select refresh_token from refresh_tokens_table where user_id = ?`
	err = r.db.QueryRowContext(ctx, query, userId).Scan(&dbRefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", http.StatusUnauthorized, fmt.Errorf("please login again")
		}
		utils.Log.ErrorContext(ctx, "error on fetching from refresh_tokens_table", "function", "GenerateAuthTokens", "error", err)
		return "", "", http.StatusInternalServerError, fmt.Errorf("please login again")
	}

	if dbRefreshToken != refreshToken {
		return "", "", http.StatusUnauthorized, fmt.Errorf("please login again")
	}

	accessToken, refreshToken, status, err := r.getAuthTokens(ctx, userId)
	return accessToken, refreshToken, status, err
}

func (r *Repository) getAuthTokens(ctx context.Context, userId int) (string, string, int, error) {
	accessToken, err := getToken(userId, r.auth, time.Now().Add(1*time.Minute).Unix())
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("please try again later")
	}
	refreshTokenExpire := time.Now().Add(24 * time.Hour).Unix()
	refreshToken, err := getToken(userId, r.auth, refreshTokenExpire)
	if err != nil {
		return "", "", http.StatusInternalServerError, fmt.Errorf("please try again later")
	}

	query := `
		INSERT INTO refresh_tokens_table (user_id, refresh_token, expire_time) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE 
		refresh_token = VALUES(refresh_token), 
		expire_time = VALUES(expire_time), 
		created_at = CURRENT_TIMESTAMP
	`
	_, err = r.db.ExecContext(ctx, query, userId, refreshToken, refreshTokenExpire)
	if err != nil {
		utils.Log.ErrorContext(ctx, "error on saving refresh token in db", "function", "LoginUser", "error", err)
		return "", "", http.StatusInternalServerError, fmt.Errorf("please try again later")
	}
	return accessToken, refreshToken, http.StatusOK, nil
}

func getToken(userId int, auth *jwtauth.JWTAuth, expireTime int64) (string, error) {
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
