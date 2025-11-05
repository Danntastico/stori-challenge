package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danntastico/stori-backend/internal/domain"
	"github.com/danntastico/stori-backend/internal/repository"
	"github.com/danntastico/stori-backend/internal/service"
)

// Test data
var testJSON = []byte(`[
	{"date": "2024-01-01", "amount": 2800, "category": "salary", "description": "Bi-weekly salary", "type": "income"},
	{"date": "2024-01-02", "amount": -1200, "category": "rent", "description": "Monthly rent", "type": "expense"},
	{"date": "2024-01-03", "amount": -85, "category": "groceries", "description": "Whole Foods", "type": "expense"},
	{"date": "2024-02-01", "amount": 2800, "category": "salary", "description": "Bi-weekly salary", "type": "income"}
]`)

// MockAnalyticsService implements AnalyticsServiceInterface for testing
// This allows us to test handlers in isolation without real services
type MockAnalyticsService struct {
	GetCategorySummaryFunc        func() (*domain.CategorySummary, error)
	GetTimelineFunc               func() (*domain.TimelineResponse, error)
	GetTransactionsFunc          func() (*domain.TransactionsResponse, error)
	GetTransactionsByDateRangeFunc func(start, end time.Time) (*domain.TransactionsResponse, error)
}

// Ensure MockAnalyticsService implements the interface (compile-time check)
var _ service.AnalyticsServiceInterface = (*MockAnalyticsService)(nil)

func (m *MockAnalyticsService) GetCategorySummary() (*domain.CategorySummary, error) {
	if m.GetCategorySummaryFunc != nil {
		return m.GetCategorySummaryFunc()
	}
	return nil, errors.New("GetCategorySummary not implemented in mock")
}

func (m *MockAnalyticsService) GetTimeline() (*domain.TimelineResponse, error) {
	if m.GetTimelineFunc != nil {
		return m.GetTimelineFunc()
	}
	return nil, errors.New("GetTimeline not implemented in mock")
}

func (m *MockAnalyticsService) GetTransactions() (*domain.TransactionsResponse, error) {
	if m.GetTransactionsFunc != nil {
		return m.GetTransactionsFunc()
	}
	return nil, errors.New("GetTransactions not implemented in mock")
}

func (m *MockAnalyticsService) GetTransactionsByDateRange(start, end time.Time) (*domain.TransactionsResponse, error) {
	if m.GetTransactionsByDateRangeFunc != nil {
		return m.GetTransactionsByDateRangeFunc(start, end)
	}
	return nil, errors.New("GetTransactionsByDateRange not implemented in mock")
}

// MockAIService implements AIServiceInterface for testing
// This allows us to test AI advice handler without calling OpenAI API
type MockAIService struct {
	GetFinancialAdviceFunc func(ctx context.Context, summary domain.CategorySummary, req service.AdviceRequest) (*service.AdviceResponse, error)
}

// Ensure MockAIService implements the interface (compile-time check)
var _ service.AIServiceInterface = (*MockAIService)(nil)

func (m *MockAIService) GetFinancialAdvice(ctx context.Context, summary domain.CategorySummary, req service.AdviceRequest) (*service.AdviceResponse, error) {
	if m.GetFinancialAdviceFunc != nil {
		return m.GetFinancialAdviceFunc(ctx, summary, req)
	}
	return nil, errors.New("GetFinancialAdvice not implemented in mock")
}

func setupTestHandlers(t *testing.T) (*TransactionHandler, *SummaryHandler) {
	t.Helper()

	repo, err := repository.NewJSONRepository(testJSON)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	analyticsService := service.NewAnalyticsService(repo)
	transactionHandler := NewTransactionHandler(analyticsService)
	summaryHandler := NewSummaryHandler(analyticsService)

	return transactionHandler, summaryHandler
}

func TestHealthHandler(t *testing.T) {
	handler := NewHealthHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response domain.HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
	}

	if response.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
}

func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	handler := NewHealthHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestTransactionHandler_GetAll(t *testing.T) {
	handler, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var response domain.TransactionsResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedCount := 4
	if response.Count != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, response.Count)
	}

	if len(response.Transactions) != expectedCount {
		t.Errorf("Expected %d transactions, got %d", expectedCount, len(response.Transactions))
	}
}

func TestTransactionHandler_GetByDateRange(t *testing.T) {
	handler, _ := setupTestHandlers(t)

	tests := []struct {
		name           string
		startDate      string
		endDate        string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "valid date range - january only",
			startDate:      "2024-01-01",
			endDate:        "2024-01-31",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:           "valid date range - all data",
			startDate:      "2024-01-01",
			endDate:        "2024-12-31",
			expectedStatus: http.StatusOK,
			expectedCount:  4,
		},
		{
			name:           "invalid start date format",
			startDate:      "01-01-2024",
			endDate:        "2024-12-31",
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:           "invalid end date format",
			startDate:      "2024-01-01",
			endDate:        "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/transactions?startDate=" + tt.startDate + "&endDate=" + tt.endDate
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response domain.TransactionsResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response.Count != tt.expectedCount {
					t.Errorf("Expected count %d, got %d", tt.expectedCount, response.Count)
				}
			}
		})
	}
}

