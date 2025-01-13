package config

import (
	"database/sql"
	"os"
	"sync"

	"github.com/go-chi/jwtauth/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/the-arcade-01/auth-flow/server/internal/utils"
)

var once = sync.Once{}
var config *AppConfig

type AppConfig struct {
	DbClient  *sql.DB
	AuthToken *jwtauth.JWTAuth
}

func NewAppConfig() *AppConfig {
	once.Do(func() {
		config = &AppConfig{}
		db, err := newDBClient()
		if err != nil {
			panic(err)
		}
		config.DbClient = db
		config.AuthToken = generateAuthToken()
	})
	return config
}

func newDBClient() (*sql.DB, error) {
	db, err := sql.Open(os.Getenv("DRIVER_NAME"), os.Getenv("DB_URL"))
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	if err != nil {
		utils.Log.Error("error establishing db conn", "error", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		utils.Log.Error("error on pinging db", "error", err)
		return nil, err
	}
	return db, err
}

func generateAuthToken() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil)
}
