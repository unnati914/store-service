# Store Service (Ephemeral Stories Backend)

Production-style backend scaffold for ephemeral Stories: text + optional media, visibility (public/friends/private), expiring after 24h. Built in Go, with Postgres, Redis, and MinIO.

## Architecture
- API (Go): REST endpoints, JWT auth, rate limits (planned), Redis caching (planned), Prometheus metrics at `/metrics`.
- Worker (Go): background ephemerality scanner (stubbed; planned to soft-delete expired stories and emit logs/metrics).
- Postgres: primary data store.
- Redis: caching, rate limiting, and pub/sub for real-time (planned).
- MinIO (S3-compatible): object storage for media (presigned uploads; planned).

## Requirements
- Go 1.22+
- Docker + Docker Compose (recommended), or
- macOS (Homebrew), Postgres 16, Redis (for local dev without Docker)

## Configuration (.env)
Copy and edit as needed:
```
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/store?sslmode=disable
REDIS_ADDR=localhost:6379
JWT_SECRET=change-me
S3_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=stories
```

## Run with Docker (recommended)
```
cp .env.example .env
make docker-up

# Apply DB migrations
docker compose exec -T postgres \
  psql -U postgres -d store -f internal/db/migrations.sql

# Health
curl -s localhost:8080/healthz
# Metrics
curl -s localhost:8080/metrics | head
```
- Logs: `docker compose logs -f api`
- Stop: `make docker-down`

## Local development (no Docker)
```
# macOS setup
brew install postgresql@16 redis
brew services start postgresql@16
brew services start redis
# ensure psql is on PATH (Apple Silicon default)
echo 'export PATH="/opt/homebrew/opt/postgresql@16/bin:$PATH"' >> ~/.zshrc && source ~/.zshrc

# Env
cp .env.example .env

# DB + migrations
createdb store || true
psql "postgres://postgres@localhost:5432/store?sslmode=disable" -f internal/db/migrations.sql

# Run API
make dev  # or: go run ./cmd/api

# Health
curl -s localhost:8080/healthz
```

## Endpoints (current)
- POST `/signup`
  - Body: `{ "email": string, "password": string }`
  - 201 Created → `{ "token": string }` (JWT)
- POST `/login`
  - Body: `{ "email": string, "password": string }`
  - 200 OK → `{ "token": string }`
- GET `/healthz`
  - 200 OK when API can reach Postgres and Redis
- GET `/metrics`
  - Prometheus metrics (e.g., `http_requests_total{route,code}`)

Notes:
- JWT secret is `JWT_SECRET`. Claims include `user_id`.
- Passwords are hashed with bcrypt.

## Roadmap (planned in this repo)
- Presigned S3/MinIO uploads with content-type/size validation
- Stories CRUD + feed, views (idempotent), reactions, RBAC/permissions
- Social graph (follow/unfollow) and friends visibility policy
- Redis token-bucket rate limits (POST /stories, POST /reactions)
- Redis caching: followee IDs, hot feed page
- WebSocket/SSE real-time events for views/reactions
- Ephemerality worker (soft-delete expired stories) + logs/metrics
- Observability: structured JSON logs; Prometheus metrics; /healthz checks
- Tests: unit tests for handlers; integration tests with containers
- Optional: search via Postgres FTS (GIN)

## Make targets
```
make dev           # run API locally
make build         # build all packages
make test          # run tests (once added)
make lint          # run linter (if installed)
make docker-up     # compose up (api, worker, postgres, redis, minio)
make docker-down   # compose down -v
```

## Quick walkthrough
1) Sign up & login → JWT
```
curl -s -X POST localhost:8080/signup -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"hunter2!!"}'

curl -s -X POST localhost:8080/login -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"hunter2!!"}'
```
2) Health & Metrics
```
curl -s localhost:8080/healthz
curl -s localhost:8080/metrics | head
```
3) (Upcoming) Presigned upload → Create story → Feed → View/React → Real-time → Worker expiry

## Repository structure
```
cmd/
  api/       # API service entrypoint
  worker/    # worker process entrypoint (expiration scanner - stub)
internal/
  auth/      # JWT + password hashing
  config/    # env config loader
  db/        # Postgres + Redis clients, migrations
  handlers/  # HTTP handlers (auth, etc.)
  logging/   # JSON structured logger (slog)
  metrics/   # Prometheus metrics
  models/    # data access (UserStore, etc.)
  server/    # HTTP server, middleware (logging + metrics)
```

## Troubleshooting
- `psql: command not found`: install Postgres client and add to PATH.
- Health shows down for DB/Redis: ensure services are running and env vars point to them.
- Docker Desktop required for `make docker-up`.