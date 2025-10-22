package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appcfg "backend/internal/config"
)

func TestLoginAndMe(t *testing.T) {
	cfg := &appcfg.Config{
		JWTSecret: "test-secret",
		JWTTTL:    5 * time.Minute,
		DemoUser:  "admin",
		DemoPass:  "password123",
	}

	// Login
	body := []byte(`{"username":"admin","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	Login(cfg).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login expected 200 got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("token")) {
		t.Fatalf("expected token in response: %s", rec.Body.String())
	}
}
