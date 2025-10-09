package main

import (
	"log"
	"net/http"
	"os"
	"auth-service/internal/database"
	"auth-service/internal/handlers"
	"github.com/gorilla/mux"
)


func main() {
	// Initialize database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("❌ Auth Service: Database connection failed:", err)
	}
	defer db.Close()

	// Initialize database schema
	if err := db.Init(); err != nil {
		log.Fatal("❌ Auth Service: Database initialization failed:", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)

	// Setup routes
	r := mux.NewRouter()
	
	
	// Log all incoming requests
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("📥 AUTH SERVICE - INCOMING REQUEST: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})
	
	// API routes
	r.HandleFunc("/api/v1/auth/signup", authHandler.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/health", authHandler.HealthCheck).Methods("GET")

	// Log all registered routes
	log.Println("📋 AUTH SERVICE - REGISTERED ROUTES:")
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
		port = "8084"
	}

	log.Printf("🚀 Auth Service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}