package handlers

import (
	"net/http"

	"github.com/danntastico/stori-backend/internal/service"
)

// SummaryHandler handles financial summary requests
type SummaryHandler struct {
	analyticsService *service.AnalyticsService
}

// NewSummaryHandler creates a new summary handler
func NewSummaryHandler(analyticsService *service.AnalyticsService) *SummaryHandler {
	return &SummaryHandler{
		analyticsService: analyticsService,
	}
}

// HandleCategorySummary handles GET /api/summary/categories
// Returns aggregated spending breakdown by category with totals and percentages
func (h *SummaryHandler) HandleCategorySummary(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get category summary from analytics service
	summary, err := h.analyticsService.GetCategorySummary()
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Send successful response
	respondWithJSON(w, http.StatusOK, summary)
}

// HandleTimeline handles GET /api/summary/timeline
// Returns monthly income vs expenses over time
func (h *SummaryHandler) HandleTimeline(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get timeline from analytics service
	timeline, err := h.analyticsService.GetTimeline()
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Send successful response
	respondWithJSON(w, http.StatusOK, timeline)
}

