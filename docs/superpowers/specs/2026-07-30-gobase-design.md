# GoBase — Go Project Template Design

## Overview

A batteries-included Go project template inspired by Laravel's `artisan`. One binary with subcommands (`serve`, `migrate`, `jwt:generate`), a chi-based HTTP API, sqlx-backed database access with golang-migrate migrations, HMAC JWT authentication, and structured logging via `slog`. Clone, run `./setup.sh`, start building.

## Tech Stack

| Concern | Choice | Rationale |
|---|---|---|
| HTTP router | `chi` | `net/http` compatible, clean route groups, whole stdlib ecosystem |
| CLI | `cobra` | Industry standard, subcommand nesting, built-in help |
| Database access | `sqlx` | Raw SQL with struct scanning, minimal magic |
| Migrations | `golang-migrate` | SQL file-based, CLI integration, no ORM coupling |
| Auth | HMAC-SHA256 JWT | Simple shared secret, stateless verification |
| Config | Env vars + `.env` | Standard 12-factor pattern |
| Logging | `slog` | Stdlib structured logging |
| Validation | `go-playground/validator` | Struct tag validation, standard in Go |
| Module rename | `setup.sh` | Interactive bootstrap script |

## Project Structure

```
gobase/
├── cmd/                     # CLI subcommands (cobra)
│   ├── root.go              # Root command, global flags, help
│   ├── serve.go             # serve — starts the HTTP API
│   ├── migrate.go           # migrate:{up,down,create,status}
│   └── jwt.go               # jwt:generate — creates signed tokens
├── internal/
│   ├── config/              # Env + .env loading, typed Config struct
│   │   └── config.go
│   ├── database/            # sqlx connection pool + migration runner
│   │   ├── database.go
│   │   └── migrate.go
│   ├── auth/                # JWT sign/verify + chi middleware
│   │   ├── jwt.go
│   │   └── middleware.go
│   ├── log/                 # slog setup (JSON/prod, text/dev)
│   │   └── log.go
│   ├── server/              # Router assembly, graceful shutdown
│   │   ├── server.go
│   │   ├── routes.go
│   │   └── middleware.go    # Request ID, request logging, recovery
│   └── user/                # Working domain module — the template pattern
│       ├── model.go         # Types: User, CreateUserRequest, etc.
│       ├── handler.go       # HTTP handlers + route registration
│       ├── service.go       # Business logic
│       └── repository.go    # sqlx data access
├── migrations/              # golang-migrate SQL files
├── .env.example             # Documented config template
├── setup.sh                 # Bootstrap: module rename, env copy, go mod tidy
├── Makefile                 # build, run, test, lint, migrate targets
├── go.mod                   # Placeholder: github.com/wesleysnt/gobase
└── main.go                  # Single entry point: cmd.Execute()
```

### Why `internal/` for everything

Go's `internal/` package boundary prevents any code outside the module from importing these packages. Application code (users, posts, billing) belongs here — it's specific to this app. Reusable libraries that should be shared across projects belong in their own module or a top-level non-`internal` package. The `internal/` boundary is the Go equivalent of "you don't know it's a library until you need it to be one."

### Domain module pattern

Each domain lives as a flat folder under `internal/` containing all its layers:

```
internal/user/
├── model.go       # Types: User, CreateUserRequest, UserResponse
├── handler.go     # HTTP handlers + RegisterRoutes(r, db, cfg)
├── service.go     # Business logic, validation, orchestration
└── repository.go  # sqlx queries, scanning
```

This is vertical slicing by feature, not horizontal by layer. Adding a new domain means creating one new folder. Removing one means `rm -rf internal/<domain>/`. Infrastructure that spans all domains (config, DB, auth, log, server) stays as shared packages.

The `user/` module is a fully working CRUD domain (with its own migration), not a stub. It serves both as a usable user-management starting point and as the reference pattern for adding new domains. Keep it or delete it — either way, it shows you the conventions.

## CLI Design

Single binary with cobra subcommands:

