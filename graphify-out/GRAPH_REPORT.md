# Graph Report - gobase  (2026-09-04)

## Corpus Check
- Corpus is ~17,431 words - fits in a single context window. You may not need a graph.

## Summary
- 192 nodes · 441 edges · 13 communities (11 shown, 2 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 45 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- User Domain Module
- JWT Auth
- Design Doc & CLI Root
- Logging
- HTTP Middleware & Response
- Database & Migrations
- HTTP Handlers
- Config
- Error Mapping
- Setup Script
- Package Root

## God Nodes (most connected - your core abstractions)
1. `GoBase Go Project Template Design` - 21 edges
2. `Load()` - 12 edges
3. `User` - 12 edges
4. `newStubRepo()` - 12 edges
5. `setupTestRouter()` - 11 edges
6. `GoBase template` - 11 edges
7. `HMAC-SHA256 JWT authentication` - 11 edges
8. `GoBase Template Implementation Plan` - 11 edges
9. `GenerateToken()` - 10 edges
10. `NewRouter()` - 10 edges

## Surprising Connections (you probably didn't know these)
- `HMAC-SHA256 JWT authentication` --references--> `GenerateToken()`  [EXTRACTED]
  docs/superpowers/specs/2026-07-30-gobase-design.md → internal/auth/jwt.go
- `HMAC-SHA256 JWT authentication` --references--> `ParseToken()`  [EXTRACTED]
  docs/superpowers/specs/2026-07-30-gobase-design.md → internal/auth/jwt.go
- `HMAC-SHA256 JWT authentication` --references--> `RequireAuth()`  [EXTRACTED]
  docs/superpowers/specs/2026-07-30-gobase-design.md → internal/auth/middleware.go
- `HMAC-SHA256 JWT authentication` --references--> `OptionalAuth()`  [EXTRACTED]
  docs/superpowers/specs/2026-07-30-gobase-design.md → internal/auth/middleware.go
- `HMAC-SHA256 JWT authentication` --references--> `GetClaims()`  [EXTRACTED]
  docs/superpowers/specs/2026-07-30-gobase-design.md → internal/auth/middleware.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **GoBase batteries-included technology stack** — docs_superpowers_specs_2026_07_30_gobase_design_gobase, docs_superpowers_specs_2026_07_30_gobase_design_chi, docs_superpowers_specs_2026_07_30_gobase_design_cobra, docs_superpowers_specs_2026_07_30_gobase_design_sqlx, docs_superpowers_specs_2026_07_30_gobase_design_golang_migrate, docs_superpowers_specs_2026_07_30_gobase_design_hmac_jwt_auth, docs_superpowers_specs_2026_07_30_gobase_design_slog, docs_superpowers_specs_2026_07_30_gobase_design_go_playground_validator, docs_superpowers_specs_2026_07_30_gobase_design_config_management [EXTRACTED 1.00]
- **user domain reference vertical slice conventions** — docs_superpowers_specs_2026_07_30_gobase_design_user_domain_module, docs_superpowers_specs_2026_07_30_gobase_design_vertical_slicing, docs_superpowers_specs_2026_07_30_gobase_design_repository_pattern, docs_superpowers_specs_2026_07_30_gobase_design_context_propagation, docs_superpowers_specs_2026_07_30_gobase_design_sentinel_error_mapping [INFERRED 0.85]

## Communities (13 total, 2 thin omitted)

### Community 0 - "User Domain Module"
Cohesion: 0.11
Nodes (20): context.Context, time.Time, newStubRepo(), TestServiceCreate(), TestServiceCreateValidation(), TestServiceDelete(), TestServiceDeleteNotFound(), TestServiceGetByID() (+12 more)

### Community 1 - "JWT Auth"
Cohesion: 0.14
Nodes (23): Claims, contextKey, Testing strategy, github.com/golang-jwt/jwt/v5.RegisteredClaims, GenerateToken(), ParseToken(), TestGenerateAndParseToken(), TestGenerateTokenWithExtraClaims() (+15 more)

### Community 2 - "Design Doc & CLI Root"
Cohesion: 0.20
Nodes (20): Execute(), GoBase Template Implementation Plan, chi HTTP router, cobra CLI framework, Config management (env + .env), context.Context propagation convention, No global state / dependency injection, GoBase Go Project Template Design (+12 more)

### Community 3 - "Logging"
Cohesion: 0.19
Nodes (19): log/slog.Level, testing.T, New(), parseLevel(), TestLevelMapping(), TestNewDebugLevel(), TestNewInvalidLevel(), TestNewJSONFormat() (+11 more)

### Community 4 - "HTTP Middleware & Response"
Cohesion: 0.16
Nodes (13): HTTP server infrastructure, log/slog.Logger, net/http.Handler, RequestLogger(), TestRequestLogger(), TestWriteError(), TestWriteJSON(), TestWriteJSONUnsupportedType() (+5 more)

### Community 5 - "Database & Migrations"
Cohesion: 0.16
Nodes (12): runMigrate(), golang-migrate migrations, github.com/jmoiron/sqlx.DB, github.com/spf13/cobra.Command, Connect(), parseDriver(), TestConnectUnreachable(), TestParseDriver() (+4 more)

### Community 6 - "HTTP Handlers"
Cohesion: 0.34
Nodes (8): net/http.Request, net/http.ResponseWriter, CodeFrom(), WriteError(), WriteJSON(), chi.Router, RegisterRoutes(), Handler

### Community 7 - "Config"
Cohesion: 0.27
Nodes (11): time.Duration, getEnv(), getEnvDuration(), getEnvInt(), Config, Load(), logFormatDefault(), TestLoadDefaults() (+3 more)

### Community 8 - "Error Mapping"
Cohesion: 0.40
Nodes (4): Sentinel error mapping, CodeFrom(), TestCodeFrom(), TestCodeFromWrapped()

## Knowledge Gaps
- **3 isolated node(s):** `github.com/wesleysnt/gobase`, `contextKey`, `setup.sh script`
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `GoBase Go Project Template Design` connect `Design Doc & CLI Root` to `Error Mapping`, `JWT Auth`, `HTTP Middleware & Response`, `Database & Migrations`?**
  _High betweenness centrality (0.100) - this node is a cross-community bridge._
- **Why does `HMAC-SHA256 JWT authentication` connect `Design Doc & CLI Root` to `JWT Auth`, `HTTP Middleware & Response`?**
  _High betweenness centrality (0.099) - this node is a cross-community bridge._
- **Why does `GetClaims()` connect `JWT Auth` to `User Domain Module`, `Design Doc & CLI Root`, `HTTP Handlers`?**
  _High betweenness centrality (0.096) - this node is a cross-community bridge._
- **Are the 4 inferred relationships involving `Load()` (e.g. with `TestLoadDefaults()` and `TestLoadFromEnv()`) actually correct?**
  _`Load()` has 4 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/wesleysnt/gobase`, `contextKey`, `setup.sh script` to the rest of the system?**
  _3 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `User Domain Module` be split into smaller, more focused modules?**
  _Cohesion score 0.1141025641025641 - nodes in this community are weakly interconnected._
- **Should `JWT Auth` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._