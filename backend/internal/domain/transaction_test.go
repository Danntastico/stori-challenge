package domain

import (
	"testing"
)

func TestTransaction_IsIncome(t *testing.T) {
	tx := Transaction{Type: "income"}
	if !tx.IsIncome() {
		t.Error("Expected IsIncome() to return true for income transaction")
	}

	tx.Type = "expense"
	if tx.IsIncome() {
		t.Error("Expected IsIncome() to return false for expense transaction")
	}
}

func TestTransaction_IsExpense(t *testing.T) {
	tx := Transaction{Type: "expense"}
	if !tx.IsExpense() {
		t.Error("Expected IsExpense() to return true for expense transaction")
	}

	tx.Type = "income"
	if tx.IsExpense() {
		t.Error("Expected IsExpense() to return false for income transaction")
	}
}

func TestTransaction_AbsoluteAmount(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected float64
	}{
		{"positive amount", 100.50, 100.50},
		{"negative amount", -100.50, 100.50},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := Transaction{Amount: tt.amount}
			if result := tx.AbsoluteAmount(); result != tt.expected {
				t.Errorf("AbsoluteAmount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTransaction_ParseDate(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		{"valid date", "2024-01-01", false},
		{"valid date 2", "2024-12-31", false},
		{"invalid format", "01-01-2024", true},
		{"invalid date", "2024-13-01", true},
		{"empty date", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := Transaction{Date: tt.date}
			_, err := tx.ParseDate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransaction_GetYearMonth(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		expected string
		wantErr  bool
	}{
		{"january", "2024-01-15", "2024-01", false},
		{"december", "2024-12-31", "2024-12", false},
		{"invalid date", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := Transaction{Date: tt.date}
			result, err := tx.GetYearMonth()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetYearMonth() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("GetYearMonth() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTransaction_Validate(t *testing.T) {
	tests := []struct {
		name        string
		transaction Transaction
		wantErr     error
	}{
		{
			name: "valid income",
			transaction: Transaction{
				Date:     "2024-01-01",
				Amount:   2800,
				Category: "salary",
				Type:     "income",
			},
			wantErr: nil,
		},
		{
			name: "valid expense",
			transaction: Transaction{
				Date:     "2024-01-01",
				Amount:   -1200,
				Category: "rent",
				Type:     "expense",
			},
			wantErr: nil,
		},
		{
			name: "empty date",
			transaction: Transaction{
				Amount:   2800,
				Category: "salary",
				Type:     "income",
			},
			wantErr: ErrInvalidDate,
		},
		{
			name: "invalid date format",
			transaction: Transaction{
				Date:     "01-01-2024",
				Amount:   2800,
				Category: "salary",
				Type:     "income",
			},
			wantErr: ErrInvalidDate,
		},
		{
			name: "empty category",
			transaction: Transaction{
				Date:   "2024-01-01",
				Amount: 2800,
				Type:   "income",
			},
			wantErr: ErrInvalidCategory,
		},
		{
			name: "invalid type",
			transaction: Transaction{
				Date:     "2024-01-01",
				Amount:   2800,
				Category: "salary",
				Type:     "transfer",
			},
			wantErr: ErrInvalidType,
		},
		{
			name: "income with negative amount",
			transaction: Transaction{
				Date:     "2024-01-01",
				Amount:   -2800,
				Category: "salary",
				Type:     "income",
			},
			wantErr: ErrInvalidAmount,
		},
		{
			name: "expense with positive amount",
			transaction: Transaction{
				Date:     "2024-01-01",
				Amount:   1200,
				Category: "rent",
				Type:     "expense",
			},
			wantErr: ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.transaction.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFinancialSummary_CalculateSavingsRate(t *testing.T) {
	tests := []struct {
		name     string
		summary  FinancialSummary
		expected float64
	}{
		{
			name: "positive savings",
			summary: FinancialSummary{
				TotalIncome:   5600,
				TotalExpenses: 4600,
				NetSavings:    1000,
			},
			expected: 17.86, // (1000 / 5600) * 100 = 17.857...
		},
		{
			name: "zero savings",
			summary: FinancialSummary{
				TotalIncome:   5000,
				TotalExpenses: 5000,
				NetSavings:    0,
			},
			expected: 0,
		},
		{
			name: "negative savings",
			summary: FinancialSummary{
				TotalIncome:   4000,
				TotalExpenses: 5000,
				NetSavings:    -1000,
			},
			expected: -25, // (-1000 / 4000) * 100
		},
		{
			name: "zero income",
			summary: FinancialSummary{
				TotalIncome:   0,
				TotalExpenses: 1000,
				NetSavings:    -1000,
			},
			expected: 0, // Should handle division by zero
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.summary.CalculateSavingsRate()
			if tt.summary.SavingsRate != tt.expected {
				t.Errorf("CalculateSavingsRate() = %v, want %v", tt.summary.SavingsRate, tt.expected)
			}
		})
	}
}

