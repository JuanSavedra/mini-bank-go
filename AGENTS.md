# Repository Guidelines

## Project Structure & Module Organization

MiniBank has two top-level modules:

- `backend/`: Go module (`go.mod`) for the banking API and domain services. Keep domain code, HTTP handlers, persistence, and migrations separated by package.
- `frontend/`: Angular 21 standalone app with SSR. App code lives in `frontend/src/app/`, global styles in `frontend/src/styles.scss`, assets in `frontend/public/`, and config in `angular.json`.
- `README.md`: product goals, roadmap, stack, and principles.

## Reference Context & Agent Role

Use `/home/savedra/documents/pessoal/mentorship-juansavedra` exclusively as a read-only source for broader context, roadmap details, and learning goals. Do not edit files there.

Act as both developer and teacher: implement the change, explain key choices, and highlight concepts to learn.

## Build, Test, and Development Commands

Run commands from the relevant module directory:

```bash
cd frontend && npm install      # install Angular dependencies
cd frontend && npm start        # run local Angular dev server
cd frontend && npm run build    # production build with SSR output
cd frontend && npm test         # run Angular unit tests via Vitest
cd backend && go test ./...     # run Go tests once packages exist
```

When using Codex shell tooling here, prefix commands with `rtk`, for example `rtk npm test`.

## Coding Style & Naming Conventions

Use Go 1.25 conventions in `backend/`: `gofmt` all files, keep package names short and lowercase, and prefer explicit `int64` cents for money values. Do not use floating point for financial amounts.

Use strict TypeScript in `frontend/`. Follow Angular standalone patterns, SCSS styles, `app-` selector prefix, and `*.spec.ts` tests beside covered code. Prefer names such as `transfer.service.ts`.

## Testing Guidelines

Frontend tests use Angular's unit-test builder with Vitest-compatible tooling. Add `*.spec.ts` tests for rendering, form validation, routing, and service edge cases.

Backend tests should use Go's `testing` package. For ledger and transfer code, add table-driven tests for double-entry invariants, idempotency, and concurrency. Run race checks with `go test -race ./...`.

## Commit & Pull Request Guidelines

Git history uses Conventional Commit style in Portuguese. Keep prefixes like `feat:`, `fix:`, `docs:`, `test:`, and `refactor:`.

Pull requests should include a concise description, test results, linked issue or roadmap phase, and screenshots for visible frontend changes. Call out security, schema, or ledger impacts.

## Security & Configuration

Never commit secrets, CPF/PII samples, private keys, or AWS credentials. Keep configuration in environment variables or managed secret stores. Preserve README principles: derived balances, idempotency keys, correlation IDs, and least privilege.
