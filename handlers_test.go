package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
)

var testStore *Store

func TestMain(m *testing.M) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://localhost:5432/url_shortener"
	}

	var err error
	testStore, err = NewStore(databaseURL)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}
	defer testStore.Close()

	os.Exit(m.Run())
}

func cleanupURLs(t *testing.T) {
	t.Helper()
	_, err := testStore.pool.Exec(context.Background(), "DELETE FROM urls")
	if err != nil {
		t.Fatalf("failed to clean urls table: %v", err)
	}
	// Reset the sequence so base62 codes are predictable
	_, err = testStore.pool.Exec(context.Background(), "ALTER SEQUENCE urls_id_seq RESTART WITH 1")
	if err != nil {
		t.Fatalf("failed to reset sequence: %v", err)
	}
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/shorten", handleShorten(testStore, "http://localhost:8080"))
	r.Get("/api/urls", handleListURLs(testStore))
	r.Delete("/api/urls/{code}", handleDeleteURL(testStore))
	r.Get("/{code}", handleRedirect(testStore))
	return r
}

func TestShortenURL(t *testing.T) {
	cleanupURLs(t)
	r := setupRouter()

	body := bytes.NewBufferString(`{"url": "https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp ShortenResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.OriginalURL != "https://example.com" {
		t.Errorf("expected original_url https://example.com, got %s", resp.OriginalURL)
	}
	if resp.ShortCode == "" {
		t.Error("expected non-empty short_code")
	}
	if resp.ShortURL == "" {
		t.Error("expected non-empty short_url")
	}
}

func TestShortenURL_EmptyURL(t *testing.T) {
	r := setupRouter()

	body := bytes.NewBufferString(`{"url": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestRedirect(t *testing.T) {
	cleanupURLs(t)
	r := setupRouter()

	// Create a URL first
	url, err := testStore.CreateURL(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("failed to create url: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/"+url.ShortCode, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusMovedPermanently {
		t.Fatalf("expected status 301, got %d", w.Code)
	}

	location := w.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("expected redirect to https://example.com, got %s", location)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestListURLs(t *testing.T) {
	cleanupURLs(t)
	r := setupRouter()

	// Create two URLs
	testStore.CreateURL(context.Background(), "https://example.com")
	testStore.CreateURL(context.Background(), "https://google.com")

	req := httptest.NewRequest(http.MethodGet, "/api/urls", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var urls []URL
	json.NewDecoder(w.Body).Decode(&urls)

	if len(urls) != 2 {
		t.Errorf("expected 2 urls, got %d", len(urls))
	}
}

func TestDeleteURL(t *testing.T) {
	cleanupURLs(t)
	r := setupRouter()

	url, _ := testStore.CreateURL(context.Background(), "https://example.com")

	req := httptest.NewRequest(http.MethodDelete, "/api/urls/"+url.ShortCode, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}

	// Verify it's gone
	_, err := testStore.GetURL(context.Background(), url.ShortCode)
	if err == nil {
		t.Error("expected url to be deleted")
	}
}

func TestDeleteURL_NotFound(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/api/urls/nonexistent", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}
