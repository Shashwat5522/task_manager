# Task Manager API

A production-ready REST API for managing tasks and users, built with Go, PostgreSQL, and best practices for database migrations, authentication, and API documentation.

## Overview

The Task Manager API is a robust backend service that provides:
- **User Management**: Registration and secure authentication
- **Task Management**: Full CRUD operations with bulk actions
- **JWT Authentication**: Secure token-based access control
- **Database Migrations**: Version-controlled schema management with automatic migrations
- **Swagger UI**: Interactive API documentation
- **Structured Logging**: Comprehensive request/response logging
- **Error Handling**: Graceful error management with informative responses

## Features

### Core Functionality
- ✅ User registration with email and password validation
- ✅ JWT-based authentication and authorization
- ✅ Task CRUD operations (Create, Read, Update, Delete)
- ✅ Task filtering by status and pagination
- ✅ Bulk task completion operations
- ✅ Automatic database migrations on startup
- ✅ Schema verification and integrity checks

### Technical Highlights
- Clean Architecture with separation of concerns (handlers, services, repositories, domain)
- Idempotent operations for reliability
- PostgreSQL with SQLX for type-safe queries
- Automated migration versioning using `golang-migrate`
- Structured logging with Zap
- Middleware for authentication, logging, and recovery
- Comprehensive API documentation with Swagger/OpenAPI

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.20+ |
| Framework | Gin Web Framework |
| Database | PostgreSQL 12+ |
| Authentication | JWT (JSON Web Tokens) |
| Password Hashing | bcrypt |
| Database Driver | SQLX |
| Migrations | golang-migrate |
| Logging | Uber Zap |
| API Documentation | Swagger/OpenAPI |


## Installation

### Prerequisites
- Go 1.20 or higher
- PostgreSQL 12 or higher
- Git

### Setup Instructions

1. **Clone the repository**
```bash
git clone <repository-url>
cd task_manager
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment variables**
Create a `.env` file or export variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=task_manager
export SERVER_PORT=8080
export SERVER_ENV=development
export JWT_SECRET=your-secret-key-here
```

4. **Initialize the database**
PostgreSQL will be automatically created and migrations will run on startup.

5. **Run the application**
```bash
go run ./cmd/api/main.go
```

The API will start on `http://localhost:8080`

## API Documentation

### Swagger UI
Access interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

### Authentication
All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Key Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - User login and token generation

### Tasks
- `GET /api/v1/tasks` - List all tasks (paginated)
- `POST /api/v1/tasks` - Create a new task
- `GET /api/v1/tasks/{id}` - Get a specific task
- `PUT /api/v1/tasks/{id}` - Update a task
- `DELETE /api/v1/tasks/{id}` - Delete a task
- `POST /api/v1/tasks/bulk-complete` - Mark multiple tasks as complete

## Task Status Values

Tasks support three status values:
- `todo` - Task is pending
- `in_progress` - Task is currently being worked on
- `done` - Task has been completed

## Database Migrations

### Automatic Migration on Startup
Migrations run automatically when the application starts. The system:
- Checks for existing schema and tables
- Applies pending migrations
- Verifies schema integrity
- Logs migration status

### Manual Migration Management
```bash
# Check migration status
go run ./cmd/api/main.go --check-migrations

# View migration details
go run ./cmd/api/main.go --migration-status
```

### Creating New Migrations
```bash
migrate create -ext sql -dir migrations -seq create_new_table
```

Then edit the generated `.up.sql` and `.down.sql` files.

## Request/Response Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

### Create Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project",
    "description": "Finish the task manager API",
    "status": "todo"
  }'
```

### List Tasks
```bash
curl -X GET "http://localhost:8080/api/v1/tasks?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Error message describing what went wrong"
}
```

HTTP Status Codes:
- `200 OK` - Successful request
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid authentication
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Logging

The application uses Zap for structured logging. Log levels:
- `DEBUG` - Detailed information for debugging
- `INFO` - General information messages
- `WARN` - Warning messages
- `ERROR` - Error messages

Set log level via `SERVER_ENV`:
- `development` - Debug level logging
- `production` - Info level logging

## Architecture

### Clean Architecture Principles
The project follows clean architecture with clear separation of concerns:

```
HTTP Request
    ↓
Handler (HTTP Layer)
    ↓
Service (Business Logic)
    ↓
Repository (Data Access)
    ↓
Database (PostgreSQL)
```

### Layers

1. **Handler Layer**: Manages HTTP requests/responses
2. **Service Layer**: Implements business logic and validations
3. **Repository Layer**: Handles database operations
4. **Domain Layer**: Contains business entities and models
5. **Middleware Layer**: Cross-cutting concerns (auth, logging, recovery)

## Security Features

- ✅ Password hashing with bcrypt
- ✅ JWT token-based authentication
- ✅ Token expiration and validation
- ✅ Protected endpoints with middleware
- ✅ Secure headers and CORS handling
- ✅ Input validation and sanitization
- ✅ SQL injection prevention with prepared statements

## Performance Considerations

- Database query optimization with indexes
- Pagination support for list endpoints
- Connection pooling with SQLX
- Efficient error handling
- Structured logging without performance impact

## Git Workflow

The project follows a gradual development strategy with focused commits:

```bash
# Each feature/component has dedicated commits:
git log --oneline
# - Fix domain models: change ID types from string to int
# - Update task status values to todo, in_progress, done
# - Improve logging output in main.go
# - Improve route registration logging
# - Implement HTTP middleware
# - Implement HTTP handlers
# - Fix bulk complete operations
# - Add Swagger API documentation
```

## Troubleshooting

### Database Connection Issues
```
Error: failed to connect to database
Solution: Verify PostgreSQL is running and credentials are correct
```

### Migration Errors
```
Error: failed to run migrations
Solution: Check migrations directory exists and files are valid SQL
```

### JWT Token Issues
```
Error: 401 Unauthorized
Solution: Verify token is included in Authorization header with Bearer prefix
```

## Contributing

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Make changes with focused commits
3. Push to the repository: `git push origin feature/your-feature`
4. Create a Pull Request

## Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Deployment

### Docker Deployment
```bash
docker build -t task-manager-api .
docker run -p 8080:8080 --env-file .env task-manager-api
```

### Production Considerations
- Use environment variables for sensitive data
- Enable HTTPS with proper certificates
- Configure appropriate CORS policies
- Set up monitoring and alerting
- Use a reverse proxy (nginx/traefik)
- Enable database backups

## License

This project is provided as a practical task implementation.

## Support

For issues or questions:
1. Check the Swagger documentation at `/swagger/index.html`
2. Review the error messages in logs
3. Verify environment configuration
4. Check database connectivity

## Version History

- **1.0.0** - Initial release with full API functionality, JWT auth, and Swagger docs