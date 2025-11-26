package safecomms

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *Client with Transport replaced to avoid network calls
func NewTestClient(fn RoundTripFunc) *Client {
	c := NewClient("test-api-key", "")
	c.httpClient.Transport = RoundTripFunc(fn)
	return c
}

func TestNewClient(t *testing.T) {
	c := NewClient("key", "")
	if c.apiKey != "key" {
		t.Errorf("Expected apiKey 'key', got %s", c.apiKey)
	}
	if c.baseURL != DefaultBaseURL {
		t.Errorf("Expected baseURL %s, got %s", DefaultBaseURL, c.baseURL)
	}
}

func TestModerateText(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		if req.URL.String() != DefaultBaseURL+"/moderation/text" {
			t.Errorf("Expected URL %s, got %s", DefaultBaseURL+"/moderation/text", req.URL.String())
		}
		if req.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header to be Bearer test-api-key, got %s", req.Header.Get("Authorization"))
		}

		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: io.NopCloser(bytes.NewBufferString(`{"flagged": true}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	resp, err := client.ModerateText(ModerateTextRequest{
		Content: "test content",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp["flagged"] != true {
		t.Errorf("Expected flagged true")
	}
}
