package config

import (
	"database/sql"
	"sync"

	"golang.org/x/oauth2"
)

var once sync.Once
var appConfig *AppConfig

type AppConfig struct {
	DBClient *sql.DB
	OauthCfg *oauth2.Config
}

func New() *AppConfig {
	once.Do(func() {
		appConfig = &AppConfig{
			OauthCfg: newOauthConfig(),
		}
		db, err := newDBClient()
		if err != nil {
			// TODO: change this to graceful shutdown
			panic(err)
		}
		appConfig.DBClient = db

	})
	return appConfig
}
