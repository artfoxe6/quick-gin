# Quick-Gin

Quick-Gin is a Gin-based Go scaffold that helps you bootstrap production-ready RESTful APIs in minutes.

## Highlights
- Layered architecture (Model → Repository → Service → Handler) with dependency injection
- JWT authentication, Redis/memory caching, and ready-to-use middleware (CORS, logging, recovery, signature)
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

## Project layout
```
quick-gin/
├── cmd/                    # Application entry points
│   └── app/               # Main application
│       └── main.go
├── internal/              # Private application code
│   ├── app/              # Application layer
│   │   ├── core/         # Core components (config, middleware, models, etc.)
│   │   │   ├── apperr/   # Error handling
│   │   │   ├── config/   # Configuration management
│   │   │   ├── middleware/ # HTTP middleware (auth, cors, recovery, etc.)
│   │   │   ├── model/    # Base models
│   │   │   ├── repository/ # Base repository with query builder
│   │   │   ├── request/  # Request DTOs
│   │   │   └── router/   # Route configuration
│   │   └── user/         # User module example
│   │       ├── dto/      # Data transfer objects
│   │       ├── handler/  # HTTP handlers
│   │       ├── model/    # Data models
│   │       ├── repo/     # Repository layer
│   │       └── service/  # Business logic layer
│   ├── pkg/              # Shared packages
│   │   ├── cache/        # Cache abstraction (memory/redis)
│   │   ├── cronjob/      # Scheduled task management
│   │   ├── db/           # Database connection
│   │   ├── kit/          # Utility functions
│   │   ├── mailer/       # Email service
│   │   ├── oss/          # Object storage service
│   │   └── token/        # JWT token handling
│   └── scaffold/         # Code generation utilities
├── data/                 # Data files
│   └── app.db           # SQLite database
├── config.ini            # Configuration file
├── main.go              # Alternative application entry
├── Makefile             # Build commands
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
└── README.md            # Project documentation
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
