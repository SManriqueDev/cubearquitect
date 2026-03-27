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
- Hot reload during development: `make dev-backend-hot` (uses `air` for file watching)

### Frontend (`frontend/`)
- Install deps: `npm install`
- Dev server: `npm run dev`
- Build: `npm run build`
- Lint: `npm run lint`
- Preview build locally: `npm run preview`
- Type checking included in build pipeline: `tsc -b && vite build`

## High-level architecture

- This repository is split into two apps coordinated by the root `Makefile`:
  - `backend/`: Fiber API service in Go
  - `frontend/`: React + Vite TypeScript app
- Backend request flow:
  - Entry point: `backend/cmd/api/main.go` loads config and initializes Fiber app
  - Config loading: `backend/internal/config/config.go` (loads env vars, provides defaults)
  - External API integration: `backend/internal/cubepath/client.go` provides generic HTTP helpers
  - CubePath client initialized per-request in `CubeTokenMiddleware` (not singleton)
  - Server/app wiring: `backend/internal/app/server.go` initializes services; `routes.go` registers handlers
  - HTTP handlers are structs under `backend/internal/handler`, wired via constructors and methods
  - Business logic lives in `backend/internal/service`, which parses CubePath responses and owns API URLs
  - `backend/internal/orchestrator` manages deployment orchestration and state tracking
- Current backend HTTP surface:
  - `GET /health` returns `{"status":"alive"}` (no auth required)
  - Protected routes (require `X-Cube-Token` header):
    - `GET /api/projects` - list CubePath projects
    - `GET /api/ssh-keys` - list SSH keys
    - `GET /api/pricing` - get pricing information
    - `POST /api/deploy` - initiate deployment (validates payload with Zod schema)
    - `GET /api/deployments` - list active deployments
    - `GET /api/deployments/:deployment_id` - get deployment status
    - `GET /api/deployments/:deployment_id/events` - WebSocket stream for real-time deployment events
- Frontend entrypoint is `frontend/src/main.tsx` mounting `App` under `StrictMode`.
- Current frontend app (`frontend/src/App.tsx`) is a full-viewport React Flow canvas with local initial nodes; no frontend-side data fetching is wired yet.

## Key repository conventions

- Backend configuration conventions:
  - `CUBE_API_URL` (optional) defaults to `https://api.cubepath.com`; can be overridden per-request via `X-Cube-API-URL` header
  - `PORT` (optional) defaults to `8080`
  - `CUBE_API_TOKEN` is passed at request-time via `X-Cube-Token` header (or `token` query param for WebSocket)
- Keep backend integration logic in `internal/cubepath/client.go`; HTTP handlers should call client methods instead of duplicating request/header/timeout logic.
- Backend middleware baseline is `logger` + `cors` in `cmd/api/main.go`; protected routes add `CubeTokenMiddleware` via a group.
- Frontend TypeScript is strict (`strict`, `noUnusedLocals`, `noUnusedParameters`, `noUncheckedSideEffectImports`); keep new code compatible with these compiler checks.
- Frontend linting uses `eslint.config.js` flat config with `typescript-eslint`, `react-hooks`, and `react-refresh`; align new files with this ruleset.
- Frontend validation: Use Zod schemas (e.g., `deployPayloadSchema` in `frontend/src/services/schemas/flow.ts`) and validate before POSTing to backend.
- Frontend state management: Zustand stores under `frontend/src/stores/` (e.g., `accountStore` for user config, `flowStore` for UI state); React Query for server state in hooks.
- Prefer root Make targets (`make install`, `make dev`, `make dev-backend`, `make dev-frontend`) for local workflows so backend/frontend commands stay consistent.
- Backend handlers follow a struct pattern (`type XHandler struct { svc *service.YService }`) with constructors (`NewXHandler`) whose methods are registered in `internal/app/routes.go`. Inject services via constructors.
- Services under `internal/service` encapsulate business rules (e.g., `ProjectsService.List`, `PricingService.GetPricing`) and own the CubePath URLs/parsing; `internal/cubepath/client.go` only yields generic `Get/Post/...` helpers.
- `internal/orchestrator` manages deployment orchestration (tracks state, emits events); keep it free of HTTP handler code and let services call it when needed.
- Frontend hooks (e.g., `useDeploy`, `useDeploymentEvents` in `frontend/src/hooks/`) encapsulate API calls and mutations; use `@tanstack/react-query` for server state synchronization.

## MCP Servers & AI Configuration

- OpenCode configuration in `opencode.json` exposes a shadcn MCP server for component management and UI generation
- Frontend uses shadcn/ui components (via `components.json` and Tailwind CSS) — when adding UI components, use shadcn patterns and consult the component registry
