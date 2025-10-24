package domain

import "errors"

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

