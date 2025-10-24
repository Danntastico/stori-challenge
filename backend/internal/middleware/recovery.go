package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery middleware recovers from panics and logs the error
// Prevents the server from crashing on unexpected errors
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				// Return 500 Internal Server Error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

