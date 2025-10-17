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
    
    
    // Log all incoming requests
    r.Use(loggingMiddleware)
    
    // API routes
    r.HandleFunc("/api/v1/tasks", taskHandler.CreateTask).Methods("POST")
    r.HandleFunc("/api/v1/tasks/{id}", taskHandler.GetTask).Methods("GET")
    r.HandleFunc("/api/v1/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")
    r.HandleFunc("/api/v1/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")
    r.HandleFunc("/api/v1/users/{user_id}/tasks", taskHandler.GetUserTasks).Methods("GET")
    r.HandleFunc("/health", taskHandler.HealthCheck).Methods("GET")

    // Handle preflight OPTIONS requests for all routes
    r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    // Log all registered routes
    log.Println("📋 REGISTERED ROUTES:")
    r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        t, err := route.GetPathTemplate()
        if err == nil {
            methods, _ := route.GetMethods()
            log.Printf("   %s %s", methods, t)
        }
        return nil
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8082"
    }

    log.Printf("🚀 Task Service starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
        next.ServeHTTP(w, r)
    })
}