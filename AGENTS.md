# AGENTS.md

This file is guidance for autonomous coding agents working in this repository.
It captures how to build, lint, test, and follow existing code conventions.

## Repository Overview

- Monorepo with a Go backend at repo root and a Vite + React + TypeScript frontend in `web/`.
- Package manager for JS/TS work is `pnpm` (workspace includes only `web`).
- Main backend entrypoint: `cmd/main.go`.
- API routes are wired in `internal/router/router.go`.
- Most backend handlers live in `internal/handlers/handlers.go`.

## Toolchain and Runtime

- Go version from `go.mod`: `go 1.24.0`.
- Node used by Docker build: `node:22-alpine`.
- Frontend build stack: TypeScript + Vite + ESLint.
- Backend database: SQLite via GORM.

## Install and Setup Commands

- Install root workspace dependencies (for orchestration scripts):
  - `pnpm install`
- Install frontend dependencies directly:
  - `pnpm --filter web install`
- Download Go modules:
  - `go mod download`

## Build Commands

- Build backend binary:
  - `go build -o server ./cmd/main.go`
- Build frontend production assets:
  - `pnpm --filter web build`
- Build both (sequential):
  - `go build -o server ./cmd/main.go && pnpm --filter web build`

## Dev Commands

- Run backend only:
  - `pnpm dev:api`
  - equivalent: `go run ./cmd/main.go`
- Run frontend only:
  - `pnpm dev:web`
  - equivalent: `pnpm --filter web dev`
- Run backend + frontend together:
  - `pnpm dev`

## Lint Commands

- Lint everything configured at root:
  - `pnpm lint`
- Lint backend only:
  - `pnpm lint:api`
  - equivalent: `go vet ./...`
- Lint frontend only:
  - `pnpm lint:web`
  - equivalent: `pnpm --filter web lint`

## Test Commands

Current state:

- No Go test files (`*_test.go`) are currently committed.
- No frontend test runner (Vitest/Jest) is currently configured.

Backend test commands to use when adding tests:

- Run all Go tests:
  - `go test ./...`
- Run tests for one package:
  - `go test ./internal/handlers`
- Run one test function (single test):
  - `go test ./internal/handlers -run '^TestPhoneLogin$' -v`
- Run subtests by name:
  - `go test ./internal/handlers -run 'TestPhoneLogin/invalid_phone' -v`

Frontend testing note:

- There is no `test` script in `web/package.json` yet.
- If a test runner is added later, add commands here and prefer package-filtered runs:
  - `pnpm --filter web test`
  - single test would usually be something like `pnpm --filter web test -- <pattern>`

## Docker Commands

- Build and run containerized stack:
  - `docker compose up --build -d`
- Stop stack:
  - `docker compose down`
- Stop and remove volumes (destructive to local DB/uploads):
  - `docker compose down -v`

## Architecture and File Layout

- `cmd/main.go`: app boot, config load, DB init, router startup.
- `internal/config`: env config, DB init, migrations, seed data.
- `internal/models`: GORM models and DB field mapping.
- `internal/middleware`: auth, CORS, request logging middleware.
- `internal/handlers`: HTTP handlers and response shaping.
- `internal/router`: all route registration.
- `web/src/api`: axios client and typed API service wrappers.
- `web/src/store`: Zustand auth store.

## Code Style: Go

- Always run `gofmt`/`goimports` style formatting (Go defaults, tabs, grouped imports).
- Keep package names short and lowercase (`config`, `handlers`, `router`).
- Exported symbols use PascalCase; internal helpers use camelCase.
- Import order should be: stdlib, internal module imports, third-party imports.
- Prefer explicit error checks with early returns.
- Wrap infrastructure errors with context using `%w` when returning errors.
- In HTTP handlers, return JSON responses consistently with existing `code`/`msg` shape.
- Validate request params before DB operations (current codebase pattern).
- Use `strconv.Atoi` + positivity checks for numeric path params.
- Use `gorm.Expr` for atomic counter increments.
- Keep route registration centralized in `internal/router/router.go`.
- Keep response mapping helper functions near handlers (existing pattern).

## Code Style: TypeScript + React

- TypeScript is strict (`strict: true`, `noUnusedLocals`, `noUnusedParameters`).
- Use functional React components and hooks.
- Prefer small local helper functions for form parsing and coercion.
- Use single quotes and no semicolons (match current files).
- Keep imports at top; group external imports before local imports.
- Use `type` aliases for local state and payload structures.
- Use `unknown` at boundaries, then narrow/cast near rendering/usage.
- Keep API calls in `web/src/api/service.ts`; do not scatter raw axios calls.
- Use `apiClient` interceptor for auth token propagation.
- Keep UI state updates with optimistic clarity: loading, result, and error tracked explicitly.
- Reuse shared helpers like `resolveAssetUrl` for asset URLs.

## Naming Conventions

- Go structs and fields mirror API payload naming via JSON tags.
- DB columns are snake_case in GORM tags.
- TS variable and function names are camelCase.
- React component names are PascalCase.
- Avoid single-letter names except in tight loop scopes.

## Error Handling Conventions

- Backend handlers:
  - Return HTTP 200 for many business errors when matching existing API contract.
  - Use `code` and `msg` fields for client-visible status.
  - Return 404-style `code` values where resource absence is semantically important.
- Middleware/auth errors may return HTTP 401 directly (existing behavior).
- Frontend:
  - Check `axios.isAxiosError` before reading response payload.
  - Prefer human-readable fallbacks (`msg` -> axios message -> generic fallback).

## Config and Secrets

- Environment variables used by backend:
  - `SERVER_PORT` (default `8080`)
  - `DB_PATH` (default `./data.db`)
  - `JWT_SECRET` (default in code; override in non-local environments)
- Do not commit real secrets.
- Treat default JWT secret as development-only.

## Agent Execution Guidelines

- Before finalizing code changes, run relevant lint/build commands for touched areas.
- If you add tests, run the narrowest relevant test first, then broader suites.
- Keep changes scoped; avoid unrelated refactors in this repo.
- Preserve API response contracts and route paths unless explicitly requested.

## Cursor/Copilot Rules Check

- Checked for `.cursorrules`: not found.
- Checked for `.cursor/rules/`: not found.
- Checked for `.github/copilot-instructions.md`: not found.
- If these files are added later, update this AGENTS.md and treat those rules as higher-priority constraints.
