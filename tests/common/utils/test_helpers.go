package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// BaseURL is the base URL for API tests
var BaseURL string

func init() {
	// Use env variable if set, otherwise default to localhost
	baseURL := os.Getenv("API_TEST_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	BaseURL = baseURL
}

// APIResponse is a generic structure for API responses
type APIResponse struct {
	StatusCode int
	Body       []byte
	Error      error
}

// MakeRequest makes an HTTP request to the API
func MakeRequest(method, endpoint string, payload interface{}) APIResponse {
	var req *http.Request
	var err error

	url := fmt.Sprintf("%s%s", BaseURL, endpoint)

	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return APIResponse{
				Error: fmt.Errorf("failed to marshal JSON: %v", err),
			}
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return APIResponse{
				Error: fmt.Errorf("failed to create request: %v", err),
			}
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return APIResponse{
				Error: fmt.Errorf("failed to create request: %v", err),
			}
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return APIResponse{
			Error: fmt.Errorf("failed to make request: %v", err),
		}
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return APIResponse{
			StatusCode: resp.StatusCode,
			Error:      fmt.Errorf("failed to read response body: %v", err),
		}
	}

	return APIResponse{
		StatusCode: resp.StatusCode,
		Body:       buf.Bytes(),
		Error:      nil,
	}
}

// AssertSuccess asserts that the API call was successful
func AssertSuccess(t *testing.T, response APIResponse) {
	if response.Error != nil {
		t.Errorf("API error: %v", response.Error)
		return
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		t.Errorf("API returned status code %d: %s", response.StatusCode, string(response.Body))
	}
	assert.NoError(t, response.Error)
	assert.GreaterOrEqual(t, response.StatusCode, 200)
	assert.Less(t, response.StatusCode, 300)
}

// ParseResponse parses the response body into the provided struct
func ParseResponse(response APIResponse, v interface{}) error {
	if response.Error != nil {
		return response.Error
	}
	return json.Unmarshal(response.Body, v)
}

// AssertError asserts that the API call returned an error or a non-2xx status code
func AssertError(t *testing.T, response APIResponse, message string) {
	if response.Error == nil && response.StatusCode >= 200 && response.StatusCode < 300 {
		t.Errorf("Expected API error, but got successful response with status code %d: %s",
			response.StatusCode, string(response.Body))
	}

	if response.Error != nil {
		t.Logf("API error as expected: %v (%s)", response.Error, message)
	} else {
		t.Logf("API error status code as expected: %d (%s)", response.StatusCode, message)
	}
}
