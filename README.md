# Gin Boilerplate

A clean, well-structured Go Gin boilerplate with authentication, user management, and essential features for building REST APIs.

## 🚀 Features

- **Clean Architecture**: Organized with separation of concerns
- **Authentication & Authorization**: JWT-based authentication system
- **User Management**: Complete CRUD operations for users
- **Database Integration**: PostgreSQL with GORM ORM
- **Redis Support**: Optional Redis integration for caching
- **Rate Limiting**: Built-in rate limiting middleware
- **Input Sanitization**: XSS and SQL injection protection with configurable sanitization
- **Request Size Limiting**: Protection against DoS attacks through large payloads
- **Error Handling**: Centralized error handling with custom error types
- **Validation**: Request validation using go-playground/validator
- **Logging**: Structured logging with Logrus
- **Configuration**: Environment-based configuration with Viper
- **API Documentation**: Swagger/OpenAPI documentation
- **Graceful Shutdown**: Proper server shutdown handling
- **CORS Support**: Cross-Origin Resource Sharing configuration
- **Health Checks**: Health check endpoints for monitoring

## 📁 Project Structure

```
gin-boilerplate/
├── boot/                   # Dependency injection and setup
├── config/                 # Configuration files
├── docs/                   # API documentation
├── infrastructure/         # Infrastructure layer
│   ├── config/            # Configuration management
│   ├── database/          # Database connection
│   ├── errors/            # Error handling
│   ├── httplib/           # HTTP utilities
│   ├── jwt/               # JWT authentication
│   ├── log/               # Logging
│   ├── middleware/        # Middleware components (auth, rate limiting, security)
│   ├── redis/             # Redis client
│   ├── sanitizer/         # Input sanitization utilities
│   └── validator/         # Request validation
├── modules/               # Business logic modules
│   ├── auth/              # Authentication module
│   ├── health/            # Health check module
│   ├── primitive/         # Common types and structures
│   └── user/              # User management module
├── router/                # Route definitions
├── utils/                 # Utility functions
├── main.go               # Application entry point
├── Makefile              # Build and development commands
└── go.mod                # Go module dependencies
```

## 🛠️ Prerequisites

- Go 1.23 or higher
- PostgreSQL 12 or higher
- Redis (optional)

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone <repository-url>
cd gin-boilerplate
```

### 2. Install dependencies

```bash
make deps
```

### 3. Setup database

Create a PostgreSQL database:

```sql
CREATE DATABASE gin_boilerplate;
```

### 4. Configure environment

Copy and modify the configuration file:

```bash
cp config.local.yaml config.local.yaml
```

Update the database configuration in `config.local.yaml`:

```yaml
postgres:
  host: "localhost"
  port: 5432
  user: "your_username"
  password: "your_password"
  dbname: "gin_boilerplate"
  sslmode: "disable"
  timezone: "UTC"
```

### 5. Run the application

```bash
make run
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

Once the server is running, you can access the Swagger documentation at:
`http://localhost:8080/swagger/index.html`

## 🔧 Available Commands

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code
make lint

# Generate swagger docs
make swagger

# Clean build files
make clean

# Install dependencies
make deps
```

## 🌍 Environment Configuration

The application supports multiple environments:

- `local` - Local development (default)
- `dev` - Development environment
- `uat` - User Acceptance Testing
- `prod` - Production environment

Run with specific environment:

```bash
go run main.go -env=dev
```

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Authentication Endpoints

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/profile` - Get user profile (requires auth)
- `PUT /api/v1/auth/profile` - Update user profile (requires auth)
- `POST /api/v1/auth/change-password` - Change password (requires auth)

## 👥 User Management

### User Endpoints

- `GET /api/v1/users` - Get all users (requires auth)
- `GET /api/v1/users/:id` - Get user by ID (requires auth)
- `POST /api/v1/users` - Create new user (requires auth)
- `PUT /api/v1/users/:id` - Update user (requires auth)
- `DELETE /api/v1/users/:id` - Delete user (requires auth)