```
gobase
├── serve              # Start the HTTP API server
│   └── --port, --env
├── migrate
│   ├── up             # Run pending migrations
│   ├── down           # Rollback last N migrations
│   ├── create <name>  # Create migration file pair (NNNNNN_name.up/down.sql)
│   └── status         # Show applied/pending migrations
├── jwt
│   └── generate       # Generate a signed JWT for a user ID
│       ├── --user-id (required)
│       ├── --expires-in (default: 24h)
│       └── --claims (extra claims JSON)
└── help               # Built into cobra
```

- `main.go` calls `cmd.Execute()` — the server is just another subcommand
- Running `gobase` with no args prints help
- Config loading: cobra flags > env vars > `.env` file
- `serve` loads the full config; `migrate` only loads `DATABASE_URL`; `jwt:generate` only needs `JWT_SECRET`

## CLI Command Specifications

### `serve`

Starts the HTTP API server with graceful shutdown on SIGINT/SIGTERM.

```
Flags:
  --port int    Server port (default: from PORT env or 8080)
  --env string  Environment: development|production|test (default: from ENV or "development")
```

Behavior:
- Loads full Config from env/.env
- Connects to database (exits with error if unreachable)
- Assembles chi router via `server.NewRouter()`
- Starts `http.Server` with Read/Write/Idle timeouts
- Listens for SIGINT/SIGTERM; drains connections within a deadline (default 30s)

### `migrate up`

Runs all pending migrations.

```
Flags:
  --steps int  Max migrations to apply (0 = all)
```

Behavior:
- Connects to database using `DATABASE_URL`
- Creates the schema_migrations table if absent (golang-migrate default)
- Runs pending `.up.sql` files in order
- Prints each migration applied

### `migrate down`

Rolls back the last N migrations (default 1).

```
Flags:
  --steps int  Number of migrations to roll back (default: 1)
```

Behavior:
- Same connection setup as up
- Runs `.down.sql` files in reverse order

### `migrate create <name>`

Creates a timestamped migration file pair in `migrations/`.

```
Args:
  <name>  Migration name in snake_case (e.g., "add_email_to_users")

Behavior:
  Creates migrations/NNNNNN_add_email_to_users.up.sql
  Creates migrations/NNNNNN_add_email_to_users.down.sql
  Timestamp prefix is Unix epoch seconds
```

### `migrate status`

Prints a table of all migrations with their applied/dirty status. Exits with code 1 if any migration is in a dirty state.

### `jwt:generate`

Generates a signed JWT for development and testing.

```
Flags:
  --user-id string     Subject claim (required)
  --expires-in duration Token expiry (default: 24h)
  --claims string      Extra claims as JSON object (optional)

Example:
  gobase jwt generate --user-id="usr_123" --expires-in=72h
  gobase jwt generate --user-id="usr_123" --claims='{"role":"admin"}'
```

Behavior:
- Reads `JWT_SECRET` from config
- Builds claims: sub, iat, exp, plus any --claims extras
- Signs with HMAC-SHA256
- Prints the token string to stdout

## API Design

### Router assembly

`internal/server/routes.go` owns the top-level router. Each domain registers itself:

```go
func NewRouter(cfg *config.Config, db *sqlx.DB, log *slog.Logger) chi.Router {
    r := chi.NewRouter()

    // Global middleware
    r.Use(middleware.RequestID)
    r.Use(serverMiddleware.RequestLogger(log))
    r.Use(middleware.Recoverer)

    r.Route("/api/v1", func(r chi.Router) {
        user.RegisterRoutes(r, db, cfg)
    })

    // Health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
    })

    return r
}
```

### Domain route registration

Each domain exposes one `RegisterRoutes` function:

```go
func RegisterRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config) {
    svc := &UserService{repo: &UserRepo{db}}
    h := &Handler{svc}

    r.Group(func(r chi.Router) {
        r.Get("/users", h.List)
        r.Post("/users", h.Create)
        r.Get("/users/{id}", h.GetByID)
        r.Patch("/users/{id}", h.Update)
        r.Delete("/users/{id}", h.Delete)
    })

    r.Group(func(r chi.Router) {
        r.Use(auth.RequireAuth(cfg))
        r.Get("/me", h.Me)
    })
}
```

