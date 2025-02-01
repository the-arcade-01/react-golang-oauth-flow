package config

import (
	"errors"
	"os"
	"server/internal/utils"
	"strconv"
	"sync"
)

var Envs *AppEnvs
var envOnce sync.Once

type AppEnvs struct {
	ENV                       string
	GOOGLE_CLIENT_ID          string
	GOOGLE_CLIENT_SECRET      string
	AUTH_REDIRECT_URL         string
	GOOGLE_USER_INFO          string
	DB_DRIVER                 string
	DB_URL                    string
	DB_MAX_IDLE_CONN          int
	DB_MAX_OPEN_CONN          int
	DB_MAX_CONN_TIME_SEC      int
	HTTP_COOKIE_HTTPONLY      bool
	HTTP_COOKIE_SECURE        bool
	HTTP_REFRESH_TOKEN_EXPIRE int
	APP_WEB_URL               string
	APP_WEB_URL_LOGIN_SUCCESS string
	APP_WEB_URL_LOGIN_ERROR   string
}

func ParseEnvs() (*AppEnvs, error) {
	var err error
	envOnce.Do(func() {
		Envs = &AppEnvs{
			ENV:                       os.Getenv("ENV"),
			GOOGLE_CLIENT_ID:          os.Getenv("GOOGLE_CLIENT_ID"),
			GOOGLE_CLIENT_SECRET:      os.Getenv("GOOGLE_CLIENT_SECRET"),
			AUTH_REDIRECT_URL:         os.Getenv("AUTH_REDIRECT_URL"),
			GOOGLE_USER_INFO:          os.Getenv("GOOGLE_USER_INFO"),
			DB_DRIVER:                 os.Getenv("DB_DRIVER"),
			DB_URL:                    os.Getenv("DB_URL"),
			APP_WEB_URL:               os.Getenv("APP_WEB_URL"),
			APP_WEB_URL_LOGIN_SUCCESS: os.Getenv("APP_WEB_URL_LOGIN_SUCCESS"),
			APP_WEB_URL_LOGIN_ERROR:   os.Getenv("APP_WEB_URL_LOGIN_ERROR"),
		}

		if Envs.ENV == "" || Envs.GOOGLE_CLIENT_ID == "" || Envs.GOOGLE_CLIENT_SECRET == "" || Envs.AUTH_REDIRECT_URL == "" || Envs.DB_DRIVER == "" || Envs.DB_URL == "" || Envs.GOOGLE_USER_INFO == "" || Envs.APP_WEB_URL == "" || Envs.APP_WEB_URL_LOGIN_SUCCESS == "" || Envs.APP_WEB_URL_LOGIN_ERROR == "" {
			utils.Log.Error("Error: One or more required environment variables are missing or empty")
			err = errors.New("missing required environment variables")
			return
		}

		dbIdleConn, parseErr := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONN"))
		if parseErr != nil {
			utils.Log.Error("Error parsing DB_MAX_IDLE_CONN", "err", parseErr)
			err = parseErr
			return
		}
		Envs.DB_MAX_IDLE_CONN = dbIdleConn

		dbOpenConn, parseErr := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONN"))
		if parseErr != nil {
			utils.Log.Error("Error parsing DB_MAX_OPEN_CONN", "err", parseErr)
			err = parseErr
			return
		}
		Envs.DB_MAX_OPEN_CONN = dbOpenConn

		dbMaxConnTimeSec, parseErr := strconv.Atoi(os.Getenv("DB_MAX_CONN_TIME_SEC"))
		if parseErr != nil {
			utils.Log.Error("Error parsing DB_MAX_CONN_TIME_SEC", "err", parseErr)
			err = parseErr
			return
		}
		Envs.DB_MAX_CONN_TIME_SEC = dbMaxConnTimeSec

		httpCookieHttpOnly, parseErr := strconv.ParseBool(os.Getenv("HTTP_COOKIE_HTTPONLY"))
		if parseErr != nil {
			utils.Log.Error("Error parsing HTTP_COOKIE_HTTPONLY", "err", parseErr)
			err = parseErr
			return
		}
		Envs.HTTP_COOKIE_HTTPONLY = httpCookieHttpOnly

		httpCookieSecure, parseErr := strconv.ParseBool(os.Getenv("HTTP_COOKIE_SECURE"))
		if parseErr != nil {
			utils.Log.Error("Error parsing HTTP_COOKIE_SECURE", "err", parseErr)
			err = parseErr
			return
		}
		Envs.HTTP_COOKIE_SECURE = httpCookieSecure

		httpRefreshTokenExpire, parseErr := strconv.Atoi(os.Getenv("HTTP_REFRESH_TOKEN_EXPIRE"))
		if parseErr != nil {
			utils.Log.Error("Error parsing HTTP_REFRESH_TOKEN_EXPIRE", "err", parseErr)
			err = parseErr
			return
		}
		Envs.HTTP_REFRESH_TOKEN_EXPIRE = httpRefreshTokenExpire
	})

	if err != nil {
		return nil, err
	}

	utils.Log.Info("env variables loaded successfully", "func", "ParseEnvs")
	// utils.Log.Info("system environment", "func", "ParseEnvs", "env", Envs.ENV)
	return Envs, nil
}
