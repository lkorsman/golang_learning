package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	userStore UserStore
	jwtManager *JWTManager
}

func NewHandler(userStore UserStore, jwtManager *JWTManager) *Handler {
	return &Handler{
		userStore: userStore,
		jwtManager: jwtManager,
	}
}

type RegisterRequest struct {
	Email	 string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email	string  `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User	 `json:"user"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, "password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	user, err := h.userStore.Create(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	token, err := h.jwtManager.Generate(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.userStore.GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	memStore, ok := h.userStore.(*MemoryUserStore)
	if !ok {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := memStore.ValidatePassword(user.Password, req.Password); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.jwtManager.Generate(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User: user,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}