### Handler signature pattern

Plain `http.HandlerFunc` — no framework binding. Each handler follows a 4-step flow:

```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    // 1. Parse
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    // 2. Validate
    if err := validate.Struct(req); err != nil {
        writeError(w, http.StatusUnprocessableEntity, err.Error())
        return
    }
    // 3. Call service
    user, err := h.svc.Create(r.Context(), req)
    if err != nil {
        writeError(w, codeFrom(err), err.Error())
        return
    }
    // 4. Respond
    writeJSON(w, http.StatusCreated, user)
}
```

### Response helpers

Two small functions in `internal/server/`:

- `writeJSON(w, statusCode, data)` — sets Content-Type, writes status, encodes JSON
- `writeError(w, statusCode, message)` — writes `{"error": "message"}`

### Error mapping

```go
// internal/server/errors.go
func codeFrom(err error) int {
    switch {
    case errors.Is(err, ErrNotFound):    return 404
    case errors.Is(err, ErrValidation):  return 422
    case errors.Is(err, ErrConflict):    return 409
    case errors.Is(err, ErrForbidden):   return 403
    case errors.Is(err, ErrUnauthorized): return 401
    default:                             return 500
    }
}
```

Services define and return sentinel errors. Unknown errors default to 500 and are logged at ERROR level.

## Database

### Connection

`internal/database/database.go`:

```go
func Connect(cfg *config.Config) (*sqlx.DB, error)
```

- Opens via `sqlx.Connect(driver, dsn)` using `DATABASE_URL` from config
- Driver inferred from the URL scheme: `postgres://` → pgx, `mysql://` → mysql, `sqlite://` or `file:` → sqlite3. Unrecognized schemes produce a clear error listing the supported schemes.
- Sets pool defaults: MaxOpenConns=25, MaxIdleConns=5, ConnMaxLifetime=5min (overridable in config)
- Pings once to verify connectivity on startup

### Migrations

`internal/database/migrate.go` uses `golang-migrate/migrate` with the `iofs` source driver. Migration files are embedded via `//go:embed` so the binary is self-contained:

```go
//go:embed migrations/*.sql
var migrationsFS embed.FS
```

SQL files in `migrations/`:

```
migrations/
├── 000001_create_users.up.sql
├── 000001_create_users.down.sql
├── 000002_add_email_to_users.up.sql
└── 000002_add_email_to_users.down.sql
```

### Repository pattern

Concrete structs, not interfaces. `sqlx` for scanning:

```go
type UserRepo struct {
    db *sqlx.DB
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*User, error) {
    var u User
    err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id = $1", id)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user %s: %w", id, ErrNotFound)
    }
    return &u, err
}

func (r *UserRepo) Create(ctx context.Context, u *User) error {
    query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at, updated_at`
    return r.db.QueryRowContext(ctx, query, u.Name, u.Email).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}
```

- `context.Context` as first parameter in every method for cancellation/deadline propagation
- `$1, $2` PostgreSQL placeholders by default (documented how to swap for MySQL `?`)
- Returns Go errors directly; service layer wraps them into domain sentinel errors

## Auth & JWT

### Token operations (`internal/auth/jwt.go`)

```go
func GenerateToken(userID string, expiry time.Duration, extra map[string]interface{}) (string, error)
func ParseToken(tokenStr string, secret []byte) (*Claims, error)
```

- HMAC-SHA256 signing with `JWT_SECRET` from config
- Claims struct: `Sub` (user ID), `Iat`, `Exp`, plus a `map[string]interface{}` for extras
- `ParseToken` validates signature AND expiry

### Middleware (`internal/auth/middleware.go`)

```go
func RequireAuth(cfg *config.Config) func(http.Handler) http.Handler
func OptionalAuth(cfg *config.Config) func(http.Handler) http.Handler
```

- `RequireAuth`: No token → 401. Invalid/expired token → 401. Valid token → injects `*Claims` into context, calls next handler
- `OptionalAuth`: Same parsing but allows unauthenticated requests through. Claims may be nil in the handler
- Token extracted from `Authorization: Bearer <token>` header

### Context helpers

```go
type contextKey string
const claimsKey contextKey = "claims"

