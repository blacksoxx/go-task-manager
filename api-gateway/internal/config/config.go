package config

import (
    "os"
    "apigateway/internal/models"
)

func LoadConfig() *models.GatewayConfig {
    return &models.GatewayConfig{
        Port: getEnv("PORT", "8080"),
        Services: map[string]models.ServiceConfig{
            "user-service": {
                Name: "user-service",
                URL:  getEnv("USER_SERVICE_URL", "http://localhost:8081"),
                Port: getEnv("USER_SERVICE_PORT", "8081"),
            },
            "task-service": {
                Name: "task-service", 
                URL:  getEnv("TASK_SERVICE_URL", "http://localhost:8082"),
                Port: getEnv("TASK_SERVICE_PORT", "8082"),
            },
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
