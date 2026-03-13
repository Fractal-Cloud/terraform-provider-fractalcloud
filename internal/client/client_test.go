package fractalCloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()
	host := server.URL
	id := "test-id"
	secret := "test-secret"
	return NewClient(nil, &host, &id, &secret)
}

func TestNewClient(t *testing.T) {
	host := "https://api.example.com"
	id := "my-id"
	secret := "my-secret"
	logger := &ClientLogger{}

	c := NewClient(logger, &host, &id, &secret)

	if c.HostURL != host {
		t.Errorf("HostURL = %q, want %q", c.HostURL, host)
	}
	if c.Auth.ServiceAccountId != id {
		t.Errorf("ServiceAccountId = %q, want %q", c.Auth.ServiceAccountId, id)
	}
	if c.Auth.ServiceAccountSecret != secret {
		t.Errorf("ServiceAccountSecret = %q, want %q", c.Auth.ServiceAccountSecret, secret)
	}
	if c.Logger != logger {
		t.Error("Logger not set")
	}
	if c.HTTPClient == nil {
		t.Error("HTTPClient is nil")
	}
}

func TestDoRequest_SetsHeaders(t *testing.T) {
	var gotHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := testClient(t, server)
	req, _ := http.NewRequest("GET", server.URL+"/test", nil)
	_, _, err := c.doRequest(context.Background(), req, []int{200})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := gotHeaders.Get("X-ClientID"); got != "test-id" {
		t.Errorf("X-ClientID = %q, want %q", got, "test-id")
	}
	if got := gotHeaders.Get("X-ClientSecret"); got != "test-secret" {
		t.Errorf("X-ClientSecret = %q, want %q", got, "test-secret")
	}
	if got := gotHeaders.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
}

func TestDoRequest_AcceptedStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	c := testClient(t, server)
	req, _ := http.NewRequest("POST", server.URL+"/test", nil)
	code, body, err := c.doRequest(context.Background(), req, []int{200, 202})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 202 {
		t.Errorf("status code = %d, want 202", code)
	}
	if string(body) != `{"ok":true}` {
		t.Errorf("body = %q, want %q", string(body), `{"ok":true}`)
	}
}

func TestDoRequest_UnexpectedStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	c := testClient(t, server)
	req, _ := http.NewRequest("GET", server.URL+"/test", nil)
	code, _, err := c.doRequest(context.Background(), req, []int{200})
	if err == nil {
		t.Fatal("expected error for unexpected status code")
	}
	if code != 500 {
		t.Errorf("status code = %d, want 500", code)
	}
	if !strings.Contains(err.Error(), "unexpected status 500") {
		t.Errorf("error = %q, should contain 'unexpected status 500'", err.Error())
	}
}

func TestDoRequest_HTTPError(t *testing.T) {
	// Use a server that's already closed to force a connection error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	c := testClient(t, server)
	req, _ := http.NewRequest("GET", server.URL+"/test", nil)
	_, _, err := c.doRequest(context.Background(), req, []int{200})
	if err == nil {
		t.Fatal("expected error for connection failure")
	}
	if !strings.Contains(err.Error(), "HTTP request failed") {
		t.Errorf("error = %q, should contain 'HTTP request failed'", err.Error())
	}
}

func TestTruncateBody(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		maxLen int
		want   string
	}{
		{"short body", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"truncated", "hello world", 5, "hello... (truncated)"},
		{"empty", "", 10, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateBody([]byte(tt.body), tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateBody(%q, %d) = %q, want %q", tt.body, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestMapAnyToMapStringJSON(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]interface{}
		want map[string]string
	}{
		{
			"string values pass through directly",
			map[string]interface{}{"key": "value", "enabled": "true"},
			map[string]string{"key": "value", "enabled": "true"},
		},
		{
			"non-string values are JSON-marshaled",
			map[string]interface{}{"count": float64(42), "active": true},
			map[string]string{"count": "42", "active": "true"},
		},
		{
			"nil map returns empty map",
			nil,
			map[string]string{},
		},
		{
			"empty map returns empty map",
			map[string]interface{}{},
			map[string]string{},
		},
		{
			"mixed types",
			map[string]interface{}{"name": "test", "port": float64(8080), "debug": false},
			map[string]string{"name": "test", "port": "8080", "debug": "false"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapAnyToMapStringJSON(nil, tt.in)
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for k, wantV := range tt.want {
				if gotV, ok := got[k]; !ok {
					t.Errorf("missing key %q", k)
				} else if gotV != wantV {
					t.Errorf("key %q = %q, want %q", k, gotV, wantV)
				}
			}
		})
	}
}

func TestLogger_NilSafety(t *testing.T) {
	// Nil logger should not panic
	c := &Client{Logger: nil}
	c.logDebug("test")
	c.logInformation("test")
	c.logWarning("test")
	c.logError("test")

	// Logger with nil functions should not panic
	c.Logger = &ClientLogger{}
	c.logDebug("test")
	c.logInformation("test")
	c.logWarning("test")
	c.logError("test")
}

func TestLogger_CallsFunctions(t *testing.T) {
	var debugMsg, infoMsg, warnMsg, errMsg string
	c := &Client{
		Logger: &ClientLogger{
			Debug:       func(s string) { debugMsg = s },
			Information: func(s string) { infoMsg = s },
			Warning:     func(s string) { warnMsg = s },
			Error:       func(s string) { errMsg = s },
		},
	}

	c.logDebug("d")
	c.logInformation("i")
	c.logWarning("w")
	c.logError("e")

	if debugMsg != "d" {
		t.Errorf("debug = %q, want %q", debugMsg, "d")
	}
	if infoMsg != "i" {
		t.Errorf("info = %q, want %q", infoMsg, "i")
	}
	if warnMsg != "w" {
		t.Errorf("warning = %q, want %q", warnMsg, "w")
	}
	if errMsg != "e" {
		t.Errorf("error = %q, want %q", errMsg, "e")
	}
}
