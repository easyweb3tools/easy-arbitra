package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPClientRetriesOnServerError(t *testing.T) {
	count := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		if count < 2 {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("temporary"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	c := NewHTTPClient(ts.URL, 2*time.Second)
	var out map[string]bool
	if err := c.GetJSON(context.Background(), "", &out); err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	if !out["ok"] {
		t.Fatalf("unexpected output: %#v", out)
	}
}
