# 🚀 Gin Boilerplate Usage Guide

Panduan lengkap untuk menggunakan boilerplate Gin dengan fitur security yang telah ditingkatkan.

## 📋 Prerequisites

Pastikan Anda telah menginstall:
- **Go 1.23+** - [Download Go](https://golang.org/dl/)
- **PostgreSQL 12+** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Redis** (optional) - [Download Redis](https://redis.io/download)
- **Git** - [Download Git](https://git-scm.com/downloads)

## 🔧 Quick Start

### 1. Clone Repository

```bash
# Clone dari repository Anda
git clone https://github.com/your-username/gin-boilerplate.git
cd gin-boilerplate

# Atau gunakan template
git clone https://github.com/your-username/gin-boilerplate.git my-new-project
cd my-new-project
```

### 2. Install Dependencies

```bash
# Download dan install dependencies
make deps

# Atau manual
go mod download
go mod tidy
```

### 3. Setup Database

```sql
-- Buat database PostgreSQL
CREATE DATABASE gin_boilerplate;
CREATE USER gin_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE gin_boilerplate TO gin_user;
```

### 4. Configure Environment

```bash
# Copy configuration file
cp config.local.yaml.example config.local.yaml

# Edit konfigurasi database
nano config.local.yaml
```

Update konfigurasi database:

```yaml
postgres:
  host: "localhost"
  port: 5432
  user: "gin_user"
  password: "your_password"
  dbname: "gin_boilerplate"
  sslmode: "disable"
  timezone: "UTC"
```

### 5. Run Application

```bash
# Development mode
make run

# Atau dengan environment spesifik
go run main.go -env=local

# Build dan run
make build
./gin-boilerplate
```

Server akan berjalan di `http://localhost:8080`

## 🌍 Environment Configuration

### Local Development

```yaml
# config.local.yaml
env: "local"
port: 8080
logLevel: "DEBUG"
logMode: true

security:
  enableSanitization: true
  strictMode: false  # Lebih lenient untuk development
  maxBodySize: 10485760  # 10MB
```

### Production

```yaml
# config.prod.yaml
env: "prod"
port: 8080
logLevel: "INFO"
logMode: true

security:
  enableSanitization: true
  strictMode: true   # Strict untuk production
  maxBodySize: 5242880  # 5MB (lebih ketat)
```

### Testing

```yaml
# config.test.yaml
env: "test"
port: 8081
logLevel: "ERROR"

security:
  enableSanitization: false  # Disabled untuk testing
  strictMode: false
```

## 🔒 Security Features

### Input Sanitization

Otomatis melindungi dari:
- **XSS Attacks**: `<script>`, event handlers, javascript: protocol
- **SQL Injection**: UNION, DROP TABLE, INSERT, dll
- **HTML Injection**: Tag HTML berbahaya

### Request Size Limiting

Melindungi dari DoS attacks:
- **Body Size**: Default 10MB (configurable)
- **Form Fields**: Max 100 fields
- **File Upload**: Max 50MB per file
- **URL Length**: Max 2048 characters
- **Query Params**: Max 50 parameters

## 📚 API Documentation

Setelah aplikasi berjalan, akses dokumentasi di:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **API Base URL**: `http://localhost:8080/api/v1`

### Authentication Endpoints

```bash
# Register user baru
POST /api/v1/auth/register
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}

# Login
POST /api/v1/auth/login
{
  "email": "john@example.com",
  "password": "password123"
}

# Get profile (requires auth)
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

### User Management

```bash
# Get all users
GET /api/v1/users
Authorization: Bearer <token>

# Create user
POST /api/v1/users
Authorization: Bearer <token>
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "password123"
}
```

### Health Checks

```bash
# Simple ping
GET /api/v1/health/ping

# Comprehensive health check
GET /api/v1/health/check
```

## 🧪 Testing

### Run Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Security tests only
go test ./infrastructure/sanitizer ./infrastructure/middleware -v

# Specific test
go test ./infrastructure/sanitizer -v -run TestSanitizeString
```

### Test Security Features

```bash
# Test XSS protection
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "<script>alert(\"xss\")</script>", "email": "test@example.com", "password": "pass123"}'

# Expected: 400 Bad Request (blocked)

# Test normal request
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "password123"}'

# Expected: 201 Created (success)
```

## 🚀 Deployment

### Docker (Coming Soon)

```bash
# Build image
docker build -t gin-boilerplate .

# Run container
docker run -p 8080:8080 gin-boilerplate
```

### Production Deployment

1. **Build for production**:
```bash
make build-prod
```

2. **Set environment variables**:
```bash
export GIN_BOILERPLATE_POSTGRES_PASSWORD="secure_password"
export GIN_BOILERPLATE_SECRET_KEY="your-super-secret-jwt-key"
export GIN_BOILERPLATE_REDIS_PASSWORD="redis_password"
```

3. **Run with production config**:
```bash
./gin-boilerplate -env=prod
```

## 🔧 Development

### Project Structure

```
gin-boilerplate/
├── boot/                   # Dependency injection
├── infrastructure/         # Infrastructure layer
│   ├── config/            # Configuration
│   ├── database/          # Database connection
│   ├── middleware/        # Security & other middlewares
│   ├── sanitizer/         # Input sanitization
│   └── ...
├── modules/               # Business logic
│   ├── auth/              # Authentication
│   ├── user/              # User management
│   └── ...
├── router/                # Route definitions
└── main.go               # Entry point
```

### Adding New Features

1. **Create new module**:
```bash
mkdir modules/your-feature
touch modules/your-feature/{model,repository,service,handler}.go
```

2. **Add to dependency injection** (`boot/boot.go`):
```go
// Initialize repository
yourRepository := yourfeature.NewRepository(db.DbConn)

// Initialize service
yourService := yourfeature.NewService(yourRepository)

// Initialize handler
yourHttp := yourfeature.NewHttp(yourService, authMiddleware)
```

3. **Add routes** (`router/router.go`):
```go
// Your feature routes
yourGroup := v1.Group("/your-feature")
hr.Setup.YourHttp.GroupYourFeature(yourGroup)
```

### Available Make Commands

```bash
make build          # Build application
make run            # Run application
make test           # Run tests
make test-coverage  # Run tests with coverage
make fmt            # Format code
make lint           # Lint code
make swagger        # Generate swagger docs
make clean          # Clean build files
make deps           # Download dependencies
make dev            # Run with live reload (requires air)
```

## 🔍 Troubleshooting

### Common Issues

1. **Database connection failed**:
   - Pastikan PostgreSQL berjalan
   - Check credentials di config file
   - Pastikan database sudah dibuat

2. **Port already in use**:
   - Change port di config file
   - Kill process: `lsof -ti:8080 | xargs kill -9`

3. **Security middleware blocking requests**:
   - Check `strictMode` setting
   - Disable sanitization untuk testing
   - Check logs untuk detail error

### Debug Mode

```bash
# Run dengan debug logging
GIN_BOILERPLATE_LOG_LEVEL=DEBUG go run main.go

# Check logs
tail -f gin-boilerplate.log
```

## 📝 Contributing

1. Fork repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push branch: `git push origin feature/amazing-feature`
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file.

## 🙏 Support

- 📖 [Documentation](README.md)
- 🔒 [Security Examples](SECURITY_EXAMPLES.md)
- 🐛 [Report Issues](https://github.com/your-username/gin-boilerplate/issues)
- 💬 [Discussions](https://github.com/your-username/gin-boilerplate/discussions)
