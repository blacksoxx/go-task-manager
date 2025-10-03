package handlers

import (
	"encoding/json"
	"net/http"
	"auth-service/internal/database"
	"auth-service/internal/models"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db *database.DB
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		http.Error(w, `{"error": "Email, password, first name, and last name are required"}`, http.StatusBadRequest)
		return
	}

	// Check if user already exists
	_, err := h.db.GetUserByEmail(req.Email)
	if err == nil {
		http.Error(w, `{"error": "User already exists"}`, http.StatusConflict)
		return
	}

	user := &models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.db.CreateUser(user, req.Password); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		User:  user,
		Token: "dummy-token",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
	}
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "Email and password are required"}`, http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	// For now, return user without token (we'll add JWT later)
	response := models.AuthResponse{
		User:  user,
		Token: "dummy-token", // We'll implement JWT in next iteration
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
	}
}

func (h *AuthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":    "healthy",
		"service":   "authservice",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
	}
}