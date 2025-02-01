package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"server/internal/config"
	"server/internal/models"
	"server/internal/utils"
)

type Repository struct {
	db *sql.DB
}

func New() *Repository {
	config := config.New()
	return &Repository{
		db: config.DBClient,
	}
}

func (r *Repository) Login(ctx context.Context, user *models.UserDB) (int, int, error) {
	query := `select user_id from users where email = ?`
	var userID int
	err := r.db.QueryRowContext(ctx, query, user.Email).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		utils.Log.ErrorContext(ctx, "error on checking existing user", "func", "Login", "err", err)
		return -1, http.StatusInternalServerError, fmt.Errorf("please try again later")
	}

	if err == sql.ErrNoRows {
		query = `insert into users (google_id, email, name, picture) values (?,?,?,?)`
		result, err := r.db.ExecContext(ctx, query, user.GoogleID, user.Email, user.Name, user.Picture)
		if err != nil {
			utils.Log.ErrorContext(ctx, "error on saving user in db", "func", "Login", "err", err)
			return -1, http.StatusInternalServerError, fmt.Errorf("please try again later")
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			utils.Log.ErrorContext(ctx, "error on getting lastID after saving user", "func", "Login", "err", err)
			return -1, http.StatusInternalServerError, fmt.Errorf("please try again later")
		}
		userID = int(lastID)
	}

	return userID, http.StatusOK, nil
}
