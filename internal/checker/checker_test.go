package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckSite(t *testing.T) {

	tests := []struct {
		name           string
		handler        http.HandlerFunc
		timeout        time.Duration
		expectedOK     bool
		expectedStatus int
		expectError    bool
	}{
		{
			name: "200 OK",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			timeout:        time.Second,
			expectedOK:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "404 Not Found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			timeout:        time.Second,
			expectedOK:     false,
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "500 Internal Server Error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			timeout:        time.Second,
			expectedOK:     false,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "timeout exceeded",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Second)
				w.WriteHeader(http.StatusOK)
			},
			timeout:     300 * time.Millisecond,
			expectedOK:  false,
			expectError: true,
		},
	}

	for _, tt := range tests {

		tt := tt // важно для t.Parallel()

		t.Run(tt.name, func(t *testing.T) {

			t.Parallel()

			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := &http.Client{
				Timeout: tt.timeout,
			}

			result := CheckSite(client, server.URL)

			// Проверка OK
			if result.OK != tt.expectedOK {
				t.Fatalf("expected OK=%v got %v",
					tt.expectedOK, result.OK)
			}

			// Проверка ошибок
			if tt.expectError {
				if result.Err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}

			if result.Err != nil {
				t.Fatalf("unexpected error: %v", result.Err)
			}

			// Проверка HTTP-кода
			if result.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status %d got %d",
					tt.expectedStatus, result.StatusCode)
			}
		})
	}
}

func TestCheckSite_ServerUnavailable(t *testing.T) {

	t.Parallel()

	client := &http.Client{
		Timeout: time.Second,
	}

	// Порт почти гарантированно закрыт
	result := CheckSite(client, "http://localhost:65534")

	if result.Err == nil {
		t.Fatal("expected connection error but got nil")
	}

	if result.OK {
		t.Fatal("expected site to be unavailable")
	}
}
