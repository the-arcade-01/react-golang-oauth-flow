package api

import (
	"context"
	"encoding/json"
	"net/http"
	"server/internal/config"
	"server/internal/handlers"
	"server/internal/models"
	"server/internal/utils"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/oauth2"
)

type Server struct {
	Router *chi.Mux
}

func (s *Server) mountMiddlewares() {
	s.Router.Use(middleware.Heartbeat("/ping"))
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(1 * time.Minute))
	s.Router.Use(requestLogger)
	s.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{config.Envs.APP_WEB_URL},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func (s *Server) mountHandlers() {
	authHandlers := handlers.NewAuthHandlers()
	authRouter := chi.NewRouter()
	authRouter.Get("/", authHandlers.Greet)
	authRouter.Get("/login", authHandlers.Login)
	authRouter.Get("/callback", authHandlers.Callback)
	authRouter.Post("/refresh-token", authHandlers.RefreshToken)
	authRouter.Group(func(r chi.Router) {
		r.Use(validateRequestToken)
		r.Get("/protected", authHandlers.Protected)
		r.Post("/logout", authHandlers.Logout)
		r.Get("/user/me", authHandlers.FetchUser)
	})
	s.Router.Mount("/api/auth", authRouter)
}

func NewServer() *Server {
	s := &Server{
		Router: chi.NewRouter(),
	}
	s.mountMiddlewares()
	s.mountHandlers()
	return s
}

func validateRequestToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("access_token")
		if err != nil || accessTokenCookie.Value == "" {
			models.ResponseWithJSON(w, http.StatusUnauthorized, &models.Response{Success: false, Status: http.StatusUnauthorized, Message: "Access token is missing"})
			return
		}

		client := config.New().OauthCfg.Client(context.Background(), &oauth2.Token{
			AccessToken: accessTokenCookie.Value,
		})

		res, err := client.Get(config.Envs.GOOGLE_USER_INFO)
		if err != nil || res.StatusCode != http.StatusOK {
			models.ResponseWithJSON(w, http.StatusUnauthorized, &models.Response{Success: false, Status: http.StatusUnauthorized, Message: "Invalid Access token"})
			return
		}
		defer res.Body.Close()

		var userInfo map[string]string
		json.NewDecoder(res.Body).Decode(&userInfo)

		ctx := context.WithValue(r.Context(), utils.UserInfoKey, userInfo)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)

		ms := time.Since(start).Milliseconds()
		utils.Log.Info("rtt",
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"ms", ms,
			"status", rec.status,
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
