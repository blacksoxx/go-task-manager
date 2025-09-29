package main

import (
    "log"
    "net/http"
    "os"
    "userservice/internal/database"
    "userservice/internal/handlers"

    "github.com/gorilla/mux"
)

func main() {
    // Initialize database
    db, err := database.NewPostgresDB()
    if err != nil {
        log.Fatal("❌ Database connection failed:", err)
    }
    defer db.Close()

    // Initialize database schema
    if err := db.Init(); err != nil {
        log.Fatal("❌ Database initialization failed:", err)
    }

    // Initialize handlers
    userHandler := handlers.NewUserHandler(db)

    // Setup routes
    r := mux.NewRouter()
    
    // API routes
    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
    api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
    api.HandleFunc("/auth/login", userHandler.Login).Methods("POST")
    
    // Health check
    r.HandleFunc("/health", userHandler.HealthCheck).Methods("GET")
    
    // Middleware
    r.Use(loggingMiddleware)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }

    log.Printf("🚀 User service starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
        next.ServeHTTP(w, r)
    })
}