package provider

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCapabilitiesRoutes(t *testing.T) {
	t.Parallel()

	server := NewServer(LoadConfig())

	for _, path := range []string{"/capabilities", "/v1.0.0/capabilities"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s returned status %d", path, rec.Code)
		}

		var payload map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
			t.Fatalf("failed to decode %s response: %v", path, err)
		}

		if payload["providerName"] != "Tata Consulting" {
			t.Fatalf("unexpected provider name for %s: %#v", path, payload["providerName"])
		}
	}
}

func TestLoginRedirectsBackToMesheryTokenHandler(t *testing.T) {
	t.Parallel()

	server := NewServer(LoadConfig())
	source := base64.RawURLEncoding.EncodeToString([]byte("https://meshery.example.com"))

	req := httptest.NewRequest(http.MethodGet, "/login?source="+url.QueryEscape(source), nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if location == "" {
		t.Fatal("expected redirect location")
	}

	redirectURL, err := url.Parse(location)
	if err != nil {
		t.Fatalf("failed to parse redirect url: %v", err)
	}

	if redirectURL.String() == "https://meshery.example.com" {
		t.Fatal("expected redirect to Meshery token callback, got source base url")
	}

	if redirectURL.Path != "/api/user/token" {
		t.Fatalf("unexpected redirect path: %s", redirectURL.Path)
	}

	query := redirectURL.Query()
	if query.Get("token") == "" {
		t.Fatal("expected token query parameter")
	}
	if query.Get("session_cookie") == "" {
		t.Fatal("expected session_cookie query parameter")
	}

	cookies := rec.Result().Cookies()
	if len(cookies) == 0 || cookies[0].Name != "session_cookie" {
		t.Fatal("expected session cookie to be set on provider response")
	}
}

func TestProfileRequiresValidBearerToken(t *testing.T) {
	t.Parallel()

	cfg := LoadConfig()
	server := NewServer(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/identity/users/profile", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", rec.Code)
	}

	token, err := server.mintToken()
	if err != nil {
		t.Fatalf("failed to mint token: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/identity/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with bearer token, got %d", rec.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode profile response: %v", err)
	}

	if payload["userId"] != cfg.DefaultUser.UserID {
		t.Fatalf("unexpected userId: %#v", payload["userId"])
	}
}
