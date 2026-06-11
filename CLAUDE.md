# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

MiniBank — didactic-but-realistic banking core in Go + Angular, built incrementally across 9 phases (F0→F8). Currently early F0/F1: backend is a skeleton API with stub domain packages; frontend is the default Angular SSR scaffold. Use `/home/savedra/documents/pessoal/mentorship-juansavedra` as **read-only** reference for roadmap, phase Definition of Done, and learning goals — never edit files there.

Agent role: act as developer *and* teacher — implement the change, explain key choices, highlight concepts to learn.

## Commands

Run from the relevant module directory. Shell tooling here prefixes commands with `rtk` (e.g. `rtk go test ./...`).

```bash
# Backend (cd backend)
go run ./cmd/api        # start API on :8080
go test ./...           # run tests
go test -race ./...     # REQUIRED for ledger/transfer/concurrency code
gofmt -w .              # format

# Frontend (cd frontend)
npm install
npm start               # ng serve dev server
npm run build           # production build with SSR output
npm test                # unit tests (ng test, Vitest-compatible)
npm run serve:ssr:frontend   # run built SSR server
```

Single Go test: `go test ./cmd/api -run TestHealthHandler`. Single frontend test: `ng test` filters via Vitest.

## Architecture

Two independent modules at repo root, no shared build:

- **`backend/`** — Go 1.25 module (`module backend`). Entry point `cmd/api/main.go` (`net/http`, `/` and `/health` handlers, tests in `main_test.go`). Domain logic lives under `internal/<domain>/` — currently stub packages `account`, `auth`, `transaction`, each exporting only a `DomainName` const. The README targets a richer set (`customer`, `account`, `ledger`, `transfer`, `statement`, `auth`, `notification`, `backoffice`); keep domain code, HTTP handlers, persistence, and migrations in separate packages as these grow.
- **`frontend/`** — Angular 21 standalone + SSR. App code in `src/app/` (`app.ts`, `app.routes.ts`, `app.config.ts`, server variants `*.server.ts`). `app-` selector prefix, SCSS, `*.spec.ts` beside covered code.

Planned stack (mostly not yet present): backend `pgx`+`sqlc`, `golang-migrate`, JWT+bcrypt+TOTP, `log/slog`, `testcontainers-go`; frontend Signals + Reactive Forms + JWT interceptor; infra on AWS via Terraform + GitHub Actions.

## Non-negotiable domain rules

These are core correctness invariants — apply them whenever touching money/ledger code:

- **Money is `int64` cents.** Never floating point for financial amounts.
- **Double-entry ledger.** Every movement is a debit/credit pair that sums to zero.
- **Balance is derived** from ledger entries, never a free mutable field; update in the same ACID transaction, guard with DB `CHECK` + locking.
- **Idempotency** via persisted `Idempotency-Key` — retries must not duplicate entries.
- **Concurrency safety** — DB lock is source of truth; run suites with `-race`. Use table-driven tests for double-entry invariants, idempotency, concurrency.
- **Security shift-left** — no secrets/CPF/PII in code; correlation IDs and JSON logs from the first endpoint; least privilege.

## Conventions

- Commits: Conventional Commits **in Portuguese** (`feat:`, `fix:`, `docs:`, `test:`, `refactor:`).
- Go: `gofmt`, short lowercase package names.
- TypeScript: strict mode, Angular standalone patterns.
- PRs: link the roadmap phase; call out security, schema, or ledger impacts; screenshots for visible frontend changes.
