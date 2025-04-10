# quick-gin

A scaffold built on top of gin`to significantly improve development efficiency.


### Overview

Gin is a minimalist and high-performance web framework. However, when building a typical web project with Gin, you usually need to manually create directories such as `model`, `service`, `controller`, and others. You'll also have to set up a mechanism to connect these modules before you can start actual business development.

Additionally, web development often requires integrating common components like Redis, MySQL, etc., which can be time-consuming to configure and integrate into your project.

This is where **quick-gin** comes in. It helps you kickstart a new project quickly, providing a well-organized directory structure and reducing setup overhead.

Moreover, quick-gin includes one-click CRUD code generation, making development even faster. The following sections describe how to use it in detail.


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

The structure under the `app` directory is self-explanatory. One point worth noting is the use of the `internal` directory.

In Go, the `internal` directory has a special meaning—it restricts access to the code inside it from outside the module. This aligns well with typical web applications, which are not designed to be shared as libraries.

For directory structure conventions, refer to:

> [GitHub - golang-standards/project-layout: Standard Go Project Layout](https://github.com/golang-standards/project-layout)

### Getting Started

- Clone this repository to your local machine.

- Create a new project based on quick-gin:

  ```
  go run cmd/clone/main.go --path=/opt/myapp --package=github.com/yourname/myapp
  ```

- Start the service:

  ```
  go run cmd/main.go --config=./config.example.ini
  ```

  If your system has `make` installed, you can also use the following commands:

  ```
  make build   # Build the service
  make run     # Start the service
  make stop    # Stop the service
  ```

### Code Generation

```
go run cmd/gen/gen.go --module=Product
```

After running the command successfully, the following files will be automatically generated. You can then start modifying the relevant business logic:

- model

- repository

- service

- handler

- route

### Integrated Components

- Redis

- MySQL

- JWT

- OSS (Object Storage)

- SendGrid (Email)
