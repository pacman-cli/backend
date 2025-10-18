## Basic New Project — Go Backend + Next.js Frontend

![Go](https://img.shields.io/badge/Go-1.25.2-00ADD8?logo=go&logoColor=white)
![Next.js](https://img.shields.io/badge/Next.js-15.5.5-000000?logo=nextdotjs)
![MySQL](https://img.shields.io/badge/MySQL-8+-4479A1?logo=mysql&logoColor=white)

A small monorepo for a blog-style API built with Go (no framework) and a modern Next.js frontend. The backend uses MySQL with JSON columns; the frontend consumes the API and renders a minimal blog UI.

### Repo Structure

```
.
├── cmd/
│   ├── migrate/          # migration runner (applies SQL files under migrations/)
│   └── server/           # backend HTTP server entrypoint
├── internal/
│   ├── config/           # env config loader
│   ├── db/               # database connection (database/sql + mysql driver)
│   ├── models/           # domain models
│   ├── repositories/     # data access layer (CRUD)
│   ├── services/         # business logic & validation
│   └── transport/http/   # handlers + middleware (logging, recovery, CORS)
├── migrations/           # SQL migrations (ordered by filename)
└── frontend/             # Next.js app consuming the API
```

---

## Backend (Go + MySQL)

### Requirements
- Go 1.25+
- MySQL 8+

### Environment Variables

| Name | Default | Description |
|------|---------|-------------|
| `APP_PORT` | `8080` | HTTP port for the API |
| `DB_HOST` | `localhost` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASS` | empty | MySQL password |
| `DB_NAME` | `blogdb` | Database name |
| `RUN_MIGRATIONS` | `false` | If `true`, attempt to run migrations on startup |

> Tip: Prefer `.env` or your process manager to inject secrets; avoid committing real credentials.

### Quick Start (Local)

1) Start MySQL and create a database (default `blogdb`).
2) Export env vars as needed (see table above).
3) Apply migrations:
```bash
go run ./cmd/migrate
```
4) Run the API server:
```bash
go run ./cmd/server
```

The server will listen on `http://localhost:8080` by default.

### API Overview
- `GET /health` — health check
- `POST /posts` — create a post
- `GET /posts` — list posts
- `GET /posts/{id}` — get a post by id
- `PUT /posts/{id}` — update a post
- `DELETE /posts/{id}` — delete a post

Create example:
```bash
curl -sS -X POST http://localhost:8080/posts \
  -H 'Content-Type: application/json' \
  -d '{
        "title": "First Post",
        "content": "Hello, world!",
        "tags": ["intro", "hello"],
        "metadata": {"featured": true}
      }'
```

### Development
```bash
go mod tidy         # fetch deps
go build ./...      # build all
go test ./...       # run tests (if/when added)
```

---

## Frontend (Next.js)

The frontend lives in `frontend/`. See an expanded guide in `frontend/README.md`.

### Minimal Setup
Create `frontend/.env.local`:
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

Then:
```bash
cd frontend
npm install
npm run dev
```

Open `http://localhost:3000`.

---

## Migrations

Migrations are plain SQL files in `migrations/`. Filenames are applied in lexicographic order.

Run them locally with:
```bash
go run ./cmd/migrate
```

> You can re-run safely if your SQL is idempotent. Prefer using `IF NOT EXISTS` and `WHERE NOT EXISTS` guards in seed scripts.

---

## Troubleshooting
- Cannot connect to DB: verify `DB_HOST`, `DB_PORT`, credentials, and that MySQL is reachable. Try `mysql -h 127.0.0.1 -P 3306 -u root -p`.
- CORS issues from the frontend: the backend applies a permissive CORS middleware for local dev. Adjust in `internal/transport/http/middleware` for production.
- 500 errors: check backend logs in the server process; ensure migrations have run.

---

## Contributing
1) Fork and create a feature branch
2) Make changes with clear commits
3) Open a PR describing rationale and testing notes

---

## Acknowledgements
- Go standard library HTTP + `database/sql`
- `github.com/go-sql-driver/mysql`
- Next.js 15 + React 19
# backend
