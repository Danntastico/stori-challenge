package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/danntastico/stori-backend/internal/domain"
)

// AIService handles AI-powered financial advice generation
type AIService struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

// NewAIService creates a new AI service instance
func NewAIService(apiKey string) *AIService {
	return &AIService{
		apiKey: apiKey,
		apiURL: "https://api.openai.com/v1/chat/completions",
		// No HTTP client timeout - rely on context cancellation from handler timeout (60s)
		// The context passed via NewRequestWithContext will control when the request is cancelled
		httpClient: &http.Client{},
	}
}

// AdviceRequest represents the request structure for advice
type AdviceRequest struct {
	Context  string `json:"context"`  // "general", "savings", "budgeting", etc.
	Category string `json:"category"` // optional, for category-specific advice
}

// AdviceResponse represents the structured advice response
type AdviceResponse struct {
	Advice          string   `json:"advice"`
	Insights        []string `json:"insights"`
	Recommendations []string `json:"recommendations"`
	Timestamp       string   `json:"timestamp"`
}

// openAIRequest represents the OpenAI API request structure
type openAIRequest struct {
	Model       string                   `json:"model"`
	Messages    []openAIMessage          `json:"messages"`
	Temperature float64                  `json:"temperature"`
	MaxTokens   int                      `json:"max_tokens"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIResponse represents the OpenAI API response structure
type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// GetFinancialAdvice generates AI-powered financial advice based on summary data
func (s *AIService) GetFinancialAdvice(ctx context.Context, summary domain.CategorySummary, req AdviceRequest) (*AdviceResponse, error) {
	// If no API key, return mock advice
	if s.apiKey == "" {
		return s.getMockAdvice(summary, req), nil
	}

	// Build the prompt
	prompt := s.buildPrompt(summary, req)

	// Call OpenAI API
	advice, err := s.callOpenAI(ctx, prompt)
	if err != nil {
		// On error, fallback to mock advice
		return s.getMockAdvice(summary, req), nil
	}

	// Parse and structure the response
	response := s.parseAdviceResponse(advice, summary)
	return response, nil
}

// buildPrompt constructs the prompt for OpenAI based on financial data
func (s *AIService) buildPrompt(summary domain.CategorySummary, req AdviceRequest) string {
	prompt := "You are a helpful and encouraging financial advisor. Analyze this user's financial data and provide personalized advice.\n\n"

	// Add income information
	prompt += fmt.Sprintf("📊 Financial Overview:\n")
	prompt += fmt.Sprintf("Period: %s to %s (%d months)\n\n", 
		summary.Period.Start, summary.Period.End, summary.Period.Months)

	prompt += fmt.Sprintf("Income:\n")
	prompt += fmt.Sprintf("- Total: $%.2f\n", summary.Summary.TotalIncome)
	prompt += fmt.Sprintf("- Average monthly: $%.2f\n\n", summary.Summary.TotalIncome/float64(summary.Period.Months))

	// Add expense breakdown
	prompt += "Expenses by Category:\n"
	for category, detail := range summary.Expenses {
		prompt += fmt.Sprintf("- %s: $%.2f (%.1f%%, %d transactions)\n",
			category, detail.Total, detail.Percentage, detail.Count)
	}

	prompt += fmt.Sprintf("\nTotal Expenses: $%.2f\n", summary.Summary.TotalExpenses)
	prompt += fmt.Sprintf("Net Savings: $%.2f\n", summary.Summary.NetSavings)
	prompt += fmt.Sprintf("Savings Rate: %.1f%%\n\n", summary.Summary.SavingsRate)

	// Add context-specific instructions
	if req.Category != "" {
		prompt += fmt.Sprintf("Focus specifically on the '%s' category.\n\n", req.Category)
	}

	prompt += `Please provide a structured response with:

1. INSIGHTS (2-3 key observations about spending patterns)
2. RECOMMENDATIONS (3-4 specific, actionable steps to improve financial health)
3. POSITIVE REINFORCEMENT (1 encouraging statement)

Format your response as:
INSIGHTS:
- [insight 1]
- [insight 2]

RECOMMENDATIONS:
- [recommendation 1]
- [recommendation 2]

POSITIVE:
[encouraging message]

Keep advice practical, specific to the data, and encouraging. Use exact dollar amounts when relevant.`

	return prompt
}

// callOpenAI makes the HTTP request to OpenAI API
func (s *AIService) callOpenAI(ctx context.Context, prompt string) (string, error) {
	reqBody := openAIRequest{
		Model:       "gpt-3.5-turbo",
		Temperature: 0.7,
		MaxTokens:   600,
		Messages: []openAIMessage{
			{
				Role:    "system",
				Content: "You are a professional financial advisor who provides clear, actionable advice.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Map OpenAI API status codes to appropriate HTTP errors
		var statusCode int
		var message string

		switch resp.StatusCode {
		case http.StatusTooManyRequests: // 429
			// Rate limit - pass through to client (they can retry)
			statusCode = http.StatusTooManyRequests
			message = "OpenAI API rate limit exceeded. Please try again later."
		case http.StatusUnauthorized: // 401
			// Invalid API key - this is our configuration issue, but expose as 500
			statusCode = http.StatusInternalServerError
			message = "AI service configuration error"
		case http.StatusServiceUnavailable: // 503
			// OpenAI is down - map to 503 for client
			statusCode = http.StatusServiceUnavailable
			message = "AI service is temporarily unavailable. Please try again later."
		default:
			// Other errors - map to 502 (Bad Gateway) since it's an external service issue
			statusCode = http.StatusBadGateway
			message = fmt.Sprintf("AI service error (status %d)", resp.StatusCode)
		}

		return "", domain.NewHTTPErrorWithCause(statusCode, message, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body)))
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// parseAdviceResponse parses the AI response into structured format
func (s *AIService) parseAdviceResponse(advice string, summary domain.CategorySummary) *AdviceResponse {
	// Simple parsing - in production, this could be more sophisticated
	insights := []string{}
	recommendations := []string{}
	
	// Extract sections from the response
	// This is a basic implementation - could use regex or more advanced parsing
	lines := splitLines(advice)
	section := ""
	
	for _, line := range lines {
		trimmed := trim(line)
		if trimmed == "" {
			continue
		}
		
		if contains(trimmed, "INSIGHTS:") {
			section = "insights"
			continue
		}
		if contains(trimmed, "RECOMMENDATIONS:") {
			section = "recommendations"
			continue
		}
		if contains(trimmed, "POSITIVE:") {
			section = "positive"
			continue
		}
		
		if startsWith(trimmed, "-") || startsWith(trimmed, "•") {
			item := trimPrefix(trimPrefix(trimmed, "-"), "•")
			item = trim(item)
			if section == "insights" {
				insights = append(insights, item)
			} else if section == "recommendations" {
				recommendations = append(recommendations, item)
			}
		}
	}
	
	// Ensure we have at least some content
	if len(insights) == 0 {
		insights = s.getDefaultInsights(summary)
	}
	if len(recommendations) == 0 {
		recommendations = s.getDefaultRecommendations(summary)
	}

	return &AdviceResponse{
		Advice:          advice,
		Insights:        insights,
		Recommendations: recommendations,
		Timestamp:       time.Now().Format(time.RFC3339),
	}
}

// getMockAdvice returns mock advice when OpenAI is not available
func (s *AIService) getMockAdvice(summary domain.CategorySummary, req AdviceRequest) *AdviceResponse {
	insights := s.getDefaultInsights(summary)
	recommendations := s.getDefaultRecommendations(summary)

	advice := "Based on your financial data analysis:\n\n"
	advice += "INSIGHTS:\n"
	for _, insight := range insights {
		advice += "- " + insight + "\n"
	}
	advice += "\nRECOMMENDATIONS:\n"
	for _, rec := range recommendations {
		advice += "- " + rec + "\n"
	}
	advice += "\nPOSITIVE:\n"
	advice += "You're tracking your finances, which is a great first step toward financial wellness!"

	return &AdviceResponse{
		Advice:          advice,
		Insights:        insights,
		Recommendations: recommendations,
		Timestamp:       time.Now().Format(time.RFC3339),
	}
}

// getDefaultInsights generates insights based on the data
func (s *AIService) getDefaultInsights(summary domain.CategorySummary) []string {
	insights := []string{}

	savingsRate := summary.Summary.SavingsRate
	if savingsRate > 20 {
		insights = append(insights, fmt.Sprintf("Excellent savings rate of %.1f%% - you're saving more than the recommended 20%%", savingsRate))
	} else if savingsRate > 10 {
		insights = append(insights, fmt.Sprintf("Your savings rate of %.1f%% is on track - aim for 20%% for optimal financial health", savingsRate))
	} else if savingsRate > 0 {
		insights = append(insights, fmt.Sprintf("Your savings rate of %.1f%% has room for improvement - consider cutting discretionary spending", savingsRate))
	} else {
		insights = append(insights, "You're currently spending more than you earn - immediate action needed to avoid debt")
	}

	// Find largest expense category
	var largestCat string
	var largestAmt float64
	for cat, detail := range summary.Expenses {
		if detail.Total > largestAmt {
			largestAmt = detail.Total
			largestCat = cat
		}
	}
	if largestCat != "" {
		insights = append(insights, fmt.Sprintf("Your largest expense is %s at $%.2f (%.1f%% of spending)", 
			largestCat, largestAmt, (largestAmt/summary.Summary.TotalExpenses)*100))
	}

	// Monthly average
	monthlyExpenses := summary.Summary.TotalExpenses / float64(summary.Period.Months)
	insights = append(insights, fmt.Sprintf("Average monthly expenses: $%.2f over %d months", 
		monthlyExpenses, summary.Period.Months))

	return insights
}

// getDefaultRecommendations generates recommendations based on the data
func (s *AIService) getDefaultRecommendations(summary domain.CategorySummary) []string {
	recommendations := []string{}

	if summary.Summary.SavingsRate < 20 {
		recommendations = append(recommendations, "Set up automatic transfers to savings account to reach a 20% savings rate")
	}

	// Check for high discretionary spending
	discretionaryTotal := 0.0
	discretionaryCategories := []string{"dining", "entertainment", "shopping", "subscriptions"}
	for cat, detail := range summary.Expenses {
		for _, discCat := range discretionaryCategories {
			if cat == discCat {
				discretionaryTotal += detail.Total
			}
		}
	}
	
	if discretionaryTotal > summary.Summary.TotalExpenses*0.2 {
		recommendations = append(recommendations, fmt.Sprintf("Consider reducing discretionary spending (dining, entertainment, shopping) - currently $%.2f", discretionaryTotal))
	}

	recommendations = append(recommendations, "Track your spending weekly to identify patterns and opportunities to save")
	recommendations = append(recommendations, "Build an emergency fund covering 3-6 months of expenses")

	return recommendations
}

// Helper functions for string manipulation
func splitLines(s string) []string {
	result := []string{}
	current := ""
	for _, char := range s {
		if char == '\n' {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trim(s string) string {
	start := 0
	end := len(s)
	
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	
	return s[start:end]
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func trimPrefix(s, prefix string) string {
	if startsWith(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

