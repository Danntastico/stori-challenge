package repository

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestWithActualData tests the repository with the real transactions.json file
func TestWithActualData(t *testing.T) {
	// Try to load the actual data file
	// Path relative to the repository package
	dataPath := filepath.Join("..", "..", "data", "transactions.json")

	data, err := os.ReadFile(dataPath)
	if err != nil {
		t.Skipf("Skipping integration test: could not read data file: %v", err)
		return
	}

	repo, err := NewJSONRepository(data)
	if err != nil {
		t.Fatalf("Failed to create repository from actual data: %v", err)
	}

	t.Run("load all transactions", func(t *testing.T) {
		transactions, err := repo.GetAll()
		if err != nil {
			t.Fatalf("GetAll() error = %v", err)
		}

		// The actual file has 112 transactions
		expectedCount := 112
		if len(transactions) != expectedCount {
			t.Errorf("Expected %d transactions, got %d", expectedCount, len(transactions))
		}
	})

	t.Run("verify date range", func(t *testing.T) {
		start, end, err := repo.GetDateRange()
		if err != nil {
			t.Fatalf("GetDateRange() error = %v", err)
		}

		// Expected range: 2024-01-01 to 2024-10-28
		expectedStart := "2024-01-01"
		expectedEnd := "2024-10-28"

		if start.Format("2006-01-02") != expectedStart {
			t.Errorf("Start date = %v, want %v", start.Format("2006-01-02"), expectedStart)
		}

		if end.Format("2006-01-02") != expectedEnd {
			t.Errorf("End date = %v, want %v", end.Format("2006-01-02"), expectedEnd)
		}
	})

	t.Run("get income transactions", func(t *testing.T) {
		income, err := repo.GetByType("income")
		if err != nil {
			t.Fatalf("GetByType() error = %v", err)
		}

		// Should have 20 bi-weekly salary payments (10 months * 2)
		expectedCount := 20
		if len(income) != expectedCount {
			t.Errorf("Expected %d income transactions, got %d", expectedCount, len(income))
		}

		// Verify all are salary category
		for _, tx := range income {
			if tx.Category != "salary" {
				t.Errorf("Unexpected income category: %s", tx.Category)
			}
			if tx.Amount <= 0 {
				t.Errorf("Income amount should be positive, got %f", tx.Amount)
			}
		}
	})

	t.Run("get expense transactions", func(t *testing.T) {
		expenses, err := repo.GetByType("expense")
		if err != nil {
			t.Fatalf("GetByType() error = %v", err)
		}

		// Should have 94 expense transactions (112 total - 20 income)
		expectedCount := 92
		if len(expenses) != expectedCount {
			t.Errorf("Expected %d expense transactions, got %d", expectedCount, len(expenses))
		}

		// Verify all have negative amounts
		for _, tx := range expenses {
			if tx.Amount >= 0 {
				t.Errorf("Expense amount should be negative, got %f for %s", tx.Amount, tx.Description)
			}
		}
	})

	t.Run("get transactions by category", func(t *testing.T) {
		categories := []string{"rent", "groceries", "utilities", "dining", "transportation"}

		for _, category := range categories {
			transactions, err := repo.GetByCategory(category)
			if err != nil {
				t.Errorf("GetByCategory(%s) error = %v", category, err)
				continue
			}

			if len(transactions) == 0 {
				t.Errorf("Expected transactions for category %s", category)
			}

			// Verify all transactions match the category
			for _, tx := range transactions {
				if tx.Category != category {
					t.Errorf("Expected category %s, got %s", category, tx.Category)
				}
			}
		}
	})

	t.Run("get transactions by date range", func(t *testing.T) {
		// Get January 2024 transactions
		start, _ := time.Parse("2006-01-02", "2024-01-01")
		end, _ := time.Parse("2006-01-02", "2024-01-31")

		transactions, err := repo.GetByDateRange(start, end)
		if err != nil {
			t.Fatalf("GetByDateRange() error = %v", err)
		}

		// January should have multiple transactions
		if len(transactions) == 0 {
			t.Error("Expected transactions in January 2024")
		}

		// Verify all transactions are in January
		for _, tx := range transactions {
			txDate, _ := tx.ParseDate()
			if txDate.Month() != time.January || txDate.Year() != 2024 {
				t.Errorf("Transaction date %s is not in January 2024", tx.Date)
			}
		}
	})

	t.Run("repository count", func(t *testing.T) {
		count := repo.Count()
		expectedCount := 112

		if count != expectedCount {
			t.Errorf("Count() = %d, want %d", count, expectedCount)
		}
	})
}

