package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danntastico/stori-backend/internal/service"
)

// AdviceHandler handles AI financial advice requests
type AdviceHandler struct {
	analyticsService *service.AnalyticsService
	aiService        *service.AIService
}

// NewAdviceHandler creates a new advice handler
func NewAdviceHandler(analyticsService *service.AnalyticsService, aiService *service.AIService) *AdviceHandler {
	return &AdviceHandler{
		analyticsService: analyticsService,
		aiService:        aiService,
	}
}

// GetAdvice handles POST /api/advice requests
func (h *AdviceHandler) GetAdvice(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req service.AdviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set default context if not provided
	if req.Context == "" {
		req.Context = "general"
	}

	// Get category summary for AI context
	summary, err := h.analyticsService.GetCategorySummary()
	if err != nil {
		log.Printf("Error getting category summary for AI: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to analyze financial data")
		return
	}

	// Generate AI advice (dereference pointer)
	advice, err := h.aiService.GetFinancialAdvice(r.Context(), *summary, req)
	if err != nil {
		log.Printf("Error generating AI advice: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate advice")
		return
	}

	respondWithJSON(w, http.StatusOK, advice)
}