func GetClaims(ctx context.Context) *Claims
func SetClaims(ctx context.Context, c *Claims) context.Context
```

## Config Management

`internal/config/config.go`:

```go
type Config struct {
    Port        int           // PORT env, default 8080
    Env         string        // ENV: development|production|test
    DatabaseURL string        // DATABASE_URL (required)
    JWTSecret   string        // JWT_SECRET (required)
    JWTExpiry   time.Duration // JWT_EXPIRY, default 24h
    LogLevel    string        // LOG_LEVEL, default "info"
    LogFormat   string        // LOG_FORMAT: json|text, default text in dev, json in prod
    // Pool settings (with defaults)
    DBMaxOpenConns    int
    DBMaxIdleConns    int
    DBConnMaxLifetime time.Duration
}
```

Loading order: flags > environment variables > `.env` file. `.env` is ignored in production (`ENV=production`).

`.env.example` ships with all keys documented and placeholder values. Running `setup.sh` copies it to `.env`.

## Logging

`internal/log/log.go`:

```go
func New(cfg *config.Config) *slog.Logger
```

- `json` format when `LOG_FORMAT=json` or `ENV=production`
- `text` format in development (colored, human-readable)
- `slog.Level` controlled by `LOG_LEVEL` (debug, info, warn, error)
- Request logging middleware logs: method, path, status, duration, request ID

## Bootstrap Flow (`setup.sh`)

```
$ ./setup.sh
Project name (github.com/org/repo): github.com/wesley/newapp
Module name set to: github.com/wesley/newapp
.env copied from .env.example — edit with your settings
✓ go mod tidy
Done. Next: cd /path/to/newapp && make dev
```

The script:
1. Prompts for the new module path
2. Replaces `github.com/wesleysnt/gobase` with the new module across all `.go` files and `go.mod`
3. Copies `.env.example` to `.env`
4. Runs `go mod tidy` to verify
5. Prints next steps

## Makefile Targets

| Target | Command |
|---|---|
| `make build` | `go build -o bin/gobase .` |
| `make dev` | `go run . serve` (with hot-reload if `air` is installed) |
| `make test` | `go test ./...` |
| `make lint` | `golangci-lint run` |
| `make migrate-up` | `go run . migrate up` |
| `make migrate-down` | `go run . migrate down` |
| `make jwt` | `go run . jwt generate --user-id=...` |

## Testing Strategy

- **Unit tests**: Standard `go test` with table-driven tests. No mocking framework — use interfaces where testing requires stubs, concrete types where it doesn't.
- **Integration tests**: Repository tests against a real database (`TEST_DATABASE_URL`). Migrations run in test setup, database is cleaned between tests.
- **Handler tests**: `httptest.NewServer` with the real chi router, real middleware. Authenticated requests build a token via `auth.GenerateToken()` in test setup.
- **Test helpers**: A `internal/testutil/` package provides `SetupTestDB(t) *sqlx.DB`, `AuthHeader(userID string) string`, and test fixture helpers.

## Conventions

- **Placeholders**: `$1, $2` (PostgreSQL) by default. Swapping to MySQL `?` requires passing a `?` rebinding function to `sqlx` — documented in `internal/database/database.go`.
- **Error sentinels**: Each service package defines its own `ErrNotFound`, `ErrValidation`, etc. The `server.codeFrom(err)` function uses `errors.Is` to map them to HTTP status codes.
- **Context propagation**: Every stack layer (handler → service → repository) passes `context.Context` as the first argument. No global request state.
- **SQL in repositories, not services**: Service layer never constructs SQL. All queries live in repository methods.
- **Package naming**: Singular (`user`, `post`), not plural (`users`, `posts`). Same convention as the Go standard library (`net/http`, `database/sql`).
