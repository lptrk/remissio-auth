package auth

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"remissio-auth/utils"
	"strings"
	"time"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

var ErrAuth = errors.New("unauthorized")

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if !isEmailValid(email) || !isUsernameValid(username) || !isPasswordValid(password) {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	exists, err := h.Service.UserExists(username, email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("[Error] UserExists check failed:", err)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		http.Error(w, "Internal server error while hashing password", http.StatusInternalServerError)
		log.Println("[Error] Error hashing password:", err)
		return
	}

	user := &User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	if err := h.Service.Create(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Println("[Error] Failed to create user:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "User registered successfully"}`))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	log.Printf("[Info] Incoming Login request for user: %s", username)

	password := r.FormValue("password")

	user, err := h.Service.GetByUsername(username)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(25)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	h.Service.SetSessionToken(sessionToken, username)
	h.Service.SetCSRFToken(csrfToken, username)

	log.Printf("[Info] Login successful for user %s", user.Username)
}

func (h *Handler) Authorize(r *http.Request) error {
	log.Println("[Info] Authorizing request")
	var username string
	switch r.Method {
	case "GET":
		username = r.URL.Query().Get("username")
	default:
		username = r.FormValue("username")
	}

	user, err := h.Service.GetByUsername(username)
	if err != nil {
		log.Fatalf("[Error] Error while fetching user by username")
	}

	log.Printf("[Info] Got user by username: %s", username)

	st, err := r.Cookie("session_token")
	if err != nil || st.Value == "" || st.Value != user.SessionToken {
		log.Printf("[Error] Error while verifying session token for user %s: expected: %s actual %s", username, user.SessionToken, st.Value)
		return ErrAuth
	}

	log.Println("[Info] Cookie was confirmed successfully")

	csrf := r.Header.Get("X-CSRF-Token")
	if csrf != user.CSRFToken || csrf == "" {
		return ErrAuth
	}

	log.Println("[Info] CSRF token was confirmed successfully")
	return nil
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	log.Println("[Info] Resetting cookie...")

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	log.Println("[Info] Resetting CSRF token...")

	err := h.Service.ResetTokens(username)
	if err != nil {
		log.Fatalf("[Error] An error occurred while resetting the tokens for user %s", username)
		return
	}
}

func isEmailValid(e string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(regex).MatchString(e)
}

func isUsernameValid(u string) bool {
	specialChars := "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	return !strings.ContainsAny(u, specialChars) && len(u) > 4
}

func isPasswordValid(p string) bool {
	hasSpecial := strings.ContainsAny(p, "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
	hasNumber := strings.ContainsAny(p, "0123456789")
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(p)
	return hasSpecial && hasNumber && hasUpper && len(p) > 8
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Info] GET params were: %s", r.URL.Query())
	if err := h.Authorize(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username := r.URL.Query().Get("username")

	log.Printf("[Info] User %s authorized successfully for route /test", username)
}
