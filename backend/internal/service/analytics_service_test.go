package service

import (
	"testing"
	"time"

	"github.com/danntastico/stori-backend/internal/domain"
	"github.com/danntastico/stori-backend/internal/repository"
)

// Test data
var testTransactionsJSON = []byte(`[
	{"date": "2024-01-01", "amount": 2800, "category": "salary", "description": "Bi-weekly salary", "type": "income"},
	{"date": "2024-01-02", "amount": -1200, "category": "rent", "description": "Monthly rent", "type": "expense"},
	{"date": "2024-01-03", "amount": -85, "category": "groceries", "description": "Whole Foods", "type": "expense"},
	{"date": "2024-01-05", "amount": -45, "category": "utilities", "description": "Electric bill", "type": "expense"},
	{"date": "2024-01-16", "amount": 2800, "category": "salary", "description": "Bi-weekly salary", "type": "income"},
	{"date": "2024-02-01", "amount": 2800, "category": "salary", "description": "Bi-weekly salary", "type": "income"},
	{"date": "2024-02-02", "amount": -1200, "category": "rent", "description": "Monthly rent", "type": "expense"},
	{"date": "2024-02-04", "amount": -110, "category": "groceries", "description": "Costco", "type": "expense"}
]`)

func setupTestService(t *testing.T) *AnalyticsService {
	t.Helper()

	repo, err := repository.NewJSONRepository(testTransactionsJSON)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	return NewAnalyticsService(repo)
}

func TestNewAnalyticsService(t *testing.T) {
	repo, _ := repository.NewJSONRepository(testTransactionsJSON)
	service := NewAnalyticsService(repo)

	if service == nil {
		t.Error("NewAnalyticsService() returned nil")
	}

	if service.repo == nil {
		t.Error("AnalyticsService.repo is nil")
	}
}

func TestAnalyticsService_GetCategorySummary(t *testing.T) {
	service := setupTestService(t)

	summary, err := service.GetCategorySummary()
	if err != nil {
		t.Fatalf("GetCategorySummary() error = %v", err)
	}

	// Verify income categories
	if len(summary.Income) != 1 {
		t.Errorf("Expected 1 income category, got %d", len(summary.Income))
	}

	salary, exists := summary.Income["salary"]
	if !exists {
		t.Fatal("Expected salary income category")
	}

	// 3 salary transactions of 2800 each = 8400
	expectedSalaryTotal := 8400.0
	if salary.Total != expectedSalaryTotal {
		t.Errorf("Salary total = %v, want %v", salary.Total, expectedSalaryTotal)
	}

	if salary.Count != 3 {
		t.Errorf("Salary count = %d, want 3", salary.Count)
	}

	if salary.Percentage != 100.0 {
		t.Errorf("Salary percentage = %v, want 100.0 (only income category)", salary.Percentage)
	}

	// Verify expense categories
	if len(summary.Expenses) != 3 {
		t.Errorf("Expected 3 expense categories, got %d", len(summary.Expenses))
	}

	// Check rent category
	rent, exists := summary.Expenses["rent"]
	if !exists {
		t.Error("Expected rent expense category")
	}
	if rent.Count != 2 {
		t.Errorf("Rent count = %d, want 2", rent.Count)
	}

	// Check groceries category
	groceries, exists := summary.Expenses["groceries"]
	if !exists {
		t.Error("Expected groceries expense category")
	}
	if groceries.Count != 2 {
		t.Errorf("Groceries count = %d, want 2", groceries.Count)
	}

	// Verify financial summary
	if summary.Summary.TotalIncome != 8400.0 {
		t.Errorf("TotalIncome = %v, want 8400.0", summary.Summary.TotalIncome)
	}

	// Total expenses: 1200 + 85 + 45 + 1200 + 110 = 2640
	expectedExpenses := 2640.0
	if summary.Summary.TotalExpenses != expectedExpenses {
		t.Errorf("TotalExpenses = %v, want %v", summary.Summary.TotalExpenses, expectedExpenses)
	}

	expectedSavings := 8400.0 - 2640.0 // 5760
	if summary.Summary.NetSavings != expectedSavings {
		t.Errorf("NetSavings = %v, want %v", summary.Summary.NetSavings, expectedSavings)
	}

	// Savings rate: (5760 / 8400) * 100 = 68.57%
	expectedSavingsRate := 68.57
	if summary.Summary.SavingsRate != expectedSavingsRate {
		t.Errorf("SavingsRate = %v, want %v", summary.Summary.SavingsRate, expectedSavingsRate)
	}

	// Verify period
	if summary.Period.Start != "2024-01-01" {
		t.Errorf("Period start = %v, want 2024-01-01", summary.Period.Start)
	}

	if summary.Period.End != "2024-02-04" {
		t.Errorf("Period end = %v, want 2024-02-04", summary.Period.End)
	}

	if summary.Period.Months != 2 {
		t.Errorf("Period months = %d, want 2", summary.Period.Months)
	}
}

