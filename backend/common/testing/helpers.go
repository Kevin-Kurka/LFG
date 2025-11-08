package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// TestDB creates a test database connection
func TestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dbURL := "postgres://lfg_test:test_password@localhost:5432/lfg_test?sslmode=disable"
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return db
}

// MakeRequest creates an HTTP request for testing
func MakeRequest(method, path string, body interface{}) *http.Request {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// MakeAuthRequest creates an authenticated HTTP request for testing
func MakeAuthRequest(method, path string, body interface{}, token string) *http.Request {
	req := MakeRequest(method, path, body)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// ParseJSONResponse parses JSON response body into target struct
func ParseJSONResponse(t *testing.T, resp *httptest.ResponseRecorder, target interface{}) {
	t.Helper()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}
}

// AssertStatus checks HTTP status code
func AssertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Status code = %d, want %d", got, want)
	}
}

// AssertJSON checks if response body contains expected JSON
func AssertJSON(t *testing.T, got, want string) {
	t.Helper()

	var gotJSON, wantJSON interface{}
	if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
		t.Fatalf("Got invalid JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(want), &wantJSON); err != nil {
		t.Fatalf("Want invalid JSON: %v", err)
	}

	gotBytes, _ := json.Marshal(gotJSON)
	wantBytes, _ := json.Marshal(wantJSON)

	if string(gotBytes) != string(wantBytes) {
		t.Errorf("JSON mismatch:\nGot:  %s\nWant: %s", gotBytes, wantBytes)
	}
}

// CleanupDB cleans up test database tables
func CleanupDB(t *testing.T, db *sqlx.DB) {
	t.Helper()

	tables := []string{"orders", "markets", "transactions", "wallets", "users"}
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " CASCADE")
		if err != nil {
			t.Logf("Warning: Failed to cleanup table %s: %v", table, err)
		}
	}
}
