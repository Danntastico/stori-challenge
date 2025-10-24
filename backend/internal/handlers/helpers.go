package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/danntastico/stori-backend/internal/domain"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// respondWithJSON sends a JSON response with the given status code
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, log it but don't try to send another response
		// as headers have already been written
		return
	}
}

// respondWithError sends an error response with the given status code and message
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(response)
}

// handleServiceError maps domain errors to HTTP status codes and sends appropriate responses
func handleServiceError(w http.ResponseWriter, err error) {
	// Map domain errors to HTTP status codes
	switch {
	case errors.Is(err, domain.ErrNoTransactions):
		// Return 200 with empty data structure rather than 404
		// This is more RESTful for "no results found" scenarios
		respondWithError(w, http.StatusOK, "No transactions found")

	case errors.Is(err, domain.ErrInvalidDateRange):
		respondWithError(w, http.StatusBadRequest, "Invalid date range: start date must be before end date")

	case errors.Is(err, domain.ErrInvalidDate):
		respondWithError(w, http.StatusBadRequest, "Invalid date format, expected YYYY-MM-DD")

	case errors.Is(err, domain.ErrInvalidCategory):
		respondWithError(w, http.StatusBadRequest, "Category cannot be empty")

	case errors.Is(err, domain.ErrInvalidType):
		respondWithError(w, http.StatusBadRequest, "Type must be either 'income' or 'expense'")

	case errors.Is(err, domain.ErrInvalidAmount):
		respondWithError(w, http.StatusBadRequest, "Amount sign must match transaction type")

	default:
		// Unknown error - return 500 Internal Server Error
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}

