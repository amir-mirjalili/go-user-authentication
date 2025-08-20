# OTP Authentication Service

A robust, production-ready OTP-based authentication service built with Go, featuring user management, rate limiting, and JWT-based security.

## Features

- ✅ **OTP-based Authentication**: Secure login/registration using phone numbers
- ✅ **Rate Limiting**: Prevents OTP spam (max 3 requests per 10 minutes)
- ✅ **JWT Security**: Token-based authentication with expiration
- ✅ **User Management**: Complete CRUD operations with pagination and search
- ✅ **Database Support**: PostgreSQL
- ✅ **API Documentation**: Comprehensive Swagger/OpenAPI documentation
- ✅ **Containerization**: Docker and docker-compose ready
- ✅ **Clean Architecture**: Separation of concerns with clear project structure

## Architecture

The service follows a clean architecture pattern with clear separation of responsibilities:

```
├── cmd/                 # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Authentication middleware
│   ├── models/          # Data models 
│   ├── services/        # Business logic layer
│   └── db/              # Data access layer
    ├── params/          # DTOs
├── pkg/jwt/             # JWT utilities
└── docs/                # API documentation
```

## Database Choice: PostgreSQL

**Why PostgreSQL?**

1. **ACID Compliance**: Essential for user authentication data integrity
2. **Concurrent Safety**: Handles multiple OTP requests safely with proper locking
3. **Performance**: Excellent performance for read-heavy operations (user lookups)
4. **Indexing**: Efficient searching and pagination with proper indexes


The service also includes an in-memory storage option for development and testing.

## Quick Start

### Local Development (with Docker)

1. **Clone and start the services:**
```bash
git clone github.com/amir-mirjalili/go-user-authentication
cd go-user-authentication
docker-compose up --build
```

2. **The service will be available at:**
- API: http://localhost:8080
- Swagger Documentation: http://localhost:8080/swagger/index.html
- Health Check: http://localhost:8080/health

### Local Development (without Docker)

1. **Install dependencies:**
```bash
go mod download
```

2. **Set environment variables:**
```bash
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_SSL_MODE=
JWT_SECRET=
```

3. **Initialize Swagger Documents:**
```bash
 swag init -g cmd/main.go -o ./docs
```

4. **Run the service:**
```bash
go run cmd/main.go
```

## API Usage Examples

### 1. Send OTP

```bash
curl -X POST http://localhost:8080/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+989123456789"}'
```

**Response:**
```json
{
  "message": "OTP sent successfully"
}
```

### 2. Verify OTP (Login/Register)

```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+1234567890",
    "code": "123456"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "phone_number": "+1234567890",
    "registered_at": "2025-08-18T10:30:00Z",
    "created_at": "2025-08-18T10:30:00Z",
    "updated_at": "2025-08-18T10:30:00Z"
  }
}
```

### 3. Get Current User

```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. List Users (with pagination and search)

```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=10&search=+123" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "phone_number": "+1234567890",
      "registered_at": "2025-08-18T10:30:00Z",
      "created_at": "2025-08-18T10:30:00Z",
      "updated_at": "2025-08-18T10:30:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```

### 5. Get User by ID

```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Rate Limiting

- **OTP Requests**: Maximum 3 requests per phone number within 10 minutes
- **Rate limit exceeded**: Returns HTTP 429 with appropriate error message


## Development Commands

```bash
# Install dependencies
go mod download

#Initialize Swagger documentation
 swag init -g cmd/main.go -o ./docs

# Run locally
go run cmd/main.go

# Build binary
go build -o bin/ cmd/main.go

# Build Docker image
docker build -t otp-auth-service .

# Run with docker-compose
docker-compose up --build

# Stop services
docker-compose down

# View logs
docker-compose logs -f app
```
