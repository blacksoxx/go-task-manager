package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"notification-service/internal/database"
	"notification-service/internal/handlers"
	"time"

	"github.com/gorilla/mux"
)

// CORS middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("❌ Notification Service: Database connection failed:", err)
	}
	defer db.Close()

	// Initialize database schema
	if err := db.Init(); err != nil {
		log.Fatal("❌ Notification Service: Database initialization failed:", err)
	}

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(db)

	// Setup routes
	r := mux.NewRouter()
	
	// Add CORS middleware
	r.Use(enableCORS)
	
	// Log all incoming requests
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("📥 NOTIFICATION SERVICE - INCOMING REQUEST: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})
	
	// API routes
	r.HandleFunc("/api/v1/notifications", notificationHandler.CreateNotification).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/notifications/{id}", notificationHandler.GetNotification).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/notifications/{id}", notificationHandler.DeleteNotification).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/v1/notifications/{id}/read", notificationHandler.MarkAsRead).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/v1/notifications/{id}/status", notificationHandler.UpdateStatus).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/v1/users/{user_id}/notifications", notificationHandler.GetUserNotifications).Methods("GET", "OPTIONS")
	
	// Enhanced health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status":    "healthy",
			"service":   "notificationservice",
			"timestamp": time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
		}
	}).Methods("GET")

	// Log all registered routes
	log.Println("📋 NOTIFICATION SERVICE - REGISTERED ROUTES:")
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
		port = "8083"
	}

	log.Printf("🚀 Notification Service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}