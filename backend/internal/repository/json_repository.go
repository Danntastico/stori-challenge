package repository

import (
	"encoding/json"
	"time"

	"github.com/danntastico/stori-backend/internal/domain"
)

// JSONRepository implements TransactionRepository using in-memory JSON data
type JSONRepository struct {
	transactions []domain.Transaction
}

// NewJSONRepository creates a new JSON-based repository from raw JSON data
// This is designed to work with embedded JSON files using go:embed
func NewJSONRepository(data []byte) (*JSONRepository, error) {
	var transactions []domain.Transaction

	if err := json.Unmarshal(data, &transactions); err != nil {
		return nil, err
	}

	// Validate all transactions on load
	for i, tx := range transactions {
		if err := tx.Validate(); err != nil {
			// Note: In production, you might want to log invalid transactions
			// For now, we trust the provided JSON data is valid
			_ = i // Placeholder for future logging
		}
	}

	return &JSONRepository{
		transactions: transactions,
	}, nil
}

// GetAll returns all transactions
func (r *JSONRepository) GetAll() ([]domain.Transaction, error) {
	if len(r.transactions) == 0 {
		return nil, domain.ErrNoTransactions
	}

	// Return a copy to prevent external modifications
	result := make([]domain.Transaction, len(r.transactions))
	copy(result, r.transactions)

	return result, nil
}

// GetByDateRange returns transactions within the specified date range (inclusive)
func (r *JSONRepository) GetByDateRange(start, end time.Time) ([]domain.Transaction, error) {
	// Validate date range
	if start.After(end) {
		return nil, domain.ErrInvalidDateRange
	}

	var filtered []domain.Transaction

	for _, tx := range r.transactions {
		txDate, err := tx.ParseDate()
		if err != nil {
			// Skip transactions with invalid dates
			continue
		}

		// Check if transaction date is within range (inclusive)
		if (txDate.Equal(start) || txDate.After(start)) &&
			(txDate.Equal(end) || txDate.Before(end)) {
			filtered = append(filtered, tx)
		}
	}

	if len(filtered) == 0 {
		return nil, domain.ErrNoTransactions
	}

	return filtered, nil
}

// GetByType returns all transactions of a specific type
func (r *JSONRepository) GetByType(txType string) ([]domain.Transaction, error) {
	var filtered []domain.Transaction

	for _, tx := range r.transactions {
		if tx.Type == txType {
			filtered = append(filtered, tx)
		}
	}

	if len(filtered) == 0 {
		return nil, domain.ErrNoTransactions
	}

	return filtered, nil
}

// GetByCategory returns all transactions for a specific category
func (r *JSONRepository) GetByCategory(category string) ([]domain.Transaction, error) {
	var filtered []domain.Transaction

	for _, tx := range r.transactions {
		if tx.Category == category {
			filtered = append(filtered, tx)
		}
	}

	if len(filtered) == 0 {
		return nil, domain.ErrNoTransactions
	}

	return filtered, nil
}

// Helper methods for analytics (not part of the interface but useful)

// GetDateRange returns the earliest and latest transaction dates
func (r *JSONRepository) GetDateRange() (start, end time.Time, err error) {
	if len(r.transactions) == 0 {
		return time.Time{}, time.Time{}, domain.ErrNoTransactions
	}

	var minDate, maxDate time.Time
	first := true

	for _, tx := range r.transactions {
		txDate, err := tx.ParseDate()
		if err != nil {
			continue
		}

		if first {
			minDate = txDate
			maxDate = txDate
			first = false
			continue
		}

		if txDate.Before(minDate) {
			minDate = txDate
		}
		if txDate.After(maxDate) {
			maxDate = txDate
		}
	}

	if first {
		return time.Time{}, time.Time{}, domain.ErrNoTransactions
	}

	return minDate, maxDate, nil
}

// Count returns the total number of transactions
func (r *JSONRepository) Count() int {
	return len(r.transactions)
}

