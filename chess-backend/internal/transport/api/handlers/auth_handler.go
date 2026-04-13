package handlers

import (
	"encoding/json"
	"net/http"

	"chess-backend/internal/app/auth"
	"chess-backend/internal/constants/appconst"
	"chess-backend/internal/transport/dto"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration info"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 409 {string} string "Email already registered"
// @Failure 500 {string} string "Internal server error"
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, appconst.MsgInvalidRequest, http.StatusBadRequest)
		return
	}

	cmd := auth.RegisterCommand{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	}

	u, token, err := h.authService.Register(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:       u.ID,
			Email:    u.Email,
			Username: u.Username,
		},
	}
	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(resp)
}

// Login godoc
// @Summary Login user
// @Description Authenticate and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid credentials"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, appconst.MsgInvalidRequest, http.StatusBadRequest)
		return
	}

	cmd := auth.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	u, token, err := h.authService.Login(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:       u.ID,
			Email:    u.Email,
			Username: u.Username,
		},
	}
	w.Header().Set(appconst.HeaderContentType, appconst.MimeTypeJSON)
	json.NewEncoder(w).Encode(resp)
}
