## Reward-Based Credit System (Go)

A minimal Go service that demonstrates a clean architecture for a reward-based credit system. It exposes an HTTP API (Gin), uses PostgreSQL via sqlx, and includes Makefile targets for local development and database migrations.

### Prerequisites
- Go 1.21+ (or compatible)
- PostgreSQL 13+ (local or remote)
- make (optional but recommended)

### Configure (env vars)
The app reads configuration from environment variables. At minimum, `DB_PASSWORD` is required.

```bash
export SERVER_HOST=localhost
export SERVER_PORT=8080

export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword   # required
export DB_NAME=srbcs
export DB_SSL_MODE=disable
```

### Install deps
```bash
go mod download
```

### Run the API
Use make (recommended):
```bash
make run
```
Or directly:
```bash
go run ./cmd/api
```
The server starts on `${SERVER_HOST}:${SERVER_PORT}` (default `localhost:8080`).

### Database migrations
Migrations live in `migrations/` and are managed with `golang-migrate` (auto-installed by the Makefile target if missing).

Apply latest migrations:
```bash
make db-migrate-up
```

Revert migrations (one by one to base):
```bash
make db-migrate-down
```

#### Run migrations without make
Install migrate (choose one):
```bash
brew install golang-migrate            # macOS (Homebrew)
# or
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Set DB URL and run:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=srbcs
export DB_SSL_MODE=disable
export DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

# apply all up migrations
migrate -path migrations -database "$DATABASE_URL" up

# rollback one step
migrate -path migrations -database "$DATABASE_URL" down 1

# show current version
migrate -path migrations -database "$DATABASE_URL" version
```

### Useful make targets
```bash
make deps          # download and tidy modules
make build         # build binary to ./bin/api
make run           # run the API locally
make db-migrate-up # apply migrations
```

---
For more details, check the source in `cmd/`, `internal/`, and `pkg/` directories.


