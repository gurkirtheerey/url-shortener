# URL Shortener with Analytics Pipeline

A production-ready URL shortener built in phases, starting with a Go API and Postgres, and growing into a multi-service system with a Python analytics pipeline, Next.js dashboard, Redis caching, and Docker orchestration.

## Architecture

The final system has three services and two data stores:

**Services:**
1. **Go API** - creates short URLs and handles redirects
2. **Python Analytics Workers** - processes click data in the background (geo lookup, device parsing, aggregation)
3. **Next.js Dashboard** - visualizes analytics

**Data Stores:**
- **Postgres** - source of truth for URLs, clicks, stats
- **Redis** - caching redirects, queueing click events, rate limiting

```
User clicks short link
        |
    [ Nginx ]
        |
    [ Go API ]
      |-- Redis Cache (URL lookup, sub-ms)
      |-- Redis Stream (publish click event)
      |-- Postgres (fallback URL lookup)
        |
    [ Python Worker ]
      |-- Redis Stream (consume events)
      |-- MaxMind (geo enrichment)
      |-- Postgres (batch write clicks + daily aggregation)
        |
    [ Next.js Dashboard ]
      |-- Postgres (read aggregated stats via Go API)
```

## Phases

The project is built incrementally. Each phase is a complete, working product.

### Phase 1: Core Shortener (Go + Postgres)
Base62-encoded short codes, CRUD API, redirect handling. No auth, no analytics, no caching.

**API Endpoints:**
```
POST /api/shorten       - Create a short URL
GET  /:code             - Redirect to original URL
GET  /api/urls          - List all shortened URLs
DELETE /api/urls/:code  - Delete a shortened URL
```

### Phase 2: Simple Analytics (Go + Postgres)
Click logging on every redirect. Stats endpoint with click counts, daily breakdown, and top referrers.

### Phase 3: Dashboard (Next.js + TypeScript)
Frontend with URL creation form, links table, and per-link analytics charts (Recharts, Tailwind, shadcn/ui).

### Phase 4: Docker + Infrastructure
Dockerfiles for each service, Docker Compose for the full stack, Nginx reverse proxy, GitHub Actions CI/CD.

### Phase 5: Redis + Performance
Cache-aside pattern for redirects, rate limiting middleware, Redis added to the stack.

### Phase 6: Analytics Pipeline (Python + Redis Streams)
Click events published to Redis Streams instead of written to Postgres during redirect. Python worker consumes events, enriches with geo/device data (MaxMind, ua-parser), batch inserts to Postgres, and runs scheduled aggregation into a daily summary table.

## Tech Stack

| Component | Technology | Why |
|-----------|-----------|-----|
| API | Go + Chi | Lightweight, idiomatic, good middleware support |
| Database | Postgres (pgx driver) | Most performant Go driver, actively maintained |
| Short codes | Base62 encoding | Deterministic, collision-free |
| Frontend | Next.js + TypeScript | Recharts for charts, Tailwind + shadcn for UI |
| Cache | Redis (cache-aside) | Industry standard pattern |
| Message queue | Redis Streams | Avoids adding another service at this scale |
| Geo lookup | MaxMind GeoLite2 | Free, local file, no API dependency |
| Orchestration | Docker Compose | Standard for multi-service local dev |

## Getting Started

### Prerequisites
- Go 1.25+
- Postgres running locally

### Run the API

```bash
# Set your database URL
export DATABASE_URL=postgres://user:pass@localhost:5432/shortener

# Create the urls table
psql $DATABASE_URL -c "
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);"

# Run the server
go run .
```

### Usage

```bash
# Shorten a URL
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

# Visit the short URL
curl -L http://localhost:8080/abc123

# List all URLs
curl http://localhost:8080/api/urls

# Delete a URL
curl -X DELETE http://localhost:8080/api/urls/abc123
```

### Run Tests

```bash
go test ./...
```
