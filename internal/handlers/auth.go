package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/cunnati/store-service/internal/auth"
	"github.com/cunnati/store-service/internal/models"
)

type AuthHandler struct {
	Users     *models.UserStore
	JWTSecret string
}

type signupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type tokenResp struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || len(req.Password) < 8 {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	ph, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "hash error", 500)
		return
	}
	_, err = h.Users.Create(ctx, req.Email, ph)
	if err != nil {
		http.Error(w, "cannot create", 409)
		return
	}
	tok, err := auth.GenerateJWT(h.JWTSecret, req.Email, 24*time.Hour)
	if err != nil {
		http.Error(w, "token error", 500)
		return
	}
	writeJSON(w, http.StatusCreated, tokenResp{Token: tok})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	u, err := h.Users.FindByEmail(ctx, req.Email)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	tok, err := auth.GenerateJWT(h.JWTSecret, u.ID, 24*time.Hour)
	if err != nil {
		http.Error(w, "token error", 500)
		return
	}
	writeJSON(w, http.StatusOK, tokenResp{Token: tok})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
