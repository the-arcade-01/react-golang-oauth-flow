package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"server/internal/config"
	"server/internal/models"
	"server/internal/service"
	"server/internal/utils"
	"time"

	"golang.org/x/oauth2"
)

type AuthHandlers struct {
	config *config.AppConfig
	svc    *service.AuthService
}

func NewAuthHandlers() *AuthHandlers {
	return &AuthHandlers{
		config: config.New(),
		svc:    service.NewAuthService(),
	}
}

func (h *AuthHandlers) Greet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World"))
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		models.ResponseWithJSON(w, http.StatusInternalServerError, "Please try again later")
		return
	}
	url := h.config.OauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandlers) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := h.config.OauthCfg.Exchange(r.Context(), code)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s=%s", config.Envs.APP_WEB_URL_LOGIN_ERROR, "exchange_failed"), http.StatusTemporaryRedirect)
		return
	}

	client := h.config.OauthCfg.Client(r.Context(), token)
	res, err := client.Get(config.Envs.GOOGLE_USER_INFO)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s=%s", config.Envs.APP_WEB_URL_LOGIN_ERROR, "userinfo_failed"), http.StatusTemporaryRedirect)
		return
	}
	defer res.Body.Close()

	var userInfo map[string]string
	json.NewDecoder(res.Body).Decode(&userInfo)

	_, errResult := h.svc.Login(r.Context(), userInfo)
	if errResult != nil {
		http.Redirect(w, r, fmt.Sprintf("%s=%s", config.Envs.APP_WEB_URL_LOGIN_ERROR, "login_failed"), http.StatusTemporaryRedirect)
		return
	}

	h.setCookies(w, token)
	http.Redirect(w, r, config.Envs.APP_WEB_URL_LOGIN_SUCCESS, http.StatusTemporaryRedirect)
}

func (h *AuthHandlers) Protected(w http.ResponseWriter, r *http.Request) {
	models.ResponseWithJSON(w, http.StatusOK, "Me protected")
}

func (h *AuthHandlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil || refreshToken.Value == "" {
		models.ResponseWithJSON(w, http.StatusUnauthorized, &models.Response{Status: http.StatusUnauthorized, Success: false, Message: "Please login again"})
		return
	}
	tokenSource := h.config.OauthCfg.TokenSource(r.Context(), &oauth2.Token{
		RefreshToken: refreshToken.Value,
	})
	newToken, err := tokenSource.Token()
	if err != nil {
		models.ResponseWithJSON(w, http.StatusUnauthorized, &models.Response{Status: http.StatusUnauthorized, Success: false, Message: "Please login again"})
		return
	}

	h.setCookies(w, newToken)
	models.ResponseWithJSON(w, http.StatusOK, &models.Response{Status: http.StatusOK, Success: true})
}

func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: config.Envs.HTTP_COOKIE_HTTPONLY,
		Secure:   config.Envs.HTTP_COOKIE_SECURE,
		SameSite: http.SameSiteLaxMode,
	})
	models.ResponseWithJSON(w, http.StatusOK, models.Response{Success: true, Status: http.StatusOK, Message: "Logout successfull"})
}

func (h *AuthHandlers) FetchUser(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(utils.UserInfoKey).(map[string]string)
	if !ok {
		models.ResponseWithJSON(w, http.StatusBadRequest, "Invalid user info")
		return
	}

	result, errResult := h.svc.Login(r.Context(), userInfo)
	if errResult != nil {
		models.ResponseWithJSON(w, errResult.Status, errResult)
		return
	}

	models.ResponseWithJSON(w, result.Status, result)
}

func (h *AuthHandlers) setCookies(w http.ResponseWriter, token *oauth2.Token) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		Expires:  token.Expiry,
		HttpOnly: config.Envs.HTTP_COOKIE_HTTPONLY,
		Secure:   config.Envs.HTTP_COOKIE_SECURE,
		SameSite: http.SameSiteLaxMode,
	})

	// NOTE: Refresh token is only issued at the first consent
	if token.RefreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    token.RefreshToken,
			Expires:  time.Now().Add(time.Duration(config.Envs.HTTP_REFRESH_TOKEN_EXPIRE) * time.Hour),
			HttpOnly: config.Envs.HTTP_COOKIE_HTTPONLY,
			Secure:   config.Envs.HTTP_COOKIE_SECURE,
			SameSite: http.SameSiteLaxMode,
		})
	}
}

func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
