# Quick-Gin

Quick-Gin is a Gin-based Go scaffold that helps you bootstrap production-ready RESTful APIs in minutes.

## Highlights
- Layered architecture (Model → Repository → Service → Handler) with dependency injection
- JWT authentication, Redis/memory caching, and ready-to-use middleware (CORS, logging, recovery, signature)
- CRUD generator that creates model, repository, service, handler, and route boilerplate
- Optional integrations for Aliyun OSS, SendGrid mailer, and cron jobs
- Works with MySQL or SQLite out of the box

## Getting Started

### Scaffold a new project
```bash
go run github.com/artfoxe6/quick-gin@latest myapp --module github.com/you/myapp
cd myapp
```
- Omit `--module` to default to the folder name.
- Add `--force` to overwrite an existing directory.
- The command copies the template and runs `go mod tidy` automatically.

## Run the app
1. Duplicate `config.ini` and adjust the settings you need (database, cache, JWT, etc.).
2. Start the server:
   ```bash
   go run cmd/app/main.go -config config.ini
   # or
   make run
   ```
3. The API listens on `http://localhost:8080` by default.

## CRUD generator
```bash
go run cmd/gen/gen.go --module Product
```
This produces matching files under `internal/app/{models,repositories,services,handlers}` and injects the routes automatically.

## Project layout
```
quick-gin/
├── cmd/            # Entry points (app, code generator, project cloner)
├── internal/app/   # Config, handlers, middleware, models, repositories, services
├── internal/pkg/   # Shared components (cache, db, mailer, oss, token, etc.)
├── config.ini      # Configuration template
├── Makefile        # Helper commands
└── README.md
```

## Tooling
- Web: Gin v1.6.0
- ORM: GORM v1.25.1
- Auth: JWT v5.0.0
- Cache: Redis v2.0.0 or in-memory
- Config: Go-ini v1.55.0

## Contributing
Issues and pull requests are welcome.

## License
MIT
