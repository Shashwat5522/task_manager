# Task Manager REST API

A RESTful API for managing tasks built with Go using Clean Architecture and raw SQL.

## Project Structure

```
task_manager/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── config/
│   └── config.go                # Configuration structures
├── internal/
│   ├── domain/                  # Domain entities
│   │   ├── user.go
│   │   └── task.go
│   ├── dto/                     # Data Transfer Objects
│   │   ├── auth_dto.go
│   │   └── task_dto.go
│   ├── repository/              # Data access interfaces
│   │   ├── user_repository.go
│   │   └── task_repository.go
│   ├── service/                 # Business logic interfaces
│   │   ├── auth_service.go
│   │   └── task_service.go
│   ├── handler/                 # HTTP handlers
│   │   ├── auth_handler.go
│   │   └── task_handler.go
│   └── middleware/              # HTTP middleware
│       ├── auth_middleware.go
│       ├── logger_middleware.go
│       └── recovery_middleware.go
├── pkg/                         # Shared packages
│   ├── database/
│   │   └── postgres.go
│   ├── logger/
│   │   └── logger.go
│   └── utils/
│       ├── jwt.go
│       ├── password.go
│       └── response.go
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

## Architecture

The project follows **Clean Architecture** with clear separation of concerns:

- **cmd/**: Application entry points
- **config/**: Configuration management
- **internal/domain/**: Core business entities
- **internal/dto/**: Request/Response data structures
- **internal/repository/**: Data access layer (interfaces)
- **internal/service/**: Business logic layer (interfaces)
- **internal/handler/**: HTTP handlers
- **internal/middleware/**: HTTP middleware
- **pkg/**: Reusable packages

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL with sqlx (raw SQL)
- **Authentication**: JWT

## Setup

1. Copy environment file:
```bash
cp .env.example .env
```

2. Edit `.env` with your configuration

3. Initialize Go module:
```bash
go mod init github.com/vedologic/task-manager
go mod tidy
```

## Next Steps

- Implement database connection
- Implement repository layer with raw SQL
- Implement service layer with business logic
- Implement handlers
- Setup routes and middleware
- Add database migrations
- Add tests
