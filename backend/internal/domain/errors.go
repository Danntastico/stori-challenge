package domain

import (
	"errors"
	"fmt"
)

// Domain-level errors
var (
	// ErrInvalidDate is returned when a transaction has an invalid date format
	ErrInvalidDate = errors.New("invalid date format, expected YYYY-MM-DD")

	// ErrInvalidCategory is returned when a transaction has an empty category
	ErrInvalidCategory = errors.New("category cannot be empty")

	// ErrInvalidType is returned when a transaction type is not "income" or "expense"
	ErrInvalidType = errors.New("type must be either 'income' or 'expense'")

	// ErrInvalidAmount is returned when amount sign doesn't match transaction type
	ErrInvalidAmount = errors.New("amount sign must match transaction type")

	// ErrNoTransactions is returned when no transactions are found
	ErrNoTransactions = errors.New("no transactions found")

	// ErrInvalidDateRange is returned when date range is invalid
	ErrInvalidDateRange = errors.New("invalid date range: start date must be before end date")
)

// HTTPError represents an error with an associated HTTP status code
// This is useful for preserving status codes from external APIs (e.g., OpenAI)
// and mapping them to appropriate HTTP responses
type HTTPError struct {
	StatusCode int
	Message    string
	Err        error
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("HTTP %d: %s (%v)", e.StatusCode, e.Message, e.Err)
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// Unwrap returns the underlying error for error wrapping support
func (e *HTTPError) Unwrap() error {
	return e.Err
}

// NewHTTPError creates a new HTTPError with status code and message
func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// NewHTTPErrorWithCause creates a new HTTPError with status code, message, and underlying error
func NewHTTPErrorWithCause(statusCode int, message string, err error) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

