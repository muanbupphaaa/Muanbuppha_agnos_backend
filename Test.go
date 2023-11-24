// test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCheckPasswordHandler tests the CheckPasswordHandler function
func TestCheckPasswordHandler(t *testing.T) {
	// Create a request body with JSON payload
	requestBody := map[string]string{"password": "TestPassword123"}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test request with the JSON payload
	req, err := http.NewRequest("POST", "/check-password", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Create a handler function from the CheckPasswordHandler
	handler := http.HandlerFunc(CheckPasswordHandler)

	// Serve the HTTP request to the response recorder
	handler.ServeHTTP(rr, req)

	// Check the HTTP status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expectedResponseBody := `{"steps":0}` // Assuming strongPasswordSteps returns 0 for the test password
	if rr.Body.String() != expectedResponseBody {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponseBody)
	}
}

// You can add more tests for other handlers or functions as needed
