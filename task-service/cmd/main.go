package main

import (
    "log"
    "net/http"
    "os"
    "taskservice/internal/database"
    "taskservice/internal/handlers"

    "github.com/gorilla/mux"
)

func main() {
    // Initialize database
    db, err := database.NewPostgresDB()
    if err != nil {
        log.Fatal("❌ Task Service: Database connection failed:", err)
    }
    defer db.Close()

    // Initialize database schema
    if err := db.Init(); err != nil {
        log.Fatal("❌ Task Service: Database initialization failed:", err)
    }

    // Initialize handlers
    taskHandler := handlers.NewTaskHandler(db)

    // Setup routes
    r := mux.NewRouter()
    
    // API routes
    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
    api.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
    api.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
    api.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")
    api.HandleFunc("/users/{user_id}/tasks", taskHandler.GetUserTasks).Methods("GET")
    
    // Health check
    r.HandleFunc("/health", taskHandler.HealthCheck).Methods("GET")
    
    // Middleware
    r.Use(loggingMiddleware)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8082"
    }

    log.Printf("🚀 Task Service starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Task Service: %s %s %s", r.RemoteAddr, r.Method, r.URL)
        next.ServeHTTP(w, r)
    })
}
