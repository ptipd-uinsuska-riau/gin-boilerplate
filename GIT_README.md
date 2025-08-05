# 🔥 Gin Boilerplate - Production Ready

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-1.10+-00ADD8?style=for-the-badge&logo=gin)](https://gin-gonic.com/)
[![Security](https://img.shields.io/badge/Security-Enhanced-green?style=for-the-badge&logo=shield)](https://github.com/your-username/gin-boilerplate)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)

> **Enterprise-ready Go Gin boilerplate dengan fitur security yang komprehensif, clean architecture, dan production-ready configuration.**

## ✨ Highlights

🛡️ **Advanced Security** - XSS & SQL Injection protection, Request size limiting  
🏗️ **Clean Architecture** - Modular design dengan dependency injection  
🔐 **JWT Authentication** - Secure authentication system  
📊 **Auto Documentation** - Swagger/OpenAPI integration  
🚀 **Production Ready** - Environment-specific configurations  
🧪 **Well Tested** - Comprehensive test suite

## 🚀 Quick Start

```bash
# 1. Clone repository
git clone https://github.com/your-username/gin-boilerplate.git
cd gin-boilerplate

# 2. Install dependencies
make deps

# 3. Setup database
createdb gin_boilerplate

# 4. Configure environment
cp config.local.yaml.example config.local.yaml
# Edit database credentials in config.local.yaml

# 5. Run application
make run
```

🎉 **Server running at** `http://localhost:8080`  
📖 **API Documentation** `http://localhost:8080/swagger/index.html`

## 🔒 Security Features

### 🛡️ Input Sanitization

- **XSS Protection**: Automatic detection and blocking of malicious scripts
- **SQL Injection Protection**: Pattern-based detection of injection attempts
- **HTML Sanitization**: Safe handling of user input with configurable strictness

### 📏 Request Limiting

- **Body Size Limiting**: Configurable request size limits (default: 10MB)
- **Form Protection**: Limits on form fields, files, and multipart data
- **DoS Prevention**: URL length and query parameter restrictions

### ⚙️ Configurable Security

```yaml
security:
  enableSanitization: true
  enableXSSDetection: true
  enableSQLInjectionDetection: true
  strictMode: true
  maxBodySize: 10485760 # 10MB
  maxFileSize: 52428800 # 50MB
```

## 🏗️ Architecture

```
gin-boilerplate/
├── 🚀 boot/                   # Dependency injection & setup
├── ⚙️ infrastructure/         # Infrastructure layer
│   ├── 🔧 config/            # Configuration management
│   ├── 🗄️ database/          # Database connections
│   ├── 🛡️ middleware/        # Security & other middlewares
│   ├── 🧹 sanitizer/         # Input sanitization utilities
│   ├── 🔐 jwt/               # JWT authentication
│   └── 📝 validator/         # Request validation
├── 📦 modules/               # Business logic modules
│   ├── 🔑 auth/              # Authentication module
│   ├── 👤 user/              # User management
│   └── 💚 health/            # Health checks
└── 🛣️ router/                # Route definitions
```

## 📚 API Endpoints

### 🔑 Authentication

```http
POST   /api/v1/auth/register     # Register new user
POST   /api/v1/auth/login        # User login
GET    /api/v1/auth/profile      # Get user profile
PUT    /api/v1/auth/profile      # Update profile
POST   /api/v1/auth/change-password  # Change password
```

### 👥 User Management

```http
GET    /api/v1/users             # List users
GET    /api/v1/users/:id         # Get user by ID
POST   /api/v1/users             # Create user
PUT    /api/v1/users/:id         # Update user
DELETE /api/v1/users/:id         # Delete user
```

### 💚 Health Checks

```http
GET    /api/v1/health/ping       # Simple ping
GET    /api/v1/health/check      # Comprehensive health check
```

## 🧪 Testing Security

```bash
# Test XSS protection
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "<script>alert(\"xss\")</script>", "email": "test@example.com"}'
# Expected: 400 Bad Request - Blocked!

# Test normal request
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "password123"}'
# Expected: 201 Created - Success!
```

## 🌍 Environment Support

| Environment    | Security | Limits    | Use Case        |
| -------------- | -------- | --------- | --------------- |
| **Local**      | Balanced | 10MB body | Development     |
| **Production** | Strict   | 5MB body  | Live deployment |
| **Testing**    | Relaxed  | 20MB body | Automated tests |

## 🔧 Available Commands

```bash
make build          # Build application
make run            # Run in development
make test           # Run test suite
make test-coverage  # Test with coverage report
make swagger        # Generate API docs
make fmt            # Format code
make lint           # Lint code
make clean          # Clean build artifacts
```

## 📦 Tech Stack

- **Framework**: [Gin](https://gin-gonic.com/) - High-performance HTTP web framework
- **Database**: [PostgreSQL](https://www.postgresql.org/) with [GORM](https://gorm.io/) ORM
- **Cache**: [Redis](https://redis.io/) (optional)
- **Authentication**: [JWT](https://github.com/golang-jwt/jwt) tokens
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **Documentation**: [Swagger](https://swagger.io/) auto-generation
- **Configuration**: [Viper](https://github.com/spf13/viper) multi-format config
- **Logging**: [Logrus](https://github.com/sirupsen/logrus) structured logging

## 🚀 Production Deployment

### Environment Variables

```bash
export GIN_BOILERPLATE_POSTGRES_PASSWORD="your-secure-password"
export GIN_BOILERPLATE_SECRET_KEY="your-jwt-secret-key"
export GIN_BOILERPLATE_REDIS_PASSWORD="redis-password"
```

### Build & Deploy

```bash
# Build for production
make build-prod

# Run with production config
./gin-boilerplate -env=prod
```

## 📖 Documentation

- 📋 **[Usage Guide](USAGE_GUIDE.md)** - Comprehensive setup and usage
- 🔒 **[Security Examples](SECURITY_EXAMPLES.md)** - Security feature examples
- 📚 **[API Documentation](http://localhost:8080/swagger/index.html)** - Interactive API docs

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⭐ Show Your Support

If this project helped you, please give it a ⭐ star!

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Logrus](https://github.com/sirupsen/logrus)
- [Viper](https://github.com/spf13/viper)
- [JWT-Go](https://github.com/golang-jwt/jwt)

---

<div align="center">

**[⬆ Back to Top](#-gin-boilerplate---production-ready)**

Made with ❤️ for the Go community

</div>
