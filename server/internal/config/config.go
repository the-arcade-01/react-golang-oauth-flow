package config

import "sync"

var once = sync.Once{}
var config *AppConfig

type AppConfig struct {
}

func NewAppConfig() *AppConfig {
	once.Do(func() {
		config = &AppConfig{}
	})
	return config
}
