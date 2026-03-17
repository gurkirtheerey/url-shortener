package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// URL represents a shortened URL record in the database.
// The json tags control how each field is named when serialized to JSON
// (e.g. the API response will use "short_code" instead of "ShortCode").
type URL struct {
	ID          int64     `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// Store is the database layer. It holds a connection pool and provides
// methods for all database operations. Every handler calls Store methods
// instead of writing SQL directly.
type Store struct {
	pool *pgxpool.Pool
}

// NewStore creates a new Store with a connection pool.
// A pool keeps multiple database connections open and reuses them across requests.
// This avoids the overhead of opening a new connection for every single query
// (each new connection takes a few milliseconds of handshaking).
func NewStore(databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Ping verifies the connection actually works (not just that the config is valid)
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

// CreateURL inserts a new shortened URL. It works in two steps:
// 1. Get the next auto-increment ID from Postgres's sequence
// 2. Convert that ID to a base62 short code and insert the full row
//
// We grab the ID first (via nextval) so we can generate the short code
// before inserting. This avoids needing a placeholder value or a second UPDATE.
func (s *Store) CreateURL(ctx context.Context, originalURL string) (*URL, error) {
	var id int64
	err := s.pool.QueryRow(ctx, "SELECT nextval('urls_id_seq')").Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to get next id: %w", err)
	}

	shortCode := encodeBase62(id)
	var createdAt time.Time

	// RETURNING lets us get back auto-generated column values (like created_at)
	// without needing a separate SELECT query.
	err = s.pool.QueryRow(ctx,
		"INSERT INTO urls (id, short_code, original_url) VALUES ($1, $2, $3) RETURNING created_at",
		id, shortCode, originalURL,
	).Scan(&createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert url: %w", err)
	}

	return &URL{
		ID:          id,
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   createdAt,
	}, nil
}

// GetURL looks up a single URL by its short code.
// If the code doesn't exist, pgx returns pgx.ErrNoRows, which the handler
// uses to return a 404.
func (s *Store) GetURL(ctx context.Context, shortCode string) (*URL, error) {
	var u URL
	err := s.pool.QueryRow(ctx,
		"SELECT id, short_code, original_url, created_at FROM urls WHERE short_code = $1",
		shortCode,
	).Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ListURLs returns all shortened URLs, newest first.
// rows.Next() iterates through results one at a time, and Scan maps each
// row's columns into our URL struct fields (order matters here, it must
// match the SELECT column order).
func (s *Store) ListURLs(ctx context.Context) ([]URL, error) {
	rows, err := s.pool.Query(ctx,
		"SELECT id, short_code, original_url, created_at FROM urls ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // always close rows when done, or you leak connections

	var urls []URL
	for rows.Next() {
		var u URL
		if err := rows.Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.CreatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

// DeleteURL removes a URL by short code.
// RowsAffected tells us if the DELETE actually matched a row. If it's 0,
// the short code didn't exist, so we return an error (which the handler
// turns into a 404).
func (s *Store) DeleteURL(ctx context.Context, shortCode string) error {
	result, err := s.pool.Exec(ctx, "DELETE FROM urls WHERE short_code = $1", shortCode)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("url not found")
	}
	return nil
}