## 🏥 Health Checks

- `GET /api/v1/health/ping` - Simple ping endpoint
- `GET /api/v1/health/check` - Comprehensive health check

## 🔒 Security Features

### Input Sanitization

The application includes comprehensive input sanitization to protect against XSS and SQL injection attacks:

- **XSS Protection**: Automatically detects and sanitizes script tags, event handlers, and dangerous HTML
- **SQL Injection Protection**: Detects common SQL injection patterns and blocks malicious requests
- **HTML Sanitization**: Removes or escapes dangerous HTML tags and characters
- **Configurable**: Can be enabled/disabled per environment with strict or lenient modes

### Request Size Limiting

Protection against DoS attacks through request size limitations:

- **Body Size Limiting**: Configurable maximum request body size
- **Form Field Limiting**: Limits number of form fields and files
- **Header Size Limiting**: Prevents oversized headers
- **URL Length Limiting**: Restricts maximum URL length
- **Query Parameter Limiting**: Limits number of query parameters

### Configuration

Security features can be configured in your `config.yaml`:

```yaml
security:
  # Input Sanitization
  enableSanitization: true
  enableXSSDetection: true
  enableSQLInjectionDetection: true
  maxStringLength: 1000
  strictMode: true

  # Request Size Limiting
  maxBodySize: 10485760 # 10MB
  maxFormFields: 100
  maxFormMemory: 33554432 # 32MB
  maxFormFiles: 10
  maxHeaderSize: 1048576 # 1MB
  maxURLLength: 2048
  maxQueryParams: 50
  maxFileSize: 52428800 # 50MB
  maxResponseSize: 104857600 # 100MB
```

### Environment-Specific Security

- **Production**: More restrictive limits for enhanced security
- **Development**: Balanced settings for development workflow
- **Testing**: Relaxed settings to facilitate testing

## 🔧 Configuration Options

Key configuration options in `config.yaml`:

```yaml
# Server configuration
port: 8080
logLevel: "DEBUG"
logMode: true
logFormat: "text"

# Database configuration
postgres:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  dbname: "gin_boilerplate"
  sslmode: "disable"
  timezone: "UTC"

# Redis configuration (optional)
redis:
  enableRedis: false
  host: "localhost"
  port: 6379
  password: ""
  db: 0

# Rate limiting
rate: 100
interval: "1m"

# JWT configuration
secretKey: "your-secret-key"
accessTokenExpiry: 24
accessTokenExpiryUnit: "hour"

# Security configuration
security:
  enableSanitization: true
  enableXSSDetection: true
  enableSQLInjectionDetection: true
  maxStringLength: 1000
  strictMode: true
  maxBodySize: 10485760 # 10MB
  maxFormFields: 100
  maxFormMemory: 33554432 # 32MB
  maxFormFiles: 10
  maxHeaderSize: 1048576 # 1MB
  maxURLLength: 2048
  maxQueryParams: 50
  maxFileSize: 52428800 # 50MB
  maxResponseSize: 104857600 # 100MB
```

## 🧪 Testing

Run tests:

```bash
make test
```

Run tests with coverage:

```bash
make test-coverage
```

### Security Testing

Test security features with sample requests:

```bash
# Test XSS protection
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "<script>alert(\"xss\")</script>", "email": "test@example.com", "password": "password123"}'

# Test SQL injection protection
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com\"; DROP TABLE users; --", "password": "password"}'

# Test request size limiting
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"name": "'$(python -c "print('A' * 2000)")'", "email": "test@example.com"}'
```

## 🚀 Deployment

### Docker (Coming Soon)

### Production Configuration

For production deployment:

1. Set environment variables for sensitive data
2. Use `config.prod.yaml` configuration
3. Enable Redis for better performance
4. Configure proper CORS settings
5. Set up proper logging and monitoring

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Logrus](https://github.com/sirupsen/logrus)
- [Viper](https://github.com/spf13/viper)
- [JWT-Go](https://github.com/golang-jwt/jwt)
