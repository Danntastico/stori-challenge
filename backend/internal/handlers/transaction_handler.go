package handlers

import (
	"net/http"
	"time"

	"github.com/danntastico/stori-backend/internal/domain"
	"github.com/danntastico/stori-backend/internal/service"
)

// TransactionHandler handles transaction-related requests
type TransactionHandler struct {
	analyticsService service.AnalyticsServiceInterface
}

// NewTransactionHandler creates a new transaction handler
// Now accepts interface instead of concrete type for better testability
func NewTransactionHandler(analyticsService service.AnalyticsServiceInterface) *TransactionHandler {
	return &TransactionHandler{
		analyticsService: analyticsService,
	}
}

// ServeHTTP handles GET /api/transactions
// Query parameters:
//   - startDate: ISO 8601 date (YYYY-MM-DD) - optional
//   - endDate: ISO 8601 date (YYYY-MM-DD) - optional
//   - type: "income" or "expense" - optional (future use)
//   - category: category name - optional (future use)
func (h *TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	startDateStr := query.Get("startDate")
	endDateStr := query.Get("endDate")

	var response *domain.TransactionsResponse
	var err error

	// If date range provided, filter by date range
	if startDateStr != "" && endDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid startDate format, expected YYYY-MM-DD")
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid endDate format, expected YYYY-MM-DD")
			return
		}

		response, _ = h.analyticsService.GetTransactionsByDateRange(startDate, endDate)
	} else {
		// Get all transactions
		response, err = h.analyticsService.GetTransactions()
	}

	// Handle errors
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Send successful response
	respondWithJSON(w, http.StatusOK, response)
}

