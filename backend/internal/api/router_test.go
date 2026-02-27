package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"easy-arbitra/backend/internal/api/handler"
	"easy-arbitra/backend/internal/auth"
	"easy-arbitra/backend/config"
)

func TestHealthzRoute(t *testing.T) {
	h := handler.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	authHandler := auth.NewHandler(nil, config.AuthConfig{JWTSecret: "test-secret", FrontendURL: "http://localhost:3000"})
	r := NewRouter(h, authHandler, "test-secret", "http://localhost:3000")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
