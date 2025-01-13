package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/the-arcade-01/auth-flow/server/internal/config"
	"github.com/the-arcade-01/auth-flow/server/internal/handlers"
	"github.com/the-arcade-01/auth-flow/server/internal/utils"
)

type Server struct {
	Router *chi.Mux
}

func (s *Server) mountMiddlewares() {
	s.Router.Use(middleware.Heartbeat("/ping"))
	s.Router.Use(requestLogger)
	s.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}

func (s *Server) mountHandlers() {
	appConfig := config.NewAppConfig()
	handlers := handlers.NewHandlers()
	s.Router.Get("/", handlers.Greet)
	s.Router.Post("/login", handlers.Login)
	s.Router.Post("/register", handlers.Register)
	s.Router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(appConfig.AuthToken))
		r.Use(jwtauth.Authenticator(appConfig.AuthToken))

		r.Post("/logout", handlers.Logout)
		r.Post("/refresh-token", handlers.GenerateRefreshToken)
		r.Delete("/refresh-token", handlers.DeleteRefreshToken)
		r.Get("/users/:userId", handlers.GetUserById)
	})
}

func NewServer() *Server {
	s := &Server{
		Router: chi.NewRouter(),
	}
	s.mountMiddlewares()
	s.mountHandlers()
	return s
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ms := time.Since(start).Milliseconds()

		utils.Log.Info("rtt",
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"ms", ms,
		)
	})
}
