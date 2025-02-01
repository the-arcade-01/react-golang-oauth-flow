package config

import (
	"database/sql"
	"server/internal/utils"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func newDBClient() (*sql.DB, error) {
	db, err := sql.Open(Envs.DB_DRIVER, Envs.DB_URL)
	if err != nil {
		utils.Log.Error("error on establishing db connection", "err", err)
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(Envs.DB_MAX_CONN_TIME_SEC) * time.Second)
	db.SetMaxOpenConns(Envs.DB_MAX_OPEN_CONN)
	db.SetMaxIdleConns(Envs.DB_MAX_IDLE_CONN)

	if err := db.Ping(); err != nil {
		utils.Log.Error("error on ping db connection", "err", err)
		return nil, err
	}
	return db, nil
}
