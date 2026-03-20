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

- This repo is a two-app setup (no monorepo tool): a Go API in `backend/` and a React/Vite app in `frontend/`, coordinated from the root `Makefile`.
- Backend entrypoint is `backend/cmd/api/main.go` using Fiber v2:
  - Loads env via `godotenv.Load()`
  - Applies global CORS middleware
  - Exposes `GET /health` returning JSON `{status, api}`
  - Binds to `PORT` env var, defaulting to `8080`
- Frontend entrypoint is `frontend/src/main.tsx` mounting `App` under `StrictMode`.
- UI root (`frontend/src/App.tsx`) is currently a full-viewport React Flow canvas (`@xyflow/react`) with initial graph nodes and built-in controls/background.
- Dev runtime assumptions:
  - Frontend Vite server runs on `5173` and `host: true`
  - Backend defaults to `8080`
  - Vite uses polling file watch (`usePolling: true`) for container/volume-friendly reload behavior

## Key repository conventions

- Keep backend executable code rooted at `cmd/api/main.go` unless/until package structure is expanded.
- Keep API health contract stable (`/health` with `status` and `api`) because it is the only explicit backend contract currently present.
- Use env-driven backend port (`PORT`) and preserve default fallback behavior.
- Frontend TypeScript is strict (`strict`, `noUnusedLocals`, `noUnusedParameters`, `noUncheckedSideEffectImports`); keep new code compatible with these compiler checks.
- Frontend linting uses `eslint.config.js` flat config with `typescript-eslint`, `react-hooks`, and `react-refresh`; align new files with this ruleset.
- Root workflow expectation is Makefile-first for common dev tasks (`make install`, `make dev`) instead of custom ad-hoc scripts.
