package main

import (
	"net/http"

	"github.com/the-arcade-01/auth-flow/server/internal/api"
	"github.com/the-arcade-01/auth-flow/server/internal/utils"
)

func main() {
	server := api.NewServer()
	utils.Log.Info("server running on port", "value", ":8080")
	http.ListenAndServe(":8080", server.Router)
}
