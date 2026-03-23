# b022hub Development Guide

This document is safe to publish and is intended for contributors working on the open-source repository.

## Project Basics

| Item | Details |
|------|---------|
| Upstream repository | `Wei-Shaw/sub2api` |
| This repository | `B022MC/b022hub` |
| Backend stack | Go, Gin, Ent, PostgreSQL, Redis |
| Frontend stack | Vue 3, Vite, TypeScript, Pinia |
| Frontend package manager | `pnpm` |

## Prerequisites

- Go `1.25.7`
- Node.js and `pnpm`
- PostgreSQL `15+`
- Redis `7+`
- Docker / Docker Compose for container-based local runs

## Local Development

### Backend

```bash
cd backend
go test -tags=unit ./...
go test -tags=integration ./...
golangci-lint run ./...
```

### Frontend

```bash
cd frontend
pnpm install --frozen-lockfile
pnpm build
pnpm test:run
```

### Docker-Based Local Run

```bash
cd deploy
cp .env.example .env
docker compose -f docker-compose.local.yml up -d
```

## Contributor Workflow

- Keep changes focused and submit small, reviewable pull requests.
- Update documentation when behavior, configuration, or deployment steps change.
- Never commit secrets, private keys, production endpoints, or populated `.env` files.
- When `frontend/package.json` changes, commit the updated `frontend/pnpm-lock.yaml`.
- When `backend/ent/schema/*.go` changes, regenerate the Ent artifacts before opening a pull request.

## Common Regeneration Steps

```bash
cd backend
go generate ./ent
```

## Common Pitfalls

### `pnpm-lock.yaml` out of sync

If `frontend/package.json` changed, run:

```bash
cd frontend
pnpm install
```

### Interface changed but tests no longer compile

Search the backend test doubles and update every stub or mock that implements the modified interface.

### Ent schema changes do not take effect

Regenerate Ent code after editing schema files:

```bash
cd backend
go generate ./ent
```

### `localhost` behaves differently from `127.0.0.1`

If your local PostgreSQL or Redis client has IPv4/IPv6 resolution issues, try `127.0.0.1` explicitly during troubleshooting.

## CI Overview

The repository uses GitHub Actions for backend CI, release builds, and security checks. Before opening a pull request, run the relevant backend and frontend checks locally.
