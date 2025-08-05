# Security Features Examples

This document provides examples of how to use the security features implemented in this Gin boilerplate.

## Input Sanitization Examples

### XSS Protection

The application automatically detects and blocks XSS attempts:

```bash
# This request will be blocked in strict mode
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "<script>alert(\"xss\")</script>John",
    "email": "john@example.com",
    "password": "password123"
  }'

# Response: 400 Bad Request
# {"success":false,"code":0,"message":"Request contains potentially malicious content"}
```

### SQL Injection Protection

SQL injection attempts are detected and blocked:

```bash
# This request will be blocked
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com\"; DROP TABLE users; --",
    "password": "password"
  }'

# Response: 400 Bad Request
# {"success":false,"code":0,"message":"Request contains potentially malicious content"}
```

### Safe Requests

Normal requests pass through without issues:

```bash
# This request will succeed
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'

# Response: 201 Created
# {"success":true,"code":1,"message":"User registered successfully","data":{...}}
```

## Request Size Limiting Examples

### Body Size Limiting

Large request bodies are rejected:

```bash
# Generate a large payload (exceeds 10MB default limit)
python3 -c "print('A' * 11000000)" > large_payload.txt

# This request will be blocked
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d @large_payload.txt

# Response: 400 Bad Request
# {"success":false,"code":0,"message":"Request body too large. Maximum allowed: 10485760 bytes"}
```

### URL Length Limiting

Extremely long URLs are rejected:

```bash
# This request will be blocked (URL too long)
curl "http://localhost:8080/api/v1/users?$(python3 -c "print('param=' + 'A' * 3000)")"

# Response: 400 Bad Request
# {"success":false,"code":0,"message":"URL too long. Maximum allowed: 2048 characters"}
```

### Query Parameter Limiting

Too many query parameters are rejected:

```bash
# Generate URL with many parameters
PARAMS=$(python3 -c "print('&'.join([f'param{i}=value{i}' for i in range(60)]))")
curl "http://localhost:8080/api/v1/users?$PARAMS"

# Response: 400 Bad Request
# {"success":false,"code":0,"message":"Too many query parameters. Maximum allowed: 50"}
```

## File Upload Security

### File Size Limiting

Large file uploads are rejected:

```bash
# Create a large file
dd if=/dev/zero of=large_file.bin bs=1M count=60

# This upload will be blocked (exceeds 50MB default limit)
curl -X POST http://localhost:8080/api/v1/upload/file \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@large_file.bin"

# Response: 400 Bad Request
# {"success":false,"code":0,"message":"File too large. Maximum allowed: 52428800 bytes"}
```

### Safe File Upload

Normal file uploads work fine:

```bash
# Create a small file
echo "Hello World" > small_file.txt

# This upload will succeed
curl -X POST http://localhost:8080/api/v1/upload/file \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@small_file.txt"

# Response: 200 OK
# {"success":true,"code":1,"message":"File upload endpoint","data":{...}}
```

## Configuration Examples

### Development Configuration

For development, you might want more lenient settings:

```yaml
# config.dev.yaml
security:
  enableSanitization: true
  enableXSSDetection: true
  enableSQLInjectionDetection: true
  maxStringLength: 2000
  strictMode: false  # Allow sanitized content through
  maxBodySize: 20971520      # 20MB
  maxFormFields: 200
  maxFormMemory: 67108864    # 64MB
  maxFormFiles: 20
  maxHeaderSize: 2097152     # 2MB
  maxURLLength: 4096
  maxQueryParams: 100
  maxFileSize: 104857600     # 100MB
  maxResponseSize: 209715200 # 200MB
```

### Production Configuration

For production, use more restrictive settings:

```yaml
# config.prod.yaml
security:
  enableSanitization: true
  enableXSSDetection: true
  enableSQLInjectionDetection: true
  maxStringLength: 500       # More restrictive
  strictMode: true           # Block malicious content
  maxBodySize: 5242880       # 5MB
  maxFormFields: 50
  maxFormMemory: 16777216    # 16MB
  maxFormFiles: 5
  maxHeaderSize: 524288      # 512KB
  maxURLLength: 1024
  maxQueryParams: 25
  maxFileSize: 10485760      # 10MB
  maxResponseSize: 52428800  # 50MB
```

### Testing Configuration

For testing, you might want to disable some features:

```yaml
# config.test.yaml
security:
  enableSanitization: false  # Disabled for easier testing
  enableXSSDetection: false
  enableSQLInjectionDetection: false
  maxStringLength: 5000      # Very lenient
  strictMode: false
  maxBodySize: 52428800      # 50MB
  maxFormFields: 500
  maxFormMemory: 134217728   # 128MB
  maxFormFiles: 50
  maxHeaderSize: 5242880     # 5MB
  maxURLLength: 8192
  maxQueryParams: 200
  maxFileSize: 209715200     # 200MB
  maxResponseSize: 419430400 # 400MB
```

## Custom Middleware Usage

You can also use the security middlewares individually:

```go
package main

import (
    "gin-boilerplate/infrastructure/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.New()

    // Use only input sanitization
    sanitizerConfig := middleware.SanitizerConfig{
        EnableXSSDetection:          true,
        EnableSQLInjectionDetection: true,
        MaxStringLength:             1000,
        StrictMode:                  true,
        SkipFields:                  []string{"password", "token"},
    }
    r.Use(middleware.SanitizerMiddleware(sanitizerConfig))

    // Use only request size limiting
    limiterConfig := middleware.RequestLimiterConfig{
        MaxBodySize:    10 * 1024 * 1024, // 10MB
        MaxFormFields:  100,
        MaxURLLength:   2048,
        MaxQueryParams: 50,
    }
    r.Use(middleware.RequestLimiterMiddleware(limiterConfig))

    // Use only file size limiting for specific routes
    uploadGroup := r.Group("/upload")
    uploadGroup.Use(middleware.FileSizeLimiterMiddleware(50 * 1024 * 1024)) // 50MB
    {
        uploadGroup.POST("/file", handleFileUpload)
    }

    r.Run(":8080")
}
```

## Testing Security Features

You can test the security features using the provided test suite:

```bash
# Test input sanitization
go test ./infrastructure/sanitizer -v

# Test security middlewares
go test ./infrastructure/middleware -v -run TestSecurity

# Test all security features
go test ./infrastructure/sanitizer ./infrastructure/middleware -v
```

## Monitoring and Logging

Security events are automatically logged:

```
time="2025-08-05T11:42:08+07:00" level=warning msg="Client error:Request contains potentially malicious content"
time="2025-08-05T11:42:08+07:00" level=warning msg="Client error:Request body too large. Maximum allowed: 10485760 bytes"
```

You can monitor these logs to detect potential attacks and adjust your security settings accordingly.
