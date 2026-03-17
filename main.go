package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Config from environment variables, with sensible defaults for local dev.
	// In production, you'd set these via your hosting platform (Fly.io, Railway, etc.)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://localhost:5432/url_shortener"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Connect to the database. If this fails, crash immediately.
	// There's no point running a web server that can't reach its database.
	store, err := NewStore(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer store.Close() // clean up connections when the program exits

	// Chi is the HTTP router. It matches incoming requests to handler functions
	// based on the method (GET, POST, DELETE) and URL path.
	r := chi.NewRouter()

	// Middleware runs on every request before your handler.
	// Logger prints a log line for each request (method, path, status, duration).
	// Recoverer catches panics in handlers and returns a 500 instead of crashing the server.
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API routes (these serve JSON)
	r.Post("/api/shorten", handleShorten(store, baseURL))
	r.Get("/api/urls", handleListURLs(store))
	r.Get("/api/urls/{code}/stats", handleGetStats(store))
	r.Delete("/api/urls/{code}", handleDeleteURL(store))

	// Redirect route. {code} is a URL parameter, like a wildcard.
	// Any GET request that doesn't match /api/* falls through to here.
	// So GET /abc will try to look up "abc" as a short code.
	r.Get("/{code}", handleRedirect(store))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on :%s\n", port)

	// ListenAndServe blocks forever, listening for incoming HTTP requests.
	// It only returns if the server crashes, which log.Fatal will print and exit.
	log.Fatal(http.ListenAndServe(":"+port, r))
}
