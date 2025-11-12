package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "context"
)

type HealthCheckResponse struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Services  map[string]string `json:"services,omitempty"`
    Errors    map[string]string `json:"errors,omitempty"`
}

// HealthCheckHandler godoc
// @Summary Health check
// @Description Check the health status of the server and its dependencies
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Failure 503 {object} HealthCheckResponse
// @Router /health [get]
func (h *Handler) HealthCheckHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()

        response := HealthCheckResponse{
            Status:    "healthy",
            Timestamp: time.Now(),
            Services:  make(map[string]string),
            Errors:    make(map[string]string),
        }

        // Check PostgreSQL connection
        if h.DB != nil {
            if err := h.DB.PingContext(ctx); err != nil {
                response.Status = "unhealthy"
                response.Errors["postgres"] = "Connection failed: " + err.Error()
                response.Services["postgres"] = "down"
            } else {
                response.Services["postgres"] = "up"
            }
        } else {
            response.Status = "unhealthy"
            response.Errors["postgres"] = "Database instance is nil"
            response.Services["postgres"] = "down"
        }

        // Check Redis connection
        if h.Redis != nil {
            if err := h.Redis.Ping(ctx).Err(); err != nil {
                response.Status = "unhealthy"
                response.Errors["redis"] = "Connection failed: " + err.Error()
                response.Services["redis"] = "down"
            } else {
                response.Services["redis"] = "up"
            }
        } else {
            response.Status = "unhealthy"
            response.Errors["redis"] = "Redis client is nil"
            response.Services["redis"] = "down"
        }

        if response.Status == "unhealthy" {
            w.WriteHeader(http.StatusServiceUnavailable)
        } else {
            w.WriteHeader(http.StatusOK)
        }

        json.NewEncoder(w).Encode(response)
    }
}