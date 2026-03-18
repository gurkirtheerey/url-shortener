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
		databaseURL = "postgres://localhost:5432/url_shortener_test"
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
	// Delete clicks first because they have a foreign key referencing urls
	_, err := testStore.pool.Exec(context.Background(), "DELETE FROM clicks")
	if err != nil {
		t.Fatalf("failed to clean clicks table: %v", err)
	}
	_, err = testStore.pool.Exec(context.Background(), "DELETE FROM urls")
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
	r.Get("/health", handleHealthCheck(testStore))
	r.Post("/api/shorten", handleShorten(testStore, "http://localhost:8080"))
	r.Get("/api/urls", handleListURLs(testStore))
	r.Get("/api/urls/{code}/stats", handleGetStats(testStore))
	r.Delete("/api/urls/{code}", handleDeleteURL(testStore))
	r.Get("/{code}", handleRedirect(testStore))
	return r
}

func TestHealthCheck(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["status"] != "ok" {
		t.Errorf("expected status ok, got %s", resp["status"])
	}
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

	if w.Code != http.StatusFound {
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

func TestListURLs_WithClickCounts(t *testing.T) {
	cleanupURLs(t)
	r := setupRouter()

	url1, _ := testStore.CreateURL(context.Background(), "https://example.com")
	testStore.CreateURL(context.Background(), "https://google.com")

	testStore.RecordClick(context.Background(), url1.ShortCode, "127.0.0.1", "TestBrowser", "")
	testStore.RecordClick(context.Background(), url1.ShortCode, "127.0.0.1", "TestBrowser", "")

	req := httptest.NewRequest(http.MethodGet, "/api/urls", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var urls []URL
	json.NewDecoder(w.Body).Decode(&urls)

	// Ordered newest first: google.com (0 clicks), example.com (2 clicks)
	if urls[0].ClickCount != 0 {
		t.Errorf("expected 0 clicks for google.com, got %d", urls[0].ClickCount)
	}
	if urls[1].ClickCount != 2 {
		t.Errorf("expected 2 clicks for example.com, got %d", urls[1].ClickCount)
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

func TestRedirect_RecordsClick(t *testing.T) {
	cleanupURLs(t)
	r := setupRouter()

	url, _ := testStore.CreateURL(context.Background(), "https://example.com")

	// Hit the redirect endpoint with a user agent and referrer
	req := httptest.NewRequest(http.MethodGet, "/"+url.ShortCode, nil)
	req.Header.Set("User-Agent", "TestBrowser/1.0")
	req.Header.Set("Referer", "https://twitter.com")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("expected status 301, got %d", w.Code)
	}

	// Verify a click was recorded
	var clickCount int
	err := testStore.pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM clicks WHERE short_code = $1", url.ShortCode,
	).Scan(&clickCount)
	if err != nil {
		t.Fatalf("failed to query clicks: %v", err)
	}
	if clickCount != 1 {
		t.Errorf("expected 1 click, got %d", clickCount)
	}
}

func TestGetStats(t *testing.T) {
	cleanupURLs(t)
	r := setupRouter()

	url, _ := testStore.CreateURL(context.Background(), "https://example.com")

	// Simulate 3 clicks with different referrers
	testStore.RecordClick(context.Background(), url.ShortCode, "127.0.0.1", "TestBrowser/1.0", "https://twitter.com")
	testStore.RecordClick(context.Background(), url.ShortCode, "127.0.0.1", "TestBrowser/1.0", "https://twitter.com")
	testStore.RecordClick(context.Background(), url.ShortCode, "127.0.0.1", "TestBrowser/1.0", "https://reddit.com")

	req := httptest.NewRequest(http.MethodGet, "/api/urls/"+url.ShortCode+"/stats", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var stats URLStats
	json.NewDecoder(w.Body).Decode(&stats)

	if stats.TotalClicks != 3 {
		t.Errorf("expected 3 total clicks, got %d", stats.TotalClicks)
	}
	if len(stats.DailyClicks) != 1 {
		t.Errorf("expected 1 day of clicks, got %d", len(stats.DailyClicks))
	}
	if len(stats.TopReferrers) != 2 {
		t.Errorf("expected 2 referrers, got %d", len(stats.TopReferrers))
	}
	// Twitter should be first (2 clicks) and Reddit second (1 click)
	if len(stats.TopReferrers) >= 2 {
		if stats.TopReferrers[0].Referrer != "https://twitter.com" {
			t.Errorf("expected top referrer to be twitter, got %s", stats.TopReferrers[0].Referrer)
		}
		if stats.TopReferrers[0].Count != 2 {
			t.Errorf("expected twitter count 2, got %d", stats.TopReferrers[0].Count)
		}
	}
}

func TestGetStats_NotFound(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/urls/nonexistent/stats", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestRedirect_SkipsPrefetchClicks(t *testing.T) {
	cleanupURLs(t)
	r := setupRouter()

	url, _ := testStore.CreateURL(context.Background(), "https://example.com")

	// Simulate a browser prefetch request (Purpose header)
	req := httptest.NewRequest(http.MethodGet, "/"+url.ShortCode, nil)
	req.Header.Set("Purpose", "prefetch")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", w.Code)
	}

	var clickCount int
	err := testStore.pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM clicks WHERE short_code = $1", url.ShortCode,
	).Scan(&clickCount)
	if err != nil {
		t.Fatalf("failed to query clicks: %v", err)
	}
	if clickCount != 0 {
		t.Errorf("expected 0 clicks for prefetch, got %d", clickCount)
	}
}

func TestRedirect_SkipsSpeculativeNavigation(t *testing.T) {
	cleanupURLs(t)
	r := setupRouter()

	url, _ := testStore.CreateURL(context.Background(), "https://example.com")

	// Simulate Chrome's address bar speculation: Fetch Metadata headers
	// are present but Sec-Fetch-User is missing (not user-initiated)
	req := httptest.NewRequest(http.MethodGet, "/"+url.ShortCode, nil)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	// Sec-Fetch-User intentionally omitted (speculative request)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", w.Code)
	}

	var clickCount int
	err := testStore.pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM clicks WHERE short_code = $1", url.ShortCode,
	).Scan(&clickCount)
	if err != nil {
		t.Fatalf("failed to query clicks: %v", err)
	}
	if clickCount != 0 {
		t.Errorf("expected 0 clicks for speculative navigation, got %d", clickCount)
	}

	// Now simulate a real user navigation: same Fetch Metadata headers
	// but WITH Sec-Fetch-User: ?1 (user pressed Enter)
	req2 := httptest.NewRequest(http.MethodGet, "/"+url.ShortCode, nil)
	req2.Header.Set("Sec-Fetch-Dest", "document")
	req2.Header.Set("Sec-Fetch-Mode", "navigate")
	req2.Header.Set("Sec-Fetch-Site", "none")
	req2.Header.Set("Sec-Fetch-User", "?1")
	w2 := httptest.NewRecorder()

	r.ServeHTTP(w2, req2)

	err = testStore.pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM clicks WHERE short_code = $1", url.ShortCode,
	).Scan(&clickCount)
	if err != nil {
		t.Fatalf("failed to query clicks: %v", err)
	}
	if clickCount != 1 {
		t.Errorf("expected 1 click for real navigation, got %d", clickCount)
	}
}

func TestGetStats_NoClicks(t *testing.T) {
	cleanupURLs(t)
	r := setupRouter()

	url, _ := testStore.CreateURL(context.Background(), "https://example.com")

	req := httptest.NewRequest(http.MethodGet, "/api/urls/"+url.ShortCode+"/stats", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var stats URLStats
	json.NewDecoder(w.Body).Decode(&stats)

	if stats.TotalClicks != 0 {
		t.Errorf("expected 0 total clicks, got %d", stats.TotalClicks)
	}
}
