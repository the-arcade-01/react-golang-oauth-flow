package main

import (
	"net/http"
	"server/internal/api"
	"server/internal/config"
	"server/internal/utils"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		utils.Log.Error("Error loading .env file", "err", err)
		panic(err)
	}
	_, err = config.ParseEnvs()
	if err != nil {
		panic(err)
	}
}

func main() {
	server := api.NewServer()
	utils.Log.Info("server running on port:8080", "func", "main")
	http.ListenAndServe(":8080", server.Router)
}
