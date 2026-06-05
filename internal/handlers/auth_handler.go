package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jtorre/qisurChallenge/internal/auth"
	"github.com/jtorre/qisurChallenge/internal/config"
	"github.com/jtorre/qisurChallenge/internal/repository"
	"github.com/jtorre/qisurChallenge/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

func NewAuthHandler(userRepo *repository.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateNonEmpty(req.Email, "email"); err != nil {
		RespondError(w, err)
		return
	}

	if err := utils.ValidateNonEmpty(req.Password, "password"); err != nil {
		RespondError(w, err)
		return
	}

	if err := utils.ValidateMinLength(req.Password, 6, "password"); err != nil {
		RespondError(w, err)
		return
	}

	// Check if user already exists
	_, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err == nil {
		RespondConflict(w, "email already registered")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondError(w, err)
		return
	}

	// Create user with 'client' role by default
	user, err := h.userRepo.Create(r.Context(), req.Email, string(hashedPassword), "client")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, h.cfg.JWTSecret, h.cfg.JWTExpirationHours)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
		Token: token,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateNonEmpty(req.Email, "email"); err != nil {
		RespondError(w, err)
		return
	}

	if err := utils.ValidateNonEmpty(req.Password, "password"); err != nil {
		RespondError(w, err)
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		RespondUnauthorized(w, "invalid credentials")
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		RespondUnauthorized(w, "invalid credentials")
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, h.cfg.JWTSecret, h.cfg.JWTExpirationHours)
	if err != nil {
		RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
		Token: token,
	})
}
