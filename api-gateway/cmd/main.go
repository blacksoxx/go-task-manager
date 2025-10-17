package main

import (
    "log"
    "net/http"
    "apigateway/internal/config"
    "apigateway/internal/middleware"
    "apigateway/internal/router"
    "apigateway/internal/jwt"       

)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize JWT service 
    jwtService := jwt.NewJWTService()
    
    // Create service router
    serviceRouter := router.NewServiceRouter(cfg.Services)
    
    // Setup routes
    mux := http.NewServeMux()
    
    // Health check endpoint
    mux.HandleFunc("/health", serviceRouter.HealthCheck)
    
    // API routes - all other requests go through the router
    mux.HandleFunc("/", serviceRouter.RouteRequest)
    
    // Apply middleware
    handler := middleware.CorsMiddleware(mux)
    handler = middleware.LoggingMiddleware(handler)
    handler = middleware.AuthMiddleware(jwtService)(handler) // ADD AUTH MIDDLEWARE

    
    port := cfg.Port
    log.Printf("🚀 API Gateway starting on port %s", port)
    log.Printf("📡 Routing to services:")
    for name, service := range cfg.Services {
        log.Printf("   - %s: %s", name, service.URL)
    }
    log.Printf("🔐 JWT Authentication: ENABLED")
    log.Fatal(http.ListenAndServe(":"+port, handler))
}
