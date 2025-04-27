# quick-gin
A scaffold based on Gin, designed to help you start writing interface code as quickly as possible without wasting time setting up the project framework.

### Overview
- Provides a typical web architecture: model, repository, service, and handler
- Integrates commonly used components such as Redis and MySQL.
- Generates CRUD template code with one click.

### Directory Structure
```bash
├── cmd
├── internal
│   ├── app
│   │   ├── handlers
│   │   ├── middleware
│   │   ├── models
│   │   ├── repositories
│   │   ├── request
│   │   ├── router
│   │   └── services
│   └── pkg
```

> refer: [GitHub - golang-standards/project-layout: Standard Go Project Layout](https://github.com/golang-standards/project-layout)

### Quick Start

- Clone the project locally
  ```
  git clone https://github.com/artfoxe6/quick-gin.git
  ```

- Create a new project based on quick-gin

  ```
  go run cmd/clone/main.go --path=/opt/myapp --package=github.com/yourname/myapp    
  ```
- Navigate to the new project root directory
- Start the service:

  ```
  go run cmd/main.go --config=./config.example.ini  
  ```

  If your system has `make` installed, you can also use the following commands

  ```
  make build   # Build the service
  make run     # Start the service
  make stop    # Stop the service
  ```

### CRUD Code Generation
Specify a module name to automatically generate related code: router → handler → service → repository → model.
```
go run cmd/gen/gen.go --module=Product 
```

### Integrated Components
- Redis
- MySQL
- JWT
- OSS (Object Storage)
- SendGrid (Email)