func TestTransactionHandler_MethodNotAllowed(t *testing.T) {
	handler, _ := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestSummaryHandler_GetCategorySummary(t *testing.T) {
	_, handler := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/api/summary/categories", nil)
	w := httptest.NewRecorder()

	handler.HandleCategorySummary(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var response domain.CategorySummary
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify structure
	if response.Income == nil {
		t.Error("Expected income map to be non-nil")
	}

	if response.Expenses == nil {
		t.Error("Expected expenses map to be non-nil")
	}

	// Verify we have income categories
	if len(response.Income) == 0 {
		t.Error("Expected at least one income category")
	}

	// Verify we have expense categories
	if len(response.Expenses) == 0 {
		t.Error("Expected at least one expense category")
	}

	// Verify financial summary
	if response.Summary.TotalIncome <= 0 {
		t.Error("Expected positive total income")
	}

	if response.Summary.TotalExpenses <= 0 {
		t.Error("Expected positive total expenses")
	}

	// Verify period information
	if response.Period.Start == "" {
		t.Error("Expected period start date")
	}

	if response.Period.End == "" {
		t.Error("Expected period end date")
	}

	if response.Period.Months <= 0 {
		t.Error("Expected positive number of months")
	}
}

func TestSummaryHandler_GetTimeline(t *testing.T) {
	_, handler := setupTestHandlers(t)

	req := httptest.NewRequest(http.MethodGet, "/api/summary/timeline", nil)
	w := httptest.NewRecorder()

	handler.HandleTimeline(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var response domain.TimelineResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify aggregation type
	if response.Aggregation != "monthly" {
		t.Errorf("Expected aggregation 'monthly', got '%s'", response.Aggregation)
	}

	// Verify we have timeline data
	if len(response.Timeline) == 0 {
		t.Error("Expected at least one timeline point")
	}

	// Verify timeline points have required fields
	for i, point := range response.Timeline {
		if point.Period == "" {
			t.Errorf("Timeline point %d has empty period", i)
		}

		// Income should be non-negative
		if point.Income < 0 {
			t.Errorf("Timeline point %d has negative income: %v", i, point.Income)
		}

		// Expenses should be non-negative (we convert to positive)
		if point.Expenses < 0 {
			t.Errorf("Timeline point %d has negative expenses: %v", i, point.Expenses)
		}
	}

	// Verify timeline is sorted
	for i := 1; i < len(response.Timeline); i++ {
		if response.Timeline[i-1].Period > response.Timeline[i].Period {
			t.Error("Timeline is not sorted chronologically")
			break
		}
	}
}

func TestSummaryHandler_MethodNotAllowed(t *testing.T) {
	_, handler := setupTestHandlers(t)

	tests := []struct {
		name    string
		path    string
		handler http.HandlerFunc
	}{
		{"categories POST", "/api/summary/categories", handler.HandleCategorySummary},
		{"timeline POST", "/api/summary/timeline", handler.HandleTimeline},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status 405, got %d", w.Code)
			}
		})
	}
}

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithError(w, http.StatusBadRequest, "Test error message")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if response.Error != "Bad Request" {
		t.Errorf("Expected error 'Bad Request', got '%s'", response.Error)
	}

	if response.Message != "Test error message" {
		t.Errorf("Expected message 'Test error message', got '%s'", response.Message)
	}
}

