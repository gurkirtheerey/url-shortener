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

# Run database migrations
go run . migrate

# Run the server
go run .
```

### Run the Dashboard

```bash
cd dashboard
npm ci
NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev
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

## Render Deployment

This repo now includes a `render.yaml` blueprint that provisions:

1. A managed Postgres database
2. A Go web service for the API
3. A Node web service for the Next.js dashboard

The API service runs `./bin/server migrate` as its Render `preDeployCommand`, so schema changes are applied before each deploy.

Automatic Render deploys are disabled in `render.yaml`. Production deploys are triggered by GitHub Actions only after the `CI` workflow passes on `main`.

### Render Setup

1. Push this repo to GitHub.
2. In Render, create a new Blueprint and point it at the repo.
3. Review `render.yaml` and adjust service names, plans, regions, or domains if needed.
4. Sync the Blueprint to create the database, API, and dashboard.
5. In GitHub Actions, add these repository secrets:
   - `RENDER_API_DEPLOY_HOOK_URL`
   - `RENDER_DASHBOARD_DEPLOY_HOOK_URL`
6. In Render, open each web service and copy its deploy hook URL into the matching GitHub secret.
7. Point your DNS records at the Render services for `api.linksmith.cc` and `app.linksmith.cc`.

### Production Cutover

1. Export your current Postgres data from the existing production host.
2. Restore it into the new Render Postgres database.
3. Merge to `main` and let GitHub Actions trigger both Render deploy hooks after CI passes.
4. Verify `/health`, URL creation, redirects, and dashboard reads.
5. Decommission any remaining host-based deploy infrastructure once traffic is cut over.
