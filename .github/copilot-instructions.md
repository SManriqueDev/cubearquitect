# Copilot instructions for `cubearquitect`

## Build, test, and lint commands

### Repository-level helpers
- Install deps for both apps: `make install`
- Run backend + frontend in parallel: `make dev`
- Run backend only: `make dev-backend`
- Run frontend only: `make dev-frontend`

### Backend (`backend/`)
- Start API locally: `go run cmd/api/main.go`
- Module maintenance: `go mod tidy`
- Full tests: `go test ./...`
- Single test (by name): `go test ./... -run TestName`
- Single package tests: `go test ./cmd/api -run TestName`

### Frontend (`frontend/`)
- Install deps: `npm install`
- Dev server: `npm run dev`
- Build: `npm run build`
- Lint: `npm run lint`
- Test commands are not configured in `package.json` yet.

## High-level architecture

- This repository is split into two apps coordinated by the root `Makefile`:
  - `backend/`: Fiber API service in Go
  - `frontend/`: React + Vite TypeScript app
- Backend request flow:
  - Entry point: `backend/cmd/api/main.go`
  - Config loading: `backend/internal/config/config.go`
  - External API integration: `backend/internal/cubepath/client.go`
  - The API initializes one CubePath client at startup and reuses it in handlers.
- Current backend HTTP surface:
  - `GET /health` returns `{"status":"alive"}`
  - `GET /api/projects` proxies to CubePath `/projects/` via the internal client
- Frontend entrypoint is `frontend/src/main.tsx` mounting `App` under `StrictMode`.
- Current frontend app (`frontend/src/App.tsx`) is a full-viewport React Flow canvas with local initial nodes; no frontend-side data fetching is wired yet.

## Key repository conventions

- Backend configuration conventions:
  - `CUBE_API_TOKEN` is required at startup (app exits if missing).
  - `CUBE_API_URL` defaults to `https://api.cubepath.com` when unset.
  - `PORT` defaults to `8080` when unset.
- Keep backend integration logic in `internal/cubepath/client.go`; HTTP handlers should call client methods instead of duplicating request/header/timeout logic.
- Backend middleware baseline is `logger` + `cors` in `cmd/api/main.go`; keep this order unless there is a clear reason to change it.
- Frontend TypeScript is strict (`strict`, `noUnusedLocals`, `noUnusedParameters`, `noUncheckedSideEffectImports`); keep new code compatible with these compiler checks.
- Frontend linting uses `eslint.config.js` flat config with `typescript-eslint`, `react-hooks`, and `react-refresh`; align new files with this ruleset.
- Prefer root Make targets (`make install`, `make dev`, `make dev-backend`, `make dev-frontend`) for local workflows so backend/frontend commands stay consistent.
