package router

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    "time"
    "apigateway/internal/models"
)

type ServiceRouter struct {
    services map[string]models.ServiceConfig
}

func NewServiceRouter(services map[string]models.ServiceConfig) *ServiceRouter {
    return &ServiceRouter{
        services: services,
    }
}

func (sr *ServiceRouter) RouteRequest(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    log.Printf("🔀 Gateway Routing: %s %s", r.Method, path)
    log.Printf("  📍 Incoming path: '%s'", path)    
    switch {
    case strings.HasPrefix(path, "/api/v1/auth/signup"):
        log.Printf("  → Routing SIGNUP to Auth Service") 
        sr.proxyRequest(w, r, "auth-service")
    case strings.HasPrefix(path, "/api/v1/auth/login"):
        log.Printf("  → Routing LOGIN to Auth Service") 
        sr.proxyRequest(w, r, "auth-service")
    case strings.HasPrefix(path, "/api/v1/auth"):
        log.Printf("  → Routing AUTH to Auth Service") 
        sr.proxyRequest(w, r, "auth-service")
    case strings.Contains(path, "/users/") && strings.Contains(path, "/tasks"):
        log.Printf("  → Routing USER TASKS to Task Service")
        sr.proxyRequest(w, r, "task-service")
    case strings.HasPrefix(path, "/api/v1/users"):
        log.Printf("  → Routing USERS to User Service") 
        sr.proxyRequest(w, r, "user-service")
    case strings.HasPrefix(path, "/api/v1/tasks"):
        log.Printf("  → Routing TASKS to Task Service")
        sr.proxyRequest(w, r, "task-service")
    case strings.HasPrefix(path, "/api/v1/notifications"):
        log.Printf("  → Routing NOTIFICATIONS to Notification Service")
        sr.proxyRequest(w, r, "notification-service")
    case path == "/health":
        sr.HealthCheck(w, r)
    default:
        log.Printf("  ❌ No route found - returning 404")
        http.NotFound(w, r)
    }
}

func (sr *ServiceRouter) proxyRequest(w http.ResponseWriter, r *http.Request, serviceName string) {
    service, exists := sr.services[serviceName]
    if !exists {
        log.Printf("  ❌ Service %s not found in configuration", serviceName)
        http.Error(w, fmt.Sprintf("Service %s not found", serviceName), http.StatusBadGateway)
        return
    }

    targetURL, err := url.Parse(service.URL)
    if err != nil {
        log.Printf("  ❌ Invalid service URL: %s", err.Error())
        http.Error(w, fmt.Sprintf("Invalid service URL: %s", err.Error()), http.StatusInternalServerError)
        return
    }

    log.Printf("  → Proxying to: %s%s", targetURL.String(), r.URL.Path)
    
    proxy := httputil.NewSingleHostReverseProxy(targetURL)
    
    // Modify the request to preserve the original path
    r.URL.Host = targetURL.Host
    r.URL.Scheme = targetURL.Scheme
    r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
    r.Host = targetURL.Host

    proxy.ServeHTTP(w, r)
}

func (sr *ServiceRouter) HealthCheck(w http.ResponseWriter, r *http.Request) {
    healthStatus := make(map[string]interface{})
    
    for serviceName, service := range sr.services {
        client := &http.Client{Timeout: 3 * time.Second}
        
        // Try to check the health endpoint
        healthURL := service.URL + "/health"
        resp, err := client.Get(healthURL)
        
        if err != nil {
            log.Printf("  ❌ Health check failed for %s: %v", serviceName, err)
            healthStatus[serviceName] = map[string]string{
                "status": "unhealthy",
                "error":  err.Error(),
            }
        } else if resp.StatusCode != http.StatusOK {
            log.Printf("  ❌ Health check failed for %s: status %d", serviceName, resp.StatusCode)
            healthStatus[serviceName] = map[string]string{
                "status": "unhealthy", 
                "error":  fmt.Sprintf("Status %d", resp.StatusCode),
            }
        } else {
            log.Printf("  ✅ Health check passed for %s", serviceName)
            healthStatus[serviceName] = map[string]string{
                "status": "healthy",
                "url":    service.URL,
            }
        }
        
        if resp != nil {
            resp.Body.Close()
        }
    }

    response := models.HealthResponse{
        Status:    "healthy",
        Service:   "api-gateway",
        Timestamp: time.Now().Format(time.RFC3339),
        Services:  healthStatus,
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
    }
}