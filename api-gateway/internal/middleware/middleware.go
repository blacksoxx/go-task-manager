package middleware

import (
    "log"
    "net/http"
    "time"
    "strings"
    "context"
    "apigateway/internal/jwt"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create a custom ResponseWriter to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start)
        log.Printf("API Gateway: %s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, duration)
    })
}

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "content-type, authorization")   
              
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// Custom response writer to capture status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}



func AuthMiddleware(jwtService *jwt.JWTService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Skip auth for public routes
            if isPublicRoute(r.URL.Path) {
                next.ServeHTTP(w, r)
                return
            }

            // Get token from header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
                return
            }

            // Extract token (format: "Bearer <token>")
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, `{"error": "Invalid authorization format. Use: Bearer <token>"}`, http.StatusUnauthorized)
                return
            }

            token := parts[1]

            // Validate JWT token
            claims, err := jwtService.ValidateToken(token)
            if err != nil {
                http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
                return
            }

            // Add user info to context for downstream services
            ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
            ctx = context.WithValue(ctx, "user_email", claims.Email)

            // Add user info to headers for downstream services
            r.Header.Set("X-User-ID", claims.UserID)
            r.Header.Set("X-User-Email", claims.Email)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func isPublicRoute(path string) bool {
    publicRoutes := []string{
        "/api/v1/auth/login",
        "/api/v1/auth/signup", 
        "/health",
        "/api/v1/auth",
    }
    
    for _, route := range publicRoutes {
        if strings.HasPrefix(path, route) {
            return true
        }
    }
    return false
}