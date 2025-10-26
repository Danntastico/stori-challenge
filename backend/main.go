package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/danntastico/stori-backend/internal/handlers"
	"github.com/danntastico/stori-backend/internal/middleware"
	"github.com/danntastico/stori-backend/internal/repository"
	"github.com/danntastico/stori-backend/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

//go:embed data/transactions.json
var transactionsData []byte

func main() {
	// Load environment variables
	config := loadConfig()

	log.Println("🚀 Starting Stori Financial Tracker API...")
	log.Printf("📊 Loaded %d bytes of transaction data", len(transactionsData))

	// Initialize repository
	repo, err := repository.NewJSONRepository(transactionsData)
	if err != nil {
		log.Fatalf("❌ Failed to initialize repository: %v", err)
	}
	log.Printf("✅ Repository initialized with %d transactions", repo.Count())

	// Initialize analytics service
	analyticsService := service.NewAnalyticsService(repo)
	log.Println("✅ Analytics service initialized")

	// Initialize AI service
	aiService := service.NewAIService(config.OpenAIAPIKey)
	if config.OpenAIAPIKey == "" {
		log.Println("⚠️  OpenAI API key not provided - using mock responses")
	} else {
		log.Println("✅ AI service initialized with OpenAI integration")
	}

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	transactionHandler := handlers.NewTransactionHandler(analyticsService)
	summaryHandler := handlers.NewSummaryHandler(analyticsService)
	adviceHandler := handlers.NewAdviceHandler(analyticsService, aiService)
	log.Println("✅ Handlers initialized")

	// Initialize chi router
	r := chi.NewRouter()

	// Register middleware (order matters!)
	r.Use(middleware.Recovery)                    // 1. Catch panics
	r.Use(middleware.Logger)                      // 2. Log requests
	r.Use(chimiddleware.RequestID)                // 3. Add request ID
	r.Use(chimiddleware.RealIP)                   // 4. Get real IP
	r.Use(middleware.CORS(config.AllowedOrigins)) // 5. Handle CORS
	r.Use(chimiddleware.Timeout(60 * time.Second)) // 6. Request timeout

	log.Println("✅ Middleware registered")

	// Register routes
	r.Get("/api/health", healthHandler.ServeHTTP)
	r.Get("/api/transactions", transactionHandler.ServeHTTP)
	r.Get("/api/summary/categories", summaryHandler.HandleCategorySummary)
	r.Get("/api/summary/timeline", summaryHandler.HandleTimeline)
	r.Post("/api/advice", adviceHandler.GetAdvice)

	// Root endpoint for API info
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"name": "Stori Financial Tracker API",
			"version": "1.0.0",
			"status": "running",
			"endpoints": {
				"health": "/api/health",
				"transactions": "/api/transactions",
				"categories": "/api/summary/categories",
				"timeline": "/api/summary/timeline",
				"advice": "/api/advice"
			}
		}`))
	})

	log.Println("✅ Routes registered")

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🌐 Server listening on http://localhost:%s", config.Port)
		log.Println("📡 API endpoints:")
		log.Println("   GET  /api/health")
		log.Println("   GET  /api/transactions")
		log.Println("   GET  /api/summary/categories")
		log.Println("   GET  /api/summary/timeline")
		log.Println("   POST /api/advice")
		log.Println("💡 Press Ctrl+C to shutdown")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n🛑 Shutdown signal received, gracefully shutting down...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}

// Config holds application configuration
type Config struct {
	Port           string
	AllowedOrigins []string
	LogLevel       string
	OpenAIAPIKey   string
}

// loadConfig loads configuration from environment variables with defaults
func loadConfig() Config {
	port := getEnv("PORT", "8080")
	originsStr := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000")
	logLevel := getEnv("LOG_LEVEL", "info")
	openAIAPIKey := getEnv("OPENAI_API_KEY", "")

	// Parse allowed origins
	var allowedOrigins []string
	if originsStr != "" {
		origins := strings.Split(originsStr, ",")
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}

	config := Config{
		Port:           port,
		AllowedOrigins: allowedOrigins,
		LogLevel:       logLevel,
		OpenAIAPIKey:   openAIAPIKey,
	}

	log.Println("⚙️  Configuration loaded:")
	log.Printf("   Port: %s", config.Port)
	log.Printf("   Allowed Origins: %v", config.AllowedOrigins)
	log.Printf("   Log Level: %s", config.LogLevel)

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

