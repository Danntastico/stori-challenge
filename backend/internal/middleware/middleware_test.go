package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	allowedOrigins := []string{"http://localhost:5173", "http://localhost:3000"}
	middleware := CORS(allowedOrigins)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	tests := []struct {
		name           string
		origin         string
		method         string
		expectOrigin   string
		expectStatus   int
		expectMethods  string
		expectHeaders  string
	}{
		{
			name:          "allowed origin - localhost:5173",
			origin:        "http://localhost:5173",
			method:        "GET",
			expectOrigin:  "http://localhost:5173",
			expectStatus:  http.StatusOK,
			expectMethods: "GET, POST, PUT, DELETE, OPTIONS",
			expectHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization",
		},
		{
			name:          "allowed origin - localhost:3000",
			origin:        "http://localhost:3000",
			method:        "GET",
			expectOrigin:  "http://localhost:3000",
			expectStatus:  http.StatusOK,
			expectMethods: "GET, POST, PUT, DELETE, OPTIONS",
			expectHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization",
		},
		{
			name:          "disallowed origin",
			origin:        "http://evil-site.com",
			method:        "GET",
			expectOrigin:  "",
			expectStatus:  http.StatusOK,
			expectMethods: "GET, POST, PUT, DELETE, OPTIONS",
			expectHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization",
		},
		{
			name:          "OPTIONS preflight request",
			origin:        "http://localhost:5173",
			method:        "OPTIONS",
			expectOrigin:  "http://localhost:5173",
			expectStatus:  http.StatusOK,
			expectMethods: "GET, POST, PUT, DELETE, OPTIONS",
			expectHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", tt.origin)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, w.Code)
			}

			allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if allowOrigin != tt.expectOrigin {
				t.Errorf("Expected origin '%s', got '%s'", tt.expectOrigin, allowOrigin)
			}

			allowMethods := w.Header().Get("Access-Control-Allow-Methods")
			if allowMethods != tt.expectMethods {
				t.Errorf("Expected methods '%s', got '%s'", tt.expectMethods, allowMethods)
			}

			allowHeaders := w.Header().Get("Access-Control-Allow-Headers")
			if allowHeaders != tt.expectHeaders {
				t.Errorf("Expected headers '%s', got '%s'", tt.expectHeaders, allowHeaders)
			}

			// Check Max-Age header
			maxAge := w.Header().Get("Access-Control-Max-Age")
			if maxAge != "86400" {
				t.Errorf("Expected Max-Age '86400', got '%s'", maxAge)
			}

			// Check Credentials header
			credentials := w.Header().Get("Access-Control-Allow-Credentials")
			if credentials != "true" {
				t.Errorf("Expected Credentials 'true', got '%s'", credentials)
			}
		})
	}
}

func TestCORS_Wildcard(t *testing.T) {
	allowedOrigins := []string{"*"}
	middleware := CORS(allowedOrigins)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://any-origin.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "http://any-origin.com" {
		t.Errorf("Expected wildcard to allow any origin, got '%s'", allowOrigin)
	}
}

func TestCORS_EmptyAllowedOrigins(t *testing.T) {
	allowedOrigins := []string{}
	middleware := CORS(allowedOrigins)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "" {
		t.Errorf("Expected no origin to be allowed with empty list, got '%s'", allowOrigin)
	}
}

func TestIsOriginAllowed(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		allowedOrigins []string
		expected       bool
	}{
		{
			name:           "exact match",
			origin:         "http://localhost:5173",
			allowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
			expected:       true,
		},
		{
			name:           "not in list",
			origin:         "http://evil.com",
			allowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
			expected:       false,
		},
		{
			name:           "wildcard",
			origin:         "http://any-origin.com",
			allowedOrigins: []string{"*"},
			expected:       true,
		},
		{
			name:           "empty list",
			origin:         "http://localhost:5173",
			allowedOrigins: []string{},
			expected:       false,
		},
		{
			name:           "trailing slash handling",
			origin:         "http://localhost:5173",
			allowedOrigins: []string{"http://localhost:5173/"},
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOriginAllowed(tt.origin, tt.allowedOrigins)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLogger(t *testing.T) {
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", body)
	}
}

func TestLogger_CapturesStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		handlerStatus  int
		expectedStatus int
	}{
		{"status 200", http.StatusOK, http.StatusOK},
		{"status 404", http.StatusNotFound, http.StatusNotFound},
		{"status 500", http.StatusInternalServerError, http.StatusInternalServerError},
		{"status 201", http.StatusCreated, http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.handlerStatus)
			}))

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestLogger_DefaultStatusCode(t *testing.T) {
	// When handler doesn't explicitly set status, should default to 200
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected default status 200, got %d", w.Code)
	}
}

func TestRecovery(t *testing.T) {
	handler := Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Should not panic - recovery should catch it
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 after panic, got %d", w.Code)
	}

	body := w.Body.String()
	if body != "Internal Server Error\n" {
		t.Errorf("Expected 'Internal Server Error' message, got '%s'", body)
	}
}

func TestRecovery_NoPanic(t *testing.T) {
	handler := Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", body)
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	rw.WriteHeader(http.StatusCreated)

	if rw.statusCode != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", rw.statusCode)
	}

	if !rw.written {
		t.Error("Expected written flag to be true")
	}
}

func TestResponseWriter_WriteHeaderOnce(t *testing.T) {
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	// First call should set the status
	rw.WriteHeader(http.StatusCreated)

	// Second call should be ignored
	rw.WriteHeader(http.StatusBadRequest)

	if rw.statusCode != http.StatusCreated {
		t.Errorf("Expected status code to remain 201, got %d", rw.statusCode)
	}
}

func TestResponseWriter_Write(t *testing.T) {
	w := httptest.NewRecorder()
	rw := newResponseWriter(w)

	data := []byte("test data")
	n, err := rw.Write(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}

	if !rw.written {
		t.Error("Expected written flag to be true after Write")
	}

	if rw.statusCode != http.StatusOK {
		t.Errorf("Expected default status 200, got %d", rw.statusCode)
	}
}

