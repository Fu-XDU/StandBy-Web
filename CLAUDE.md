# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

**Full build (backend + frontend):**
```bash
make build          # builds Go binary, Vue frontend, and syncs assets
make start          # build + run server at http://127.0.0.1:1423/v1/web/
```

**Development (hot reload):**
```bash
make dev            # runs Go server + Vite dev server concurrently
```

**Frontend only:**
```bash
cd web && npm run dev        # Vite dev server (proxies /v1 to :1423)
cd web && npm run build      # type-check + build to web/dist/
cd web && npm run type-check # TypeScript check only
```

**Backend only:**
```bash
make build-backend  # outputs to bin/server
```

**Docker:**
```bash
make docker         # builds Docker image
```

## Architecture

This is a full-stack web clock app — a Vue 3 SPA served by a Go/Gin backend with a REST API for multi-device config sync.

### Backend (Go)

- `server.go` — CLI entry point using `urfave/cli` with port/host flags from `mingfu_go_common`
- `routes/` — Gin route setup; `web.go` serves the SPA with `index.html` fallback; `remote.go` defines `/v1/remote/*` API routes; `cors.go` adds custom CORS middleware (includes PATCH)
- `controller/remote.go` — HTTP handlers: register device, sync config, get config, patch config
- `store/remote/store.go` — In-memory config store (map + mutex); restarts clear all state
- `store/remote/config_float.go` — Float page config schema, field validation, and merge logic (timestamp-based last-write-wins)

The backend stores config only in memory — no database. The server runs on port 1423 by default.

### Frontend (Vue 3 + TypeScript)

- `web/src/App.vue` — Root component; selects clock style, applies anti-burn transforms (every 30s), manages dynamic routing between clock styles
- `web/src/components/Float.vue` — Main large-digit clock display (37vw font, 14 color presets)
- `web/src/components/Numerical.vue` — Alternate numerical clock style
- `web/src/components/ClockLayout.vue` — Control panel: clock style selector, color picker, brightness slider, night mode schedule
- `web/src/composables/useFloatConfig.ts` — Central config state (night mode, invisible mode, brightness, per-style colors)
- `web/src/composables/usePageRemoteSync.ts` — Syncs config with the server on load and every minute
- `web/src/composables/useLocalStorageSync.ts` — Persists reactive refs to `localStorage`
- `web/src/composables/useRemoteDevice.ts` — Generates/stores UUID `deviceId` in `localStorage` as `standby_device_id`
- `web/src/config/floatConfig.ts` — Float config type definitions, defaults, validation, and equality checks

### Key Data Flow

1. On page load, the frontend registers its `deviceId` via `POST /v1/remote/register`
2. Config syncs to/from server every minute via `POST /v1/remote/sync` (full config with `updatedAt` timestamps; server returns merged result)
3. External clients can push partial updates via `PATCH /v1/remote/config`; the browser picks them up on the next sync cycle
4. Config is also persisted in `localStorage` for offline/reload resilience

### Asset Serving

- Dev: Vite dev server proxies `/v1/*` to the Go backend at `:1423`
- Production: `make build` builds the frontend to `web/dist/` and `make sync-frontend` copies it to `assets/web/v1/`; the Go server serves this directory and falls back to `index.html` for SPA routing
- The Vite `base: './'` setting allows hosting under any URL prefix (e.g. `/standby/`, `/clock/`)

### Remote Config API

| Method | Path | Description |
|--------|------|-------------|
| POST | `/v1/remote/register` | Register deviceId |
| POST | `/v1/remote/sync` | Full config sync with updatedAt timestamps |
| GET | `/v1/remote/config?deviceId=&pageId=float` | Read config |
| PATCH | `/v1/remote/config` | Partial field update (`pageId` defaults to `float`) |

Patchable fields: `brightness`, `selectedColorIndex`, `autoNightMode`, `nightModeRange`, `autoInvisible`, `invisibleRange`, `invisibleDayEnable`, `invisibleDay`.
