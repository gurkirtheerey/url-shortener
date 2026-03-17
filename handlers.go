package main

import (
	"encoding/json"
	"net/http"

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
func handleRedirect(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		url, err := store.GetURL(r.Context(), code)
		if err != nil {
			// pgx.ErrNoRows means the query returned zero results (code doesn't exist)
			if err == pgx.ErrNoRows {
				writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "short url not found"})
				return
			}
			// Any other error is an actual database problem
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to look up url"})
			return
		}

		// 301 tells the browser "this URL has permanently moved to the new location."
		// The browser will follow the redirect automatically.
		http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
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