func TestAnalyticsService_GetTimeline(t *testing.T) {
	service := setupTestService(t)

	timeline, err := service.GetTimeline()
	if err != nil {
		t.Fatalf("GetTimeline() error = %v", err)
	}

	if timeline.Aggregation != "monthly" {
		t.Errorf("Aggregation = %v, want monthly", timeline.Aggregation)
	}

	if len(timeline.Timeline) != 2 {
		t.Fatalf("Expected 2 months in timeline, got %d", len(timeline.Timeline))
	}

	// Check January (2024-01)
	jan := timeline.Timeline[0]
	if jan.Period != "2024-01" {
		t.Errorf("First period = %v, want 2024-01", jan.Period)
	}

	// January income: 2800 + 2800 = 5600
	expectedJanIncome := 5600.0
	if jan.Income != expectedJanIncome {
		t.Errorf("January income = %v, want %v", jan.Income, expectedJanIncome)
	}

	// January expenses: 1200 + 85 + 45 = 1330
	expectedJanExpenses := 1330.0
	if jan.Expenses != expectedJanExpenses {
		t.Errorf("January expenses = %v, want %v", jan.Expenses, expectedJanExpenses)
	}

	expectedJanNet := 5600.0 - 1330.0 // 4270
	if jan.Net != expectedJanNet {
		t.Errorf("January net = %v, want %v", jan.Net, expectedJanNet)
	}

	// Check February (2024-02)
	feb := timeline.Timeline[1]
	if feb.Period != "2024-02" {
		t.Errorf("Second period = %v, want 2024-02", feb.Period)
	}

	// February income: 2800
	if feb.Income != 2800.0 {
		t.Errorf("February income = %v, want 2800.0", feb.Income)
	}

	// February expenses: 1200 + 110 = 1310
	expectedFebExpenses := 1310.0
	if feb.Expenses != expectedFebExpenses {
		t.Errorf("February expenses = %v, want %v", feb.Expenses, expectedFebExpenses)
	}

	expectedFebNet := 2800.0 - 1310.0 // 1490
	if feb.Net != expectedFebNet {
		t.Errorf("February net = %v, want %v", feb.Net, expectedFebNet)
	}

	// Verify timeline is sorted chronologically
	if timeline.Timeline[0].Period > timeline.Timeline[1].Period {
		t.Error("Timeline is not sorted chronologically")
	}
}

func TestAnalyticsService_GetTransactions(t *testing.T) {
	service := setupTestService(t)

	response, err := service.GetTransactions()
	if err != nil {
		t.Fatalf("GetTransactions() error = %v", err)
	}

	expectedCount := 8
	if response.Count != expectedCount {
		t.Errorf("Count = %d, want %d", response.Count, expectedCount)
	}

	if len(response.Transactions) != expectedCount {
		t.Errorf("Transactions length = %d, want %d", len(response.Transactions), expectedCount)
	}

	if response.Period.Start != "2024-01-01" {
		t.Errorf("Period start = %v, want 2024-01-01", response.Period.Start)
	}

	if response.Period.End != "2024-02-04" {
		t.Errorf("Period end = %v, want 2024-02-04", response.Period.End)
	}
}

