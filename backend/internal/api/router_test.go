package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"easy-arbitra/backend/internal/api/handler"
)

func TestHealthzRoute(t *testing.T) {
	h := handler.New(nil, nil, nil, nil, nil, nil)
	r := NewRouter(h, "http://localhost:3000")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
