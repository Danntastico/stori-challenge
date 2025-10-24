package service

import (
	"math"
	"sort"
	"time"

	"github.com/danntastico/stori-backend/internal/domain"
	"github.com/danntastico/stori-backend/internal/repository"
)

// AnalyticsService provides business logic for financial data analysis
type AnalyticsService struct {
	repo repository.TransactionRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(repo repository.TransactionRepository) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
	}
}

// GetCategorySummary calculates spending breakdown by category with totals and percentages
func (s *AnalyticsService) GetCategorySummary() (*domain.CategorySummary, error) {
	// Fetch all transactions
	transactions, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Initialize maps for income and expense categories
	incomeCategories := make(map[string]*domain.CategoryDetail)
	expenseCategories := make(map[string]*domain.CategoryDetail)

	var totalIncome float64
	var totalExpenses float64

	// Aggregate transactions by category
	for _, tx := range transactions {
		if tx.IsIncome() {
			totalIncome += tx.Amount
			s.aggregateCategory(incomeCategories, tx)
		} else if tx.IsExpense() {
			totalExpenses += tx.AbsoluteAmount()
			s.aggregateCategory(expenseCategories, tx)
		}
	}

	// Calculate percentages for income categories
	incomeMap := s.calculatePercentages(incomeCategories, totalIncome)

	// Calculate percentages for expense categories
	expenseMap := s.calculatePercentages(expenseCategories, totalExpenses)

	// Get date range
	start, end, err := s.getDateRangeFromTransactions(transactions)
	if err != nil {
		return nil, err
	}

	// Calculate number of months
	months := s.calculateMonthsBetween(start, end)

	// Create financial summary
	summary := domain.FinancialSummary{
		TotalIncome:   roundToTwo(totalIncome),
		TotalExpenses: roundToTwo(totalExpenses),
		NetSavings:    roundToTwo(totalIncome - totalExpenses),
	}
	summary.CalculateSavingsRate()

	return &domain.CategorySummary{
		Income:   incomeMap,
		Expenses: expenseMap,
		Summary:  summary,
		Period: domain.Period{
			Start:  start.Format("2006-01-02"),
			End:    end.Format("2006-01-02"),
			Months: months,
		},
	}, nil
}

// GetTimeline calculates monthly income vs expenses over time
func (s *AnalyticsService) GetTimeline() (*domain.TimelineResponse, error) {
	// Fetch all transactions
	transactions, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Group transactions by month
	monthlyData := make(map[string]*domain.TimelinePoint)

	for _, tx := range transactions {
		yearMonth, err := tx.GetYearMonth()
		if err != nil {
			// Skip transactions with invalid dates
			continue
		}

		// Initialize month if not exists
		if _, exists := monthlyData[yearMonth]; !exists {
			monthlyData[yearMonth] = &domain.TimelinePoint{
				Period:   yearMonth,
				Income:   0,
				Expenses: 0,
				Net:      0,
			}
		}

		// Aggregate by type
		if tx.IsIncome() {
			monthlyData[yearMonth].Income += tx.Amount
		} else if tx.IsExpense() {
			monthlyData[yearMonth].Expenses += tx.AbsoluteAmount()
		}
	}

	// Calculate net for each month and round values
	for _, point := range monthlyData {
		point.Income = roundToTwo(point.Income)
		point.Expenses = roundToTwo(point.Expenses)
		point.Net = roundToTwo(point.Income - point.Expenses)
	}

	// Convert map to sorted slice
	timeline := make([]domain.TimelinePoint, 0, len(monthlyData))
	for _, point := range monthlyData {
		timeline = append(timeline, *point)
	}

	// Sort by period (chronologically)
	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Period < timeline[j].Period
	})

	return &domain.TimelineResponse{
		Timeline:    timeline,
		Aggregation: "monthly",
	}, nil
}

// GetTransactions returns all transactions with metadata
func (s *AnalyticsService) GetTransactions() (*domain.TransactionsResponse, error) {
	transactions, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	start, end, err := s.getDateRangeFromTransactions(transactions)
	if err != nil {
		return nil, err
	}

	return &domain.TransactionsResponse{
		Transactions: transactions,
		Count:        len(transactions),
		Period: domain.Period{
			Start: start.Format("2006-01-02"),
			End:   end.Format("2006-01-02"),
		},
	}, nil
}

// GetTransactionsByDateRange returns filtered transactions within a date range
func (s *AnalyticsService) GetTransactionsByDateRange(start, end time.Time) (*domain.TransactionsResponse, error) {
	transactions, err := s.repo.GetByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	return &domain.TransactionsResponse{
		Transactions: transactions,
		Count:        len(transactions),
		Period: domain.Period{
			Start: start.Format("2006-01-02"),
			End:   end.Format("2006-01-02"),
		},
	}, nil
}

// Helper methods

// aggregateCategory adds a transaction to the category aggregation
func (s *AnalyticsService) aggregateCategory(categories map[string]*domain.CategoryDetail, tx domain.Transaction) {
	if _, exists := categories[tx.Category]; !exists {
		categories[tx.Category] = &domain.CategoryDetail{
			Total:      0,
			Count:      0,
			Percentage: 0,
		}
	}

	categories[tx.Category].Total += tx.AbsoluteAmount()
	categories[tx.Category].Count++
}

// calculatePercentages converts category map to final format with percentages
func (s *AnalyticsService) calculatePercentages(categories map[string]*domain.CategoryDetail, total float64) map[string]domain.CategoryDetail {
	result := make(map[string]domain.CategoryDetail)

	for category, detail := range categories {
		percentage := 0.0
		if total > 0 {
			percentage = (detail.Total / total) * 100
		}

		result[category] = domain.CategoryDetail{
			Total:      roundToTwo(detail.Total),
			Count:      detail.Count,
			Percentage: roundToTwo(percentage),
		}
	}

	return result
}

// getDateRangeFromTransactions finds the min and max dates from a slice of transactions
func (s *AnalyticsService) getDateRangeFromTransactions(transactions []domain.Transaction) (time.Time, time.Time, error) {
	if len(transactions) == 0 {
		return time.Time{}, time.Time{}, domain.ErrNoTransactions
	}

	var minDate, maxDate time.Time
	first := true

	for _, tx := range transactions {
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

// calculateMonthsBetween calculates the number of months between two dates
func (s *AnalyticsService) calculateMonthsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())

	// Add 1 because we want inclusive count (e.g., Jan to Feb is 2 months)
	return years*12 + months + 1
}

// roundToTwo rounds a float64 to 2 decimal places
func roundToTwo(val float64) float64 {
	return math.Round(val*100) / 100
}

