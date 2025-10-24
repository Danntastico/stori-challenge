package domain

import (
	"math"
	"time"
)

// Transaction represents a single financial transaction
type Transaction struct {
	Date        string  `json:"date"`        // ISO 8601 format (YYYY-MM-DD)
	Amount      float64 `json:"amount"`      // Positive for income, negative for expenses
	Category    string  `json:"category"`    // e.g., "salary", "rent", "groceries"
	Description string  `json:"description"` // Human-readable description
	Type        string  `json:"type"`        // "income" or "expense"
}

// Period represents a time range
type Period struct {
	Start  string `json:"start"`  // ISO 8601 format
	End    string `json:"end"`    // ISO 8601 format
	Months int    `json:"months"` // Number of months in period
}

// CategoryDetail holds aggregated data for a single category
type CategoryDetail struct {
	Total      float64 `json:"total"`      // Total amount for this category
	Count      int     `json:"count"`      // Number of transactions
	Percentage float64 `json:"percentage"` // Percentage of total expenses/income
}

// FinancialSummary provides high-level financial metrics
type FinancialSummary struct {
	TotalIncome   float64 `json:"total_income"`   // Sum of all income
	TotalExpenses float64 `json:"total_expenses"` // Sum of all expenses (positive value)
	NetSavings    float64 `json:"net_savings"`    // Income - Expenses
	SavingsRate   float64 `json:"savings_rate"`   // (NetSavings / TotalIncome) * 100
}

// CategorySummary contains category-wise breakdown and overall summary
type CategorySummary struct {
	Income   map[string]CategoryDetail `json:"income"`   // Income categories
	Expenses map[string]CategoryDetail `json:"expenses"` // Expense categories
	Summary  FinancialSummary          `json:"summary"`  // Overall financial summary
	Period   Period                    `json:"period"`   // Time period covered
}

// TimelinePoint represents aggregated data for a specific time period
type TimelinePoint struct {
	Period   string  `json:"period"`   // "YYYY-MM" for monthly
	Income   float64 `json:"income"`   // Total income for period
	Expenses float64 `json:"expenses"` // Total expenses for period (positive value)
	Net      float64 `json:"net"`      // Income - Expenses
}

// TimelineResponse contains the timeline data
type TimelineResponse struct {
	Timeline    []TimelinePoint `json:"timeline"`    // Ordered time series data
	Aggregation string          `json:"aggregation"` // "monthly" or "weekly"
}

// TransactionsResponse contains transactions with metadata
type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"` // List of transactions
	Count        int           `json:"count"`        // Total count
	Period       Period        `json:"period"`       // Time period covered
}

// AIAdviceRequest represents a request for financial advice
type AIAdviceRequest struct {
	Context  string `json:"context"`  // "general", "savings", "budgeting", "specific_category"
	Category string `json:"category"` // Optional: specific category for targeted advice
}

// AIAdviceResponse contains AI-generated financial advice
type AIAdviceResponse struct {
	Advice          string    `json:"advice"`          // Main advice text
	Insights        []string  `json:"insights"`        // Key insights discovered
	Recommendations []string  `json:"recommendations"` // Actionable recommendations
	Timestamp       time.Time `json:"timestamp"`       // When advice was generated
}

// HealthResponse represents API health status
type HealthResponse struct {
	Status    string    `json:"status"`    // "healthy" or "unhealthy"
	Timestamp time.Time `json:"timestamp"` // Current server time
}

// Helper methods

// IsIncome returns true if the transaction is income
func (t *Transaction) IsIncome() bool {
	return t.Type == "income"
}

// IsExpense returns true if the transaction is an expense
func (t *Transaction) IsExpense() bool {
	return t.Type == "expense"
}

// AbsoluteAmount returns the absolute value of the amount
func (t *Transaction) AbsoluteAmount() float64 {
	return math.Abs(t.Amount)
}

// ParseDate parses the transaction date into a time.Time
func (t *Transaction) ParseDate() (time.Time, error) {
	return time.Parse("2006-01-02", t.Date)
}

// GetYearMonth returns the year-month string (YYYY-MM) for timeline aggregation
func (t *Transaction) GetYearMonth() (string, error) {
	date, err := t.ParseDate()
	if err != nil {
		return "", err
	}
	return date.Format("2006-01"), nil
}

// Validate checks if the transaction has valid data
func (t *Transaction) Validate() error {
	if t.Date == "" {
		return ErrInvalidDate
	}
	if _, err := t.ParseDate(); err != nil {
		return ErrInvalidDate
	}
	if t.Category == "" {
		return ErrInvalidCategory
	}
	if t.Type != "income" && t.Type != "expense" {
		return ErrInvalidType
	}
	// Validate amount sign matches type
	if t.Type == "income" && t.Amount < 0 {
		return ErrInvalidAmount
	}
	if t.Type == "expense" && t.Amount > 0 {
		return ErrInvalidAmount
	}
	return nil
}

// CalculateSavingsRate computes the savings rate percentage
func (fs *FinancialSummary) CalculateSavingsRate() {
	if fs.TotalIncome > 0 {
		fs.SavingsRate = roundToTwoDecimals((fs.NetSavings / fs.TotalIncome) * 100)
	} else {
		fs.SavingsRate = 0
	}
}

// Helper function to round to 2 decimal places
func roundToTwoDecimals(val float64) float64 {
	return math.Round(val*100) / 100
}

