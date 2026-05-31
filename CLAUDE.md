# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

QuestionHelper (题小助) is a full-stack question bank and exam management platform for education. It consists of three sub-projects: a Go backend API, a Vue 3 admin web panel, and a UniApp mobile app. Documentation in `docs/` is written in Chinese.

## Repository Structure

```
questionhelper-server/   # Go backend (Gin + GORM + MySQL + Redis)
questionhelper-admin/    # Vue 3 admin panel (Element Plus + TypeScript)
questionhelper-app/      # UniApp mobile app (Vue 3 + TypeScript)
docs/                    # Requirements and technical design docs (Chinese)
```

## Common Commands

### Backend (questionhelper-server/)

```bash
cd questionhelper-server
make run          # Run dev server (go run cmd/server/main.go)
make build        # Build binary to bin/server
make test         # Run all tests (go test ./... -v)
make fmt          # Format code
make lint         # Lint with golangci-lint
make swagger      # Regenerate Swagger docs
make migrate      # Run database migrations
```

API integration tests: `bash test_api.sh` (requires running server at localhost:8080).

### Admin Frontend (questionhelper-admin/)

Package manager: **pnpm** (enforced — do not use npm/yarn). Node requirement: `^20.19.0 || >=22.12.0`.

```bash
cd questionhelper-admin
pnpm dev              # Vite dev server
pnpm build            # Type check + production build
pnpm build-only       # Production build (skip type check)
pnpm type-check       # vue-tsc --noEmit
pnpm lint             # ESLint + Prettier + Stylelint
```

### Mobile App (questionhelper-app/)

Package manager: **npm**.

```bash
cd questionhelper-app
npm run dev:h5          # Dev server (H5/web)
npm run dev:mp-weixin   # Dev server (WeChat mini-program)
npm run build:h5        # Production build (H5)
npm run build:mp-weixin # Production build (WeChat)
```

### Docker (full stack)

```bash
docker-compose up -d    # Starts MySQL 8.0, Redis 7, Go server, Nginx admin panel
```

## Backend Architecture

Layered architecture under `questionhelper-server/internal/`:

```
router/       → Route definitions (Gin route groups, per-domain files)
controller/   → HTTP handlers (one package per domain)
service/      → Business logic layer
repository/   → Data access layer (GORM queries)
dto/          → Request/response structs with validation tags
model/        → GORM model definitions (database entities)
middleware/   → Auth, CORS, rate limiting, RBAC, operation logging, sensitive word filtering
ws/           → WebSocket hub/client/message (real-time communication)
job/          → Background jobs (exam submit, import/export, notifications)
task/         → Scheduled tasks (cleanup, reminders, expiry)
```

Shared packages under `questionhelper-server/pkg/`: cache, captcha, config, consts, database, encrypt, errors, excel, jwt, logger, mq, response, sensitive, upload, validator.

**Entry point**: `cmd/server/main.go` — loads config, initializes MySQL/Redis/logger/JWT, starts Gin server.

**Route tiers** (all under `/api/v1`):
- **Public** (`/api/v1/auth/...`) — login, register, captcha, token refresh
- **Authorized** (`/api/v1/...`) — regular user endpoints (JWT required)
- **Admin** (`/api/v1/admin/...`) — admin-only endpoints (JWT + admin role)

**Auth**: JWT access + refresh token pattern; token blacklisting via Redis.

**Config**: YAML file at `config/config.yaml`, loaded via `pkg/config`.

## Admin Frontend Architecture

Feature-based Vue 3 SPA using Element Plus. Key directories under `src/`:

- `api/` — API modules per domain (Axios instances)
- `views/` — Page components per domain
- `stores/` — Pinia stores (app, user, permission, settings, dict, tags-view, tenant)
- `router/` — Static + dynamic routes; **routes are fetched from the backend based on user permissions** and injected dynamically
- `components/` — Shared components including a CURD framework
- `composables/` — Vue composables (SSE, table selection, recent menus)
- `layouts/` — Layout variants (Left, Top, Mix, Base)

**Style conventions**: SCSS + UnoCSS; Prettier enforces 100-char width, 2-space indent, double quotes, semicolons, trailing commas. ESLint enforces PascalCase component names in templates, Vue block order (template → script → style).

**Commit convention**: Conventional Commits enforced via commitlint + Commitizen (cz-git). Use `pnpm commit` for interactive commits.

## Mobile App Architecture

UniApp (Vue 3 + TypeScript) targeting H5, WeChat Mini Program, and native iOS/Android. Page routing is defined in `pages.json` (40+ pages). Tab bar: Home, Question Bank, Practice, Class, My Profile.

Key directories: `pages/`, `components/`, `api/`, `store/` (Pinia), `hooks/`, `utils/`.

## Key Technical Details

- **Database**: MySQL 8.0 with GORM auto-migration
- **Cache**: Redis 7 for session/token blacklisting and general caching
- **File storage**: Local filesystem or MinIO (S3-compatible), configurable
- **Real-time**: WebSocket for notifications; SSE (Server-Sent Events) used in admin frontend
- **Excel**: Backend uses excelize for import/export; admin frontend uses ExcelJS
- **Rich text**: @wangeditor-next/editor in admin; custom RichText component in app
- **RBAC**: Role-based access control with menu/permission assignment; middleware checks permissions per route
- **Sensitive word filtering**: Built-in content moderation middleware
- **Background processing**: In-memory message queue (`pkg/mq`) with job workers and cron-like scheduled tasks

## Infrastructure

Default local development requires MySQL and Redis. The server config is at `questionhelper-server/config/config.yaml`. Default admin credentials: `admin / admin123`.