func TestAnalyticsService_GetTransactionsByDateRange(t *testing.T) {
	service := setupTestService(t)

	tests := []struct {
		name          string
		start         string
		end           string
		expectedCount int
		wantErr       bool
	}{
		{
			name:          "january only",
			start:         "2024-01-01",
			end:           "2024-01-31",
			expectedCount: 5,
			wantErr:       false,
		},
		{
			name:          "february only",
			start:         "2024-02-01",
			end:           "2024-02-29",
			expectedCount: 3,
			wantErr:       false,
		},
		{
			name:          "single day",
			start:         "2024-01-01",
			end:           "2024-01-01",
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name:          "no transactions",
			start:         "2023-01-01",
			end:           "2023-12-31",
			expectedCount: 0,
			wantErr:       true,
		},
		{
			name:          "invalid range",
			start:         "2024-12-31",
			end:           "2024-01-01",
			expectedCount: 0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, _ := time.Parse("2006-01-02", tt.start)
			end, _ := time.Parse("2006-01-02", tt.end)

			response, err := service.GetTransactionsByDateRange(start, end)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionsByDateRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if response.Count != tt.expectedCount {
					t.Errorf("Count = %d, want %d", response.Count, tt.expectedCount)
				}

				if len(response.Transactions) != tt.expectedCount {
					t.Errorf("Transactions length = %d, want %d", len(response.Transactions), tt.expectedCount)
				}

				// Verify all transactions are within range
				for _, tx := range response.Transactions {
					txDate, _ := tx.ParseDate()
					if txDate.Before(start) || txDate.After(end) {
						t.Errorf("Transaction date %s is outside range %s to %s", tx.Date, tt.start, tt.end)
					}
				}
			}
		})
	}
}

func TestAnalyticsService_CalculateMonthsBetween(t *testing.T) {
	service := setupTestService(t)

	tests := []struct {
		name     string
		start    string
		end      string
		expected int
	}{
		{
			name:     "same month",
			start:    "2024-01-01",
			end:      "2024-01-31",
			expected: 1,
		},
		{
			name:     "two months",
			start:    "2024-01-01",
			end:      "2024-02-15",
			expected: 2,
		},
		{
			name:     "full year",
			start:    "2024-01-01",
			end:      "2024-12-31",
			expected: 12,
		},
		{
			name:     "across year boundary",
			start:    "2023-12-01",
			end:      "2024-01-15",
			expected: 2,
		},
		{
			name:     "ten months (Jan to Oct)",
			start:    "2024-01-01",
			end:      "2024-10-28",
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, _ := time.Parse("2006-01-02", tt.start)
			end, _ := time.Parse("2006-01-02", tt.end)

			result := service.calculateMonthsBetween(start, end)

			if result != tt.expected {
				t.Errorf("calculateMonthsBetween() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestAnalyticsService_RoundingAccuracy(t *testing.T) {
	service := setupTestService(t)

	summary, err := service.GetCategorySummary()
	if err != nil {
		t.Fatalf("GetCategorySummary() error = %v", err)
	}

	// Verify all monetary values are rounded to 2 decimal places
	checkRounding := func(val float64, name string) {
		rounded := roundToTwo(val)
		if rounded != val {
			t.Errorf("%s value %v is not rounded to 2 decimal places", name, val)
		}
	}

	checkRounding(summary.Summary.TotalIncome, "TotalIncome")
	checkRounding(summary.Summary.TotalExpenses, "TotalExpenses")
	checkRounding(summary.Summary.NetSavings, "NetSavings")
	checkRounding(summary.Summary.SavingsRate, "SavingsRate")

	for category, detail := range summary.Income {
		checkRounding(detail.Total, "Income."+category+".Total")
		checkRounding(detail.Percentage, "Income."+category+".Percentage")
	}

	for category, detail := range summary.Expenses {
		checkRounding(detail.Total, "Expenses."+category+".Total")
		checkRounding(detail.Percentage, "Expenses."+category+".Percentage")
	}
}

func TestAnalyticsService_EmptyData(t *testing.T) {
	emptyJSON := []byte(`[]`)
	repo, err := repository.NewJSONRepository(emptyJSON)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	service := NewAnalyticsService(repo)

	t.Run("GetCategorySummary with empty data", func(t *testing.T) {
		_, err := service.GetCategorySummary()
		if err != domain.ErrNoTransactions {
			t.Errorf("Expected ErrNoTransactions, got %v", err)
		}
	})

	t.Run("GetTimeline with empty data", func(t *testing.T) {
		_, err := service.GetTimeline()
		if err != domain.ErrNoTransactions {
			t.Errorf("Expected ErrNoTransactions, got %v", err)
		}
	})

	t.Run("GetTransactions with empty data", func(t *testing.T) {
		_, err := service.GetTransactions()
		if err != domain.ErrNoTransactions {
			t.Errorf("Expected ErrNoTransactions, got %v", err)
		}
	})
}

