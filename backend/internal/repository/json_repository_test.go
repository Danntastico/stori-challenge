package repository

import (
	"testing"
	"time"

	"github.com/danntastico/stori-backend/internal/domain"
)

// Sample test data
var testJSON = []byte(`[
	{"date": "2024-01-01", "amount": 2800, "category": "salary", "description": "Bi-weekly salary", "type": "income"},
	{"date": "2024-01-02", "amount": -1200, "category": "rent", "description": "Monthly rent", "type": "expense"},
	{"date": "2024-01-03", "amount": -85, "category": "groceries", "description": "Whole Foods", "type": "expense"},
	{"date": "2024-02-01", "amount": 2800, "category": "salary", "description": "Bi-weekly salary", "type": "income"},
	{"date": "2024-02-02", "amount": -1200, "category": "rent", "description": "Monthly rent", "type": "expense"}
]`)

func TestNewJSONRepository(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "valid JSON",
			data:    testJSON,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			data:    []byte(`invalid json`),
			wantErr: true,
		},
		{
			name:    "empty JSON array",
			data:    []byte(`[]`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := NewJSONRepository(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJSONRepository() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && repo == nil {
				t.Error("Expected repository to be non-nil")
			}
		})
	}
}

func TestJSONRepository_GetAll(t *testing.T) {
	repo, err := NewJSONRepository(testJSON)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	transactions, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v", err)
	}

	expectedCount := 5
	if len(transactions) != expectedCount {
		t.Errorf("GetAll() returned %d transactions, want %d", len(transactions), expectedCount)
	}

	// Test that modifications don't affect repository
	transactions[0].Amount = 9999
	checkTransactions, _ := repo.GetAll()
	if checkTransactions[0].Amount == 9999 {
		t.Error("GetAll() should return a copy, not the original slice")
	}
}

func TestJSONRepository_GetAll_Empty(t *testing.T) {
	repo, err := NewJSONRepository([]byte(`[]`))
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	_, err = repo.GetAll()
	if err != domain.ErrNoTransactions {
		t.Errorf("GetAll() with empty data should return ErrNoTransactions, got %v", err)
	}
}

func TestJSONRepository_GetByDateRange(t *testing.T) {
	repo, err := NewJSONRepository(testJSON)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	tests := []struct {
		name          string
		start         string
		end           string
		expectedCount int
		wantErr       error
	}{
		{
			name:          "january only",
			start:         "2024-01-01",
			end:           "2024-01-31",
			expectedCount: 3,
			wantErr:       nil,
		},
		{
			name:          "february only",
			start:         "2024-02-01",
			end:           "2024-02-29",
			expectedCount: 2,
			wantErr:       nil,
		},
		{
			name:          "all transactions",
			start:         "2024-01-01",
			end:           "2024-12-31",
			expectedCount: 5,
			wantErr:       nil,
		},
		{
			name:          "no transactions in range",
			start:         "2023-01-01",
			end:           "2023-12-31",
			expectedCount: 0,
			wantErr:       domain.ErrNoTransactions,
		},
		{
			name:          "invalid range (start after end)",
			start:         "2024-12-31",
			end:           "2024-01-01",
			expectedCount: 0,
			wantErr:       domain.ErrInvalidDateRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, _ := time.Parse("2006-01-02", tt.start)
			end, _ := time.Parse("2006-01-02", tt.end)

			transactions, err := repo.GetByDateRange(start, end)

			if err != tt.wantErr {
				t.Errorf("GetByDateRange() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && len(transactions) != tt.expectedCount {
				t.Errorf("GetByDateRange() returned %d transactions, want %d", len(transactions), tt.expectedCount)
			}
		})
	}
}

func TestJSONRepository_GetByType(t *testing.T) {
	repo, err := NewJSONRepository(testJSON)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	tests := []struct {
		name          string
		txType        string
		expectedCount int
		wantErr       error
	}{
		{
			name:          "income transactions",
			txType:        "income",
			expectedCount: 2,
			wantErr:       nil,
		},
		{
			name:          "expense transactions",
			txType:        "expense",
			expectedCount: 3,
			wantErr:       nil,
		},
		{
			name:          "invalid type",
			txType:        "transfer",
			expectedCount: 0,
			wantErr:       domain.ErrNoTransactions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactions, err := repo.GetByType(tt.txType)

			if err != tt.wantErr {
				t.Errorf("GetByType() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && len(transactions) != tt.expectedCount {
				t.Errorf("GetByType() returned %d transactions, want %d", len(transactions), tt.expectedCount)
			}

			// Verify all returned transactions have correct type
			if tt.wantErr == nil {
				for _, tx := range transactions {
					if tx.Type != tt.txType {
						t.Errorf("GetByType() returned transaction with type %s, want %s", tx.Type, tt.txType)
					}
				}
			}
		})
	}
}

func TestJSONRepository_GetByCategory(t *testing.T) {
	repo, err := NewJSONRepository(testJSON)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	tests := []struct {
		name          string
		category      string
		expectedCount int
		wantErr       error
	}{
		{
			name:          "salary category",
			category:      "salary",
			expectedCount: 2,
			wantErr:       nil,
		},
		{
			name:          "rent category",
			category:      "rent",
			expectedCount: 2,
			wantErr:       nil,
		},
		{
			name:          "groceries category",
			category:      "groceries",
			expectedCount: 1,
			wantErr:       nil,
		},
		{
			name:          "non-existent category",
			category:      "entertainment",
			expectedCount: 0,
			wantErr:       domain.ErrNoTransactions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactions, err := repo.GetByCategory(tt.category)

			if err != tt.wantErr {
				t.Errorf("GetByCategory() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && len(transactions) != tt.expectedCount {
				t.Errorf("GetByCategory() returned %d transactions, want %d", len(transactions), tt.expectedCount)
			}

			// Verify all returned transactions have correct category
			if tt.wantErr == nil {
				for _, tx := range transactions {
					if tx.Category != tt.category {
						t.Errorf("GetByCategory() returned transaction with category %s, want %s", tx.Category, tt.category)
					}
				}
			}
		})
	}
}

func TestJSONRepository_GetDateRange(t *testing.T) {
	repo, err := NewJSONRepository(testJSON)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	start, end, err := repo.GetDateRange()
	if err != nil {
		t.Errorf("GetDateRange() error = %v", err)
	}

	expectedStart := "2024-01-01"
	expectedEnd := "2024-02-02"

	if start.Format("2006-01-02") != expectedStart {
		t.Errorf("GetDateRange() start = %v, want %v", start.Format("2006-01-02"), expectedStart)
	}

	if end.Format("2006-01-02") != expectedEnd {
		t.Errorf("GetDateRange() end = %v, want %v", end.Format("2006-01-02"), expectedEnd)
	}
}

func TestJSONRepository_GetDateRange_Empty(t *testing.T) {
	repo, err := NewJSONRepository([]byte(`[]`))
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	_, _, err = repo.GetDateRange()
	if err != domain.ErrNoTransactions {
		t.Errorf("GetDateRange() with empty data should return ErrNoTransactions, got %v", err)
	}
}

func TestJSONRepository_Count(t *testing.T) {
	repo, err := NewJSONRepository(testJSON)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	count := repo.Count()
	expectedCount := 5

	if count != expectedCount {
		t.Errorf("Count() = %d, want %d", count, expectedCount)
	}
}

