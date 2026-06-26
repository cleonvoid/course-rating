# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run locally (requires Docker)
docker-compose up --build

# Build binary
go build -o server .

# Run after building (requires DATABASE_URL env var)
DATABASE_URL=postgres://postgres:postgres@localhost:5432/course_rating?sslmode=disable ./server

# Regenerate database code after changing SQL queries or migrations
sqlc generate

# Verify code compiles
go build ./...
```

There are no automated tests.

## Architecture

This is a Go web app using server-rendered HTML with HTMX for partial updates. The stack is `net/http` + `html/template` + PostgreSQL via pgx.

### Request flow

`main.go` registers all routes and wires up handlers. Each handler struct holds `*db.Queries`, `*template.Template`, and `sessionSecret`. Auth is a plain HMAC-signed cookie (see `internal/session/session.go`) — no middleware, no context values. Every handler that needs the current user calls `session.Get(r, secret)` directly, then looks up the user.

### Database layer

SQL queries live in `sql/queries/*.sql` with sqlc annotations (`:one`, `:many`, `:exec`). Running `sqlc generate` produces `internal/db/*.sql.go`. **Never edit `internal/db/` by hand** — it is fully generated. The schema source of truth is `migrations/*.up.sql`; sqlc reads all files in `migrations/` to build the schema.

Migrations run automatically on startup via `golang-migrate` using an embedded FS (`//go:embed migrations`).

### Templates

All templates are embedded at build time (`//go:embed internal/templates`). They are parsed once at startup in `parseTemplates()` in `main.go` — **new template files must be added to that function's `ParseFS` call** or they won't be available.

`base.html` defines the `{{template "head" .}}` and `{{template "nav" .}}` blocks reused by every page. It also holds **all global CSS** — there is no separate stylesheet.

HTMX is used for the star-rating widget and enroll button: those handlers return HTML fragments (partials), not full pages. Partials live in `internal/templates/partials/`.

### Adding a new feature

The typical pattern:
1. Add a `migrations/XXXXXX_<name>.up.sql` (and matching `.down.sql`)
2. Add queries to `sql/queries/<name>.sql`
3. Run `sqlc generate`
4. Write a handler in `internal/handler/<name>.go`
5. Register routes in `main.go`
6. Add template file and register it in `parseTemplates()`
7. Add CSS to the `<style>` block in `base.html`

### Seed data

The `.up.sql` migration files contain `INSERT` statements with demo users and content. Demo passwords are in a comment in `000003_create_users.up.sql` (pattern: `<name>123`).
