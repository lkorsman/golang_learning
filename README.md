# Go Store API

A REST API built with Go for learning purposes. Features JWT authentication, MySQL persistence, and full CRUD operations for managing products.

## Features

- **JWT Authentication** - Secure user registration and login
- **MySQL Database** - Persistent storage with context-aware queries
- **Request Validation** - Input validation for all endpoints
- **Comprehensive Testing** - Table-driven tests and benchmarks
- **Structured Logging** - Clean, readable logs
- **Middleware** - Request timing, authentication, recovery
- **Clean Architecture** - Organized internal packages

## Tech Stack

- **Framework**: [Chi Router](https://github.com/go-chi/chi)
- **Database**: MySQL with `database/sql`
- **Authentication**: JWT tokens with `golang-jwt/jwt`
- **Password Hashing**: bcrypt
- **Config**: Environment variables

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
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── http/
│   │   ├── context.go        # Context utilities
│   │   └── middleware.go     # HTTP middleware
│   └── product/
│       ├── handler.go        # Product endpoints
│       ├── mysql_store.go    # MySQL implementation
│       ├── product.go        # Product model
│       ├── store.go          # Store interface
│       └── validation.go     # Input validation
├── .env                      # Environment variables
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/lkorsman/golang_learning.git
cd golang_learning
```

2. **Install dependencies**
```bash
go mod download
```

3. **Start MySQL**
```bash
brew services start mysql
# or
mysql.server start
```

4. **Create the database**
```bash
mysql -u root -p
CREATE DATABASE store;
exit
```

5. **Configure environment variables**

Create a `.env` file in the project root:
```env
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
DATABASE_URL=root:yourpassword@tcp(localhost:3306)/store?parseTime=true
ENVIRONMENT=development
```

6. **Run the application**
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

### Code Formatting
```bash
# Format all code
go fmt ./...

# Check for common mistakes
go vet ./...
```

### Running with In-Memory Store

To run without MySQL (useful for development/testing), simply comment out or remove the `DATABASE_URL` from your `.env` file:

```env
# DATABASE_URL=root:yourpassword@tcp(localhost:3306)/store?parseTime=true
```

The app will automatically fall back to an in-memory store.

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
| `JWT_SECRET` | Secret key for JWT signing | `your-secret-key-change-in-production` |
| `ENVIRONMENT` | Environment (development/production) | `development` |

## What I Learned

This project helped me learn:
- Building REST APIs with Go
- Chi router and middleware patterns
- JWT authentication and authorization
- MySQL integration with `database/sql`
- Context-aware database queries
- Password hashing with bcrypt
- Input validation
- Table-driven testing in Go
- Interface-based design
- Clean architecture patterns
- Graceful server shutdown

## Future Improvements

- [ ] Add refresh tokens
- [ ] Implement role-based access control (RBAC)
- [ ] Add pagination to product listing
- [ ] Store users in MySQL instead of memory
- [ ] Add API rate limiting
- [ ] Implement CORS middleware
- [ ] Add Docker support
- [ ] Set up CI/CD pipeline
- [ ] Add API documentation with Swagger

## License

This is a learning project - feel free to use it however you'd like!

## Acknowledgments

Built while learning Go through hands-on practice and experimentation.