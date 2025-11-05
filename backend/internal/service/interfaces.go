package service

import (
	"context"
	"time"

	"github.com/danntastico/stori-backend/internal/domain"
)

// AnalyticsServiceInterface defines the contract for analytics operations
// This interface allows us to:
// 1. Mock services in handler tests
// 2. Swap implementations (e.g., cached service, rate-limited service)
// 3. Test error scenarios easily
type AnalyticsServiceInterface interface {
	GetCategorySummary() (*domain.CategorySummary, error)
	GetTimeline() (*domain.TimelineResponse, error)
	GetTransactions() (*domain.TransactionsResponse, error)
	GetTransactionsByDateRange(start, end time.Time) (*domain.TransactionsResponse, error)
}

// AIServiceInterface defines the contract for AI-powered advice generation
// This interface allows us to:
// 1. Mock AI service in handler tests (no API calls in tests)
// 2. Swap implementations (e.g., different LLM providers)
// 3. Test error scenarios (API failures, timeouts)
type AIServiceInterface interface {
	GetFinancialAdvice(ctx context.Context, summary domain.CategorySummary, req AdviceRequest) (*AdviceResponse, error)
}

// Ensure concrete types implement interfaces (compile-time check)
var (
	_ AnalyticsServiceInterface = (*AnalyticsService)(nil)
	_ AIServiceInterface         = (*AIService)(nil)
)

