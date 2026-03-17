package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// Request/response types. These structs define the shape of JSON the API
// accepts and returns. The json tags map Go field names to JSON keys.
type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL    string `json:"short_url"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// handleShorten returns a handler function that creates a new short URL.
//
// This pattern (a function that returns a function) is called a "closure."
// The inner function needs access to `store` and `baseURL`, but Go's
// http.HandlerFunc signature only gives you (w, r). So we wrap it:
// the outer function captures store/baseURL, and the inner function
// uses them when handling requests.
func handleShorten(store *Store, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the JSON request body into our struct
		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
			return
		}

		if req.URL == "" {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "url is required"})
			return
		}

		// If the user didn't include a protocol, prepend https://.
		// Without this, "gsheerey.com" would redirect to "http://localhost:8080/gsheerey.com"
		// instead of the actual website.
		if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
			req.URL = "https://" + req.URL
		}

		url, err := store.CreateURL(r.Context(), req.URL)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to create short url"})
			return
		}

		// 201 Created is the standard status code for "new resource created"
		writeJSON(w, http.StatusCreated, ShortenResponse{
			ShortURL:    baseURL + "/" + url.ShortCode,
			ShortCode:   url.ShortCode,
			OriginalURL: url.OriginalURL,
		})
	}
}

// handleRedirect looks up a short code and sends the user to the original URL.
// This is the core feature. When someone visits localhost:8080/abc, Chi
// extracts "abc" as the URL parameter "code", we look it up, and return
// a 301 (permanent redirect) to the original URL.
//
// Phase 2 addition: we also log the click (IP, user agent, referrer) before
// redirecting. If logging fails, we still redirect. Analytics shouldn't
// break the user experience.
func handleRedirect(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		url, err := store.GetURL(r.Context(), code)
		if err != nil {
			if err == pgx.ErrNoRows {
				writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "short url not found"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to look up url"})
			return
		}

		// Extract the client's IP address. RemoteAddr includes the port (e.g. "127.0.0.1:52341"),
		// so we split it off. This is the IP of whoever made the request.
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		// Log the click. If this fails, we just log the error and continue with the redirect.
		if err := store.RecordClick(r.Context(), code, ip, r.UserAgent(), r.Referer()); err != nil {
			log.Printf("failed to record click for %s: %v", code, err)
		}

		// 302 (temporary redirect) instead of 301 (permanent). A 301 tells the browser
		// to cache the redirect forever, so if the mapping changes (or gets deleted),
		// the browser would still redirect to the old URL without hitting our server.
		// 302 forces the browser to check with us every time.
		http.Redirect(w, r, url.OriginalURL, http.StatusFound)
	}
}

// handleGetStats returns analytics for a single short URL.
// It checks that the URL exists first (404 if not), then returns
// total clicks, daily breakdown, and top referrers.
func handleGetStats(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		// Verify the short code exists before querying stats
		if _, err := store.GetURL(r.Context(), code); err != nil {
			if err == pgx.ErrNoRows {
				writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "short url not found"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to look up url"})
			return
		}

		stats, err := store.GetURLStats(r.Context(), code)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to get stats"})
			return
		}

		writeJSON(w, http.StatusOK, stats)
	}
}

func handleListURLs(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls, err := store.ListURLs(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to list urls"})
			return
		}

		writeJSON(w, http.StatusOK, urls)
	}
}

func handleDeleteURL(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		if err := store.DeleteURL(r.Context(), code); err != nil {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "url not found"})
			return
		}

		// 204 No Content means "success, but there's nothing to send back"
		w.WriteHeader(http.StatusNoContent)
	}
}

// writeJSON is a helper that every handler uses to send JSON responses.
// It handles three things that are easy to forget:
// 1. Set Content-Type so the client knows it's JSON
// 2. Set the status code (200, 201, 400, 404, 500, etc.)
// 3. Encode the data as JSON into the response body
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
