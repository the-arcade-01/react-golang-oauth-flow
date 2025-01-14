package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/the-arcade-01/auth-flow/server/docs"
	"github.com/the-arcade-01/auth-flow/server/internal/api"
	"github.com/the-arcade-01/auth-flow/server/internal/utils"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		utils.Log.Error("error on loading env file", "error", err)
		os.Exit(1)
	}
	utils.Log.Info("env file loaded successfully", "function", "init")
}

// @title Auth Flow API
// @version 1.0
// @description This is a simple authentication flow API
// @contact.name Aashish Koshti
// @contact.url https://github.com/the-arcade-01/auth-flow
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	server := api.NewServer()
	utils.Log.Info("server running on port", "value", ":8080")
	http.ListenAndServe(":8080", server.Router)
}
