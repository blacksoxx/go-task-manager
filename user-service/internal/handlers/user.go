package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "userservice/internal/database"
    "userservice/internal/models"

    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
    db *database.DB
}

func NewUserHandler(db *database.DB) *UserHandler {
    return &UserHandler{db: db}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req models.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
        return
    }

    // Validate required fields
    if req.Email == "" || req.FirstName == "" || req.LastName == "" || req.Password == "" {
        http.Error(w, `{"error": "All fields are required"}`, http.StatusBadRequest)
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, `{"error": "Error processing password"}`, http.StatusInternalServerError)
        return
    }

    user := &models.User{
        Email:     req.Email,
        FirstName: req.FirstName,
        LastName:  req.LastName,
    }

    if err := h.db.CreateUser(user, string(hashedPassword)); err != nil {
        // Check if it's a duplicate email error
        if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
            http.Error(w, `{"error": "Email already exists"}`, http.StatusConflict)
            return
        }
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(user); err != nil {
        http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
    }
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["id"]

    user, err := h.db.GetUserByID(userID)
    if err != nil {
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user); err != nil {
        http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
    }
}

func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{
        "status":    "healthy",
        "service":   "user-service",
        "timestamp": time.Now().Format(time.RFC3339),
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
    }
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
        return
    }

    user, passwordHash, err := h.db.GetUserByEmail(req.Email)
    if err != nil {
        http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
        return
    }

    // Compare password
    if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
        http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
        return
    }

    // For now, return a simple token. We'll implement JWT properly later.
    response := models.LoginResponse{
        User:  user,
        Token: "mock-jwt-token-for-now", // We'll implement proper JWT in next steps
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
    }
}
