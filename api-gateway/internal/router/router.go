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
    
    // Route based on URL path
    switch {
    case strings.HasPrefix(path, "/api/v1/users") && strings.HasSuffix(path, "/tasks"):
        // This is /api/v1/users/{user_id}/tasks - route to Task Service
        log.Printf("  → Routing USER TASKS to Task Service")
        sr.proxyRequest(w, r, "task-service")
    case strings.HasPrefix(path, "/api/v1/users"), strings.HasPrefix(path, "/api/v1/auth"):
        log.Printf("  → Routing to User Service")
        sr.proxyRequest(w, r, "user-service")
    case strings.HasPrefix(path, "/api/v1/tasks"):
        log.Printf("  → Routing to Task Service")
        sr.proxyRequest(w, r, "task-service")
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
    healthStatus := make(map[string]string)
    
    for serviceName, service := range sr.services {
        client := &http.Client{Timeout: 2 * time.Second}
        resp, err := client.Get(service.URL + "/health")
        if err != nil || resp.StatusCode != http.StatusOK {
            healthStatus[serviceName] = "unhealthy"
        } else {
            healthStatus[serviceName] = "healthy"
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
