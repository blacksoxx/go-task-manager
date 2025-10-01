package models

type ServiceConfig struct {
    Name string `json:"name"`
    URL  string `json:"url"`
    Port string `json:"port"`
}

type GatewayConfig struct {
    Port     string                   `json:"port"`
    Services map[string]ServiceConfig `json:"services"`
}

type HealthResponse struct {
    Status    string            `json:"status"`
    Service   string            `json:"service"`
    Timestamp string            `json:"timestamp"`
    Services  map[string]string `json:"services,omitempty"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}