func TestHandleServiceError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "ErrNoTransactions",
			err:            domain.ErrNoTransactions,
			expectedStatus: http.StatusOK, // We return 200 for "no results"
		},
		{
			name:           "ErrInvalidDateRange",
			err:            domain.ErrInvalidDateRange,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrInvalidDate",
			err:            domain.ErrInvalidDate,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unknown error",
			err:            errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			handleServiceError(w, tt.err)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// ============================================================================
// Tests demonstrating the power of interfaces and mocking
// ============================================================================

// TestAdviceHandler_ServiceError tests that the handler correctly handles
// errors from the analytics service. This was difficult to test before
// without interfaces because we couldn't easily simulate service errors.
func TestAdviceHandler_ServiceError(t *testing.T) {
	// Create mock analytics service that returns an error
	mockAnalytics := &MockAnalyticsService{
		GetCategorySummaryFunc: func() (*domain.CategorySummary, error) {
			return nil, errors.New("database connection failed")
		},
	}

	// Create mock AI service (not needed for this test, but required by handler)
	mockAI := &MockAIService{}

	// Create handler with mocks
	handler := NewAdviceHandler(mockAnalytics, mockAI)

	// Create request with valid JSON body
	reqBody := `{"context": "general"}`
	req := httptest.NewRequest(http.MethodPost, "/api/advice", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	handler.GetAdvice(w, req)

	// Verify error response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if response.Message != "Failed to analyze financial data" {
		t.Errorf("Expected message 'Failed to analyze financial data', got '%s'", response.Message)
	}
}

// TestAdviceHandler_AIError tests that the handler correctly handles
// errors from the AI service (e.g., API timeout, network error).
func TestAdviceHandler_AIError(t *testing.T) {
	// Create mock analytics service that returns valid data
	mockAnalytics := &MockAnalyticsService{
		GetCategorySummaryFunc: func() (*domain.CategorySummary, error) {
			return &domain.CategorySummary{
				Income:   make(map[string]domain.CategoryDetail),
				Expenses: make(map[string]domain.CategoryDetail),
				Summary: domain.FinancialSummary{
					TotalIncome:   1000,
					TotalExpenses: 500,
					NetSavings:    500,
				},
				Period: domain.Period{
					Start:  "2024-01-01",
					End:    "2024-01-31",
					Months: 1,
				},
			}, nil
		},
	}

	// Create mock AI service that returns an error (simulating API failure)
	mockAI := &MockAIService{
		GetFinancialAdviceFunc: func(ctx context.Context, summary domain.CategorySummary, req service.AdviceRequest) (*service.AdviceResponse, error) {
			return nil, errors.New("OpenAI API timeout")
		},
	}

	// Create handler with mocks
	handler := NewAdviceHandler(mockAnalytics, mockAI)

	// Create request with valid JSON body
	reqBody := `{"context": "general"}`
	req := httptest.NewRequest(http.MethodPost, "/api/advice", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	handler.GetAdvice(w, req)

	// Verify error response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if response.Message != "Failed to generate advice" {
		t.Errorf("Expected message 'Failed to generate advice', got '%s'", response.Message)
	}
}

// TestAdviceHandler_Success tests a successful advice request with mocks.
// This demonstrates how easy it is to test handlers with controlled data.
func TestAdviceHandler_Success(t *testing.T) {
	// Create mock analytics service
	mockAnalytics := &MockAnalyticsService{
		GetCategorySummaryFunc: func() (*domain.CategorySummary, error) {
			return &domain.CategorySummary{
				Income:   make(map[string]domain.CategoryDetail),
				Expenses: make(map[string]domain.CategoryDetail),
				Summary: domain.FinancialSummary{
					TotalIncome:   1000,
					TotalExpenses: 500,
					NetSavings:    500,
				},
				Period: domain.Period{
					Start:  "2024-01-01",
					End:    "2024-01-31",
					Months: 1,
				},
			}, nil
		},
	}

	// Create mock AI service that returns advice
	mockAI := &MockAIService{
		GetFinancialAdviceFunc: func(ctx context.Context, summary domain.CategorySummary, req service.AdviceRequest) (*service.AdviceResponse, error) {
			return &service.AdviceResponse{
				Advice:          "You're doing well!",
				Insights:        []string{"You save 50% of your income"},
				Recommendations: []string{"Keep up the good work"},
				Timestamp:       time.Now().Format(time.RFC3339),
			}, nil
		},
	}

	// Create handler with mocks
	handler := NewAdviceHandler(mockAnalytics, mockAI)

	// Create request with valid JSON body
	reqBody := `{"context": "general"}`
	req := httptest.NewRequest(http.MethodPost, "/api/advice", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute handler
	handler.GetAdvice(w, req)

	// Verify success response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response service.AdviceResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Advice == "" {
		t.Error("Expected non-empty advice")
	}

	if len(response.Insights) == 0 {
		t.Error("Expected at least one insight")
	}
}

// TestSummaryHandler_ServiceError demonstrates testing error scenarios
// with the summary handler using mocks.
func TestSummaryHandler_ServiceError(t *testing.T) {
	// Create mock service that returns an error
	mockAnalytics := &MockAnalyticsService{
		GetCategorySummaryFunc: func() (*domain.CategorySummary, error) {
			return nil, errors.New("service unavailable")
		},
	}

	// Create handler with mock
	handler := NewSummaryHandler(mockAnalytics)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/summary/categories", nil)
	w := httptest.NewRecorder()

	// Execute handler
	handler.HandleCategorySummary(w, req)

	// Verify error response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// Benefits of using interfaces demonstrated above:
// 1. ✅ Easy to test error scenarios (no need to break real services)
// 2. ✅ Fast tests (no API calls, no database, no file I/O)
// 3. ✅ Isolated tests (test handler logic independently)
// 4. ✅ Predictable tests (controlled data, no randomness)
// 5. ✅ No external dependencies (tests run offline)

