package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/danntastico/stori-backend/internal/domain"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health check handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// ServeHTTP handles GET /api/health
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := domain.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

