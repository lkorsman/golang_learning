# Go Store API

A production-ready REST API built with Go. Features JWT authentication, MySQL persistence, Redis caching, database migrations, comprehensive observability, and full CRUD operations for managing products.

## Features

- **JWT Authentication** - Secure user registration and login with bcrypt password hashing
- **MySQL Database** - Persistent storage with context-aware queries and migrations
- **Redis Caching** - Lightning-fast responses with intelligent cache invalidation
- **Database Migrations** - Version-controlled schema management with golang-migrate
- **Observability** - Prometheus metrics for monitoring performance and health
- **Docker Support** - Complete containerization with docker-compose
- **Request Validation** - Input validation for all endpoints
- **Comprehensive Testing** - Table-driven tests and benchmarks
- **Structured Logging** - Clean, readable logs
- **Middleware** - Request timing, authentication, recovery, metrics
- **Clean Architecture** - Organized internal packages

## Tech Stack

- **Framework**: [Chi Router](https://github.com/go-chi/chi)
- **Database**: MySQL with `database/sql`
- **Cache**: Redis with `go-redis`
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Metrics**: Prometheus client library
- **Authentication**: JWT tokens with `golang-jwt/jwt`
- **Password Hashing**: bcrypt
- **Containerization**: Docker & Docker Compose

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── auth/
│   │   ├── handler.go        # Auth endpoints (register, login)
│   │   ├── jwt.go            # JWT token management
│   │   ├── store.go          # User storage
│   │   └── user.go           # User model
│   ├── cache/
│   │   └── redis.go          # Redis cache implementation
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── database/
│   │   └── migrate.go        # Database migration runner
│   ├── http/
│   │   ├── context.go        # Context utilities
│   │   └── middleware.go     # HTTP middleware (auth, metrics, timing)
│   ├── metrics/
│   │   └── metrics.go        # Prometheus metrics definitions
│   └── product/
│       ├── handler.go        # Product endpoints
│       ├── mysql_store.go    # MySQL implementation
│       ├── product.go        # Product model
│       ├── store.go          # Store interface
│       └── validation.go     # Input validation
├── migrations/
│   ├── 000001_create_products_table.up.sql
│   ├── 000001_create_products_table.down.sql
│   └── ...                   # Migration files
├── .env                      # Environment variables
├── .env.example              # Environment template
├── docker-compose.yml        # Docker services configuration
├── Dockerfile                # Multi-stage Docker build
├── Makefile                  # Common commands
├── prometheus.yml            # Prometheus configuration
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- MySQL 5.7 or higher

### Installation

#### Option 1: Docker (Recommended)

1. **Clone the repository**
```bash
git clone https://github.com/lkorsman/golang_learning.git
cd golang_learning
```

2. **Start all services with Docker Compose**
```bash
make docker-up
```

This will:
- Start MySQL container
- Start Redis container
- Build and start the API container
- Run database migrations automatically
- Expose the API on http://localhost:8080

3. **View logs**
```bash
make docker-logs
```

#### Option 2: Local Development

1. **Clone the repository**
```bash
git clone https://github.com/lkorsman/golang_learning.git
cd golang_learning
```

2. **Install dependencies**
```bash
go mod download

# Install migration tool
brew install golang-migrate
```

3. **Start MySQL**
```bash
brew services start mysql
```

4. **Start Redis**
```bash
brew services start redis
```

5. **Create the database**
```bash
mysql -u root -p
CREATE DATABASE store;
exit
```

6. **Configure environment variables**

Create a `.env` file in the project root:
```env
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
DATABASE_URL=root:yourpassword@tcp(localhost:3306)/store?parseTime=true
REDIS_URL=localhost:6379
ENVIRONMENT=development
```

7. **Run migrations**
```bash
make migrate-up
```

8. **Run the application**
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Authentication

#### Register a New User
```bash
POST /auth/register

# Request
{
  "email": "alice@example.com",
  "password": "password123"
}

# Response (201 Created)
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "alice@example.com"
  }
}
```

#### Login
```bash
POST /auth/login

# Request
{
  "email": "alice@example.com",
  "password": "password123"
}

# Response (200 OK)
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "alice@example.com"
  }
}
```

### Products

**Note:** Create, Update, and Delete operations require JWT authentication via the `Authorization: Bearer <token>` header.

#### List All Products
```bash
GET /products

# Response (200 OK)
[
  {
    "id": 1,
    "name": "Laptop",
    "price": 1200.50
  },
  {
    "id": 2,
    "name": "Mouse",
    "price": 25.99
  }
]
```

#### Get Single Product
```bash
GET /products/{id}

# Response (200 OK)
{
  "id": 1,
  "name": "Laptop",
  "price": 1200.50
}
```

#### Create Product (Protected)
```bash
POST /products
Authorization: Bearer <your-jwt-token>

# Request
{
  "name": "Laptop",
  "price": 1200.50
}

# Response (201 Created)
{
  "id": 1,
  "name": "Laptop",
  "price": 1200.50
}
```

#### Update Product (Protected)
```bash
PUT /products/{id}
Authorization: Bearer <your-jwt-token>

# Request
{
  "name": "Gaming Laptop",
  "price": 1500.00
}

# Response (200 OK)
{
  "id": 1,
  "name": "Gaming Laptop",
  "price": 1500.00
}
```

#### Delete Product (Protected)
```bash
DELETE /products/{id}
Authorization: Bearer <your-jwt-token>

# Response (204 No Content)
```

#### Prometheus Metrics
```bash
GET /metrics

# Response: Prometheus format metrics
# http_requests_total - Total HTTP requests by method, path, status
# http_request_duration_seconds - Request latency histogram
# cache_hits_total - Cache hits by key
# cache_misses_total - Cache misses by key
# db_queries_total - Database queries by operation
# db_query_duration_seconds - Database query latency
# products_created_total - Total products created
# products_deleted_total - Total products deleted
# user_registrations_total - Total user registrations
# login_attempts_total - Login attempts (success/failure)
```

## Example Usage

### Complete Flow Example

```bash
# 1. Register a user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"password123"}'

# Copy the token from the response

# 2. Create a product (use your token)
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGc..." \
  -d '{"name":"Laptop","price":1200.50}'

# 3. List all products (no auth needed)
curl http://localhost:8080/products

# 4. Get a specific product
curl http://localhost:8080/products/1

# 5. Update a product (requires auth)
curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGc..." \
  -d '{"name":"Gaming Laptop","price":1500.00}'

# 6. Delete a product (requires auth)
curl -X DELETE http://localhost:8080/products/1 \
  -H "Authorization: Bearer eyJhbGc..."
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test ./... -bench=. -benchmem

# Run specific package tests
go test ./internal/product -v
```

## Development

### Database Migrations

```bash
# Create a new migration
make migrate-create
# Enter migration name when prompted

# Run all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Check current migration version
make migrate-version

# Force a specific version (if stuck)
make migrate-force
```

### Cache Management

```bash
# Connect to Redis CLI
make redis-cli

# Common Redis commands:
KEYS *              # List all keys
GET products:list   # Get a cached value
TTL products:list   # Check time to live
FLUSHALL           # Clear all cache
```

### Observability

```bash
# View metrics
curl http://localhost:8080/metrics

# Example Prometheus queries (if using Prometheus):
# - rate(http_requests_total[1m])           # Requests per second
# - http_request_duration_seconds           # Latency percentiles
# - cache_hits_total / (cache_hits_total + cache_misses_total)  # Cache hit rate./...
```

### Running with In-Memory Store

To run without MySQL (useful for development/testing), simply comment out or remove the `DATABASE_URL` from your `.env` file:

```env
# DATABASE_URL=root:yourpassword@tcp(localhost:3306)/store?parseTime=true
```

The app will automatically fall back to an in-memory store.

### Running without Redis

To run without Redis caching, comment out or remove the `REDIS_URL`:

```env
# REDIS_URL=localhost:6379
```

The app will log a warning and continue without caching.

## Validation Rules

### Product
- **name**: Required, max 100 characters
- **price**: Required, must be > 0 and < 999,999.99

### User Registration
- **email**: Required
- **password**: Required, minimum 6 characters

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | MySQL connection string | _(empty - uses in-memory)_ |
| `REDIS_URL` | Redis connection string | `localhost:6379` |
| `JWT_SECRET` | Secret key for JWT signing | `your-secret-key-change-in-production` |
| `ENVIRONMENT` | Environment (development/production) | `development` |

## What I Learned

This project helped me learn:
- Building REST APIs with Go and Chi router
- JWT authentication and authorization
- MySQL integration with `database/sql`
- Context-aware database queries
- Password hashing with bcrypt
- Input validation patterns
- Table-driven testing in Go
- Interface-based design for testability
- Clean architecture patterns
- Graceful server shutdown
- **Redis caching** - Cache-aside pattern, TTL, invalidation strategies
- **Database migrations** - Schema versioning with golang-migrate
- **Docker & Docker Compose** - Multi-stage builds, containerization
- **Observability** - Prometheus metrics, monitoring production systems
- **Middleware patterns** - Request timing, metrics collection, authentication
- **Concurrency patterns** - Safe concurrent access with mutexes

## Future Improvements

- [ ] Add refresh tokens
- [ ] Implement role-based access control (RBAC)
- [ ] Add pagination to product listing
- [ ] Store users in MySQL instead of memory
- [ ] Add API rate limiting
- [ ] Implement CORS middleware
- [X] Add Docker support
- [ ] Set up CI/CD pipeline
- [ ] Add API documentation with Swagger

## License

This is a learning project - feel free to use it however you'd like!

## Acknowledgments

Built while learning Go through hands-on practice and experimentation.