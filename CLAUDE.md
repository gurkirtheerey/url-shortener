# URL Shortener

## Project Overview

A multi-phase URL shortener project. Currently in **Phase 1** (core shortener with Go + Postgres). See the Obsidian PRD at `My Vault/7 - Projects/URL Shortener with Analytics Pipeline.md` for the full roadmap.

## Current Architecture

Single Go service with Postgres. Flat file structure, no packages yet.

## Tech Stack (Phase 1)

- **Language:** Go 1.25
- **Router:** Chi (`github.com/go-chi/chi/v5`)
- **Database:** Postgres via pgx (`github.com/jackc/pgx/v5`)
- **Short codes:** Base62 encoding from auto-incrementing Postgres IDs (collision-free)

## Project Structure

```
main.go          - server setup, routing, config
handlers.go      - HTTP handler functions
db.go            - database connection and queries
base62.go        - base62 encoding/decoding
base62_test.go   - base62 unit tests
handlers_test.go - handler integration tests
```

## API Endpoints

```
POST /api/shorten       - accepts { "url": "https://..." }, returns short URL
GET  /:code             - 301 redirect to original URL
GET  /api/urls          - list all shortened URLs
DELETE /api/urls/:code  - delete a shortened URL
```

## Database Schema

```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## Commands

```bash
# Run the server
go run .

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...
```

## Environment Variables

- `DATABASE_URL` - Postgres connection string (e.g., `postgres://user:pass@localhost:5432/shortener`)

## Bruno API Collection

Bruno request files live in `bruno/url-shortener/`. When adding a new API endpoint, always create a corresponding `.bru` file in that directory so the collection stays in sync with the codebase.

## Key Design Decisions

- **Base62 encoding** over random strings: deterministic, no collision checks needed. Postgres auto-increment ID is converted to base62.
- **Chi router** over Gin: closer to the standard library, idiomatic Go.
- **pgx** over database/sql: better performance, native Postgres types.
- **Flat structure** for now: no packages until there's a reason. Keep it simple.
- **Integration tests hit a real database**, not mocks.

## Phase Roadmap

1. **Core Shortener** - Go API + Postgres
2. **Simple Analytics** - click logging, stats endpoint
3. **Dev/Prod Environments + Deployment** - separate test DB, deploy to Fly.io
4. **Dashboard** - Next.js + TypeScript frontend
5. **Docker + Infrastructure** - Dockerfiles, Compose, Nginx, CI/CD
6. **Redis + Performance** - caching, rate limiting
7. **Analytics Pipeline** - Python workers, Redis Streams, geo enrichment
