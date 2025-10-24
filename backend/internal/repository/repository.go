package repository

import (
	"time"

	"github.com/danntastico/stori-backend/internal/domain"
)

// TransactionRepository defines the interface for transaction data access
// This abstraction allows us to swap implementations (JSON -> Database) without
// changing the service or handler layers.
type TransactionRepository interface {
	// GetAll returns all transactions from the data source
	GetAll() ([]domain.Transaction, error)

	// GetByDateRange returns transactions within the specified date range (inclusive)
	// Returns ErrInvalidDateRange if start is after end
	// Returns ErrNoTransactions if no transactions found in range
	GetByDateRange(start, end time.Time) ([]domain.Transaction, error)

	// GetByType returns all transactions of a specific type ("income" or "expense")
	GetByType(txType string) ([]domain.Transaction, error)

	// GetByCategory returns all transactions for a specific category
	GetByCategory(category string) ([]domain.Transaction, error)

	// Future methods for write operations (Phase 2):
	// Create(tx domain.Transaction) error
	// Update(id string, tx domain.Transaction) error
	// Delete(id string) error
}

