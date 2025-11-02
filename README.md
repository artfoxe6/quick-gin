# Quick-Gin

🚀 A Go web application scaffold based on Gin framework, designed for rapid development of RESTful API services.

[![Go](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-1.6.0-red.svg)](https://gin-gonic.com)
[![GORM](https://img.shields.io/badge/GORM-1.25.1-green.svg)](https://gorm.io)

## ✨ Features

- 🏗️ **Layered Architecture** - Model → Repository → Service → Handler
- 🔄 **Multi-Database Support** - MySQL / SQLite easy switching
- 💾 **Multi-Cache Support** - Redis / Memory cache switching
- 🔐 **Complete Authentication** - JWT token authentication with role-based access control
- 📝 **CRUD Code Generation** - One-click generation of complete CRUD code
- 🛠️ **Rich Middleware** - CORS, authentication, logging, signature verification
- 📦 **External Services** - OSS object storage, SendGrid email
- ⚡ **High Performance** - Connection pool, caching strategy optimization

## 📁 Project Structure

```
quick-gin/
├── cmd/                    # Application entry points
│   ├── app/               # Main application service
│   ├── gen/               # CRUD code generation tool
│   └── clone/             # Project cloning tool
├── internal/              # Internal packages
│   ├── app/               # Application layer code
│   │   ├── config/        # Configuration management
│   │   ├── handlers/      # HTTP handlers
│   │   ├── middleware/    # Middleware
│   │   ├── models/        # Data models
│   │   ├── repositories/  # Data access layer
│   │   ├── request/       # Request structs
│   │   ├── router/        # Router configuration
│   │   └── services/      # Business logic layer
│   └── pkg/               # Common packages
│       ├── cache/         # Cache operations
│       ├── cronjob/       # Scheduled tasks
│       ├── db/            # Database connection
│       ├── kit/           # Utility functions
│       ├── mailer/        # Email service
│       ├── oss/           # Object storage
│       └── token/         # JWT token
├── config.example.ini     # Configuration file template
├── Makefile              # Build scripts
└── README.md             # Project documentation
```

> 📖 Reference: [Go Standard Project Layout](https://github.com/golang-standards/project-layout)

## 🚀 Quick Start

### 1. Clone the Project

```bash
git clone https://github.com/artfoxe6/quick-gin.git
cd quick-gin
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configuration Setup

Copy the configuration template and modify as needed:

```bash
cp config.example.ini config.ini
```

#### Database Configuration

**MySQL Configuration:**
```ini
[database]
Type = mysql
Host = localhost
Port = 3306
User = root
Password = your_password
Db = your_database
```

**SQLite Configuration:**
```ini
[database]
Type = sqlite
DbFile = data/app.db
```

#### Cache Configuration

**Redis Cache:**
```ini
[cache]
Type = redis
Host = 127.0.0.1
Port = 6379
Password =
Db = 1
```

**Memory Cache:**
```ini
[cache]
Type = memory
```

### 4. Start the Service

```bash
# Development mode
go run cmd/app/main.go -config config.ini

# Or use Makefile
make run
```

The service will start at `http://localhost:8080`.

## 🛠️ CRUD Code Generation

Generate complete CRUD code by specifying a module name:

```bash
# Generate complete code for Product module
go run cmd/gen/gen.go --module=Product
```

This will automatically generate:
- Model: `internal/app/models/product.go`
- Repository: `internal/app/repositories/product.go`
- Service: `internal/app/services/product.go`
- Handler: `internal/app/handlers/product.go`
- Router registration code

## 🔧 Configuration Details

### Application Configuration

```ini
[app]
AppMode = debug          ; Running mode: debug/release
Listen = :8080          ; Listen address
ReadTimeout = 60        ; Read timeout (seconds)
WriteTimeout = 60       ; Write timeout (seconds)
LogDir = log            ; Log directory
SignKey = your_sign_key ; Signature key
```

### JWT Configuration

```ini
[jwt]
Secret = your_jwt_secret     ; JWT secret key
Exp = 240                    ; Access token expiration (hours)
RefreshExp = 720             ; Refresh token expiration (hours)
```

### Database Detailed Configuration

```ini
[database]
Type = mysql                 ; Database type: mysql/sqlite
ConnMaxLifeTime = 15         ; Connection max lifetime (minutes)
MaxPoolSize = 10             ; Max connection pool size
MaxIdle = 10                 ; Max idle connections
```

### Cache Detailed Configuration

```ini
[cache]
Type = redis                 ; Cache type: memory/redis
MaxIdle = 5                  ; Redis max idle connections
MaxActive = 10               ; Redis max active connections
```

## 📚 Available APIs

### User Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh token
- `POST /api/auth/logout` - User logout

### User Management
- `GET /api/users` - Get user list
- `GET /api/users/:id` - Get user details
- `PUT /api/users/:id` - Update user information
- `DELETE /api/users/:id` - Delete user

### News Management
- `GET /api/news` - Get news list
- `POST /api/news` - Create news
- `GET /api/news/:id` - Get news details
- `PUT /api/news/:id` - Update news
- `DELETE /api/news/:id` - Delete news

## 🔨 Build Commands

Use Makefile to manage the project:

```bash
make build    # Build application
make run      # Start service
make stop     # Stop service
make restart  # Restart service
make test     # Run tests
```

## 🧩 Integrated Components

- **Web Framework**: Gin v1.6.0
- **ORM**: GORM v1.25.1
- **Authentication**: JWT v5.0.0
- **Cache**: Redis v2.0.0 / Memory Cache
- **Database**: MySQL v1.5.0 / SQLite
- **Configuration**: Go-ini v1.55.0
- **Object Storage**: Aliyun OSS v3.0.2
- **Email Service**: SendGrid v3.16.0
- **Scheduled Tasks**: Cron v3.0.1

## 🎯 Use Cases

- **Microservices Architecture** - As an independent service providing API interfaces
- **Rapid Prototyping** - Quickly validate ideas and concepts
- **Enterprise Applications** - Complete authentication and authorization
- **Learning Projects** - Learn Go web development best practices

## 🤝 Contributing

Issues and Pull Requests are welcome!

## 📄 License

MIT License