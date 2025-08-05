# 🚀 Deployment Guide

This guide covers various deployment options for the Gin Boilerplate application.

## 📋 Prerequisites

- Go 1.23+
- PostgreSQL 12+
- Redis (optional)
- Reverse proxy (Nginx/Apache) for production

## 🏗️ Build for Production

### 1. Build Binary

```bash
# Build for current platform
make build

# Build for Linux (from any platform)
GOOS=linux GOARCH=amd64 go build -o gin-boilerplate-linux .

# Build for Windows (from any platform)  
GOOS=windows GOARCH=amd64 go build -o gin-boilerplate.exe .

# Build with optimizations
go build -ldflags="-s -w" -o gin-boilerplate .
```

### 2. Prepare Configuration

```bash
# Copy production config
cp config.prod.yaml.example config.prod.yaml

# Edit production settings
nano config.prod.yaml
```

## 🐳 Docker Deployment

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gin-boilerplate .

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/gin-boilerplate .
COPY --from=builder /app/config.prod.yaml ./config.prod.yaml

# Create non-root user
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8080

CMD ["./gin-boilerplate", "-env=prod"]
```

### Docker Compose

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_BOILERPLATE_POSTGRES_HOST=postgres
      - GIN_BOILERPLATE_POSTGRES_PASSWORD=secure_password
      - GIN_BOILERPLATE_SECRET_KEY=your-super-secret-jwt-key
      - GIN_BOILERPLATE_REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=gin_boilerplate
      - POSTGRES_USER=gin_user
      - POSTGRES_PASSWORD=secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass redis_password
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped

volumes:
  postgres_data:
```

### Deploy with Docker

```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f app

# Scale application
docker-compose up -d --scale app=3
```

## ☁️ Cloud Deployment

### AWS EC2

```bash
# 1. Launch EC2 instance (Ubuntu 20.04+)
# 2. Install dependencies
sudo apt update
sudo apt install -y postgresql-client redis-tools nginx

# 3. Install Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 4. Deploy application
git clone https://github.com/your-username/gin-boilerplate.git
cd gin-boilerplate
make build

# 5. Create systemd service
sudo tee /etc/systemd/system/gin-boilerplate.service > /dev/null <<EOF
[Unit]
Description=Gin Boilerplate API
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/gin-boilerplate
ExecStart=/home/ubuntu/gin-boilerplate/gin-boilerplate -env=prod
Restart=always
RestartSec=5
Environment=GIN_BOILERPLATE_POSTGRES_PASSWORD=your_password
Environment=GIN_BOILERPLATE_SECRET_KEY=your_secret_key

[Install]
WantedBy=multi-user.target
EOF

# 6. Start service
sudo systemctl daemon-reload
sudo systemctl enable gin-boilerplate
sudo systemctl start gin-boilerplate
```

### Google Cloud Platform

```bash
# 1. Build for Linux
GOOS=linux GOARCH=amd64 go build -o gin-boilerplate .

# 2. Create app.yaml for App Engine
cat > app.yaml <<EOF
runtime: go123

env_variables:
  GIN_BOILERPLATE_POSTGRES_HOST: /cloudsql/PROJECT_ID:REGION:INSTANCE_ID
  GIN_BOILERPLATE_SECRET_KEY: your-secret-key

automatic_scaling:
  min_instances: 1
  max_instances: 10
EOF

# 3. Deploy
gcloud app deploy
```

### Heroku

```bash
# 1. Create Procfile
echo "web: ./gin-boilerplate -env=prod" > Procfile

# 2. Create heroku app
heroku create your-app-name

# 3. Add PostgreSQL addon
heroku addons:create heroku-postgresql:hobby-dev

# 4. Set environment variables
heroku config:set GIN_BOILERPLATE_SECRET_KEY=your-secret-key
heroku config:set GIN_BOILERPLATE_PORT=$PORT

# 5. Deploy
git push heroku main
```

## 🔧 Nginx Configuration

### Basic Configuration

```nginx
# /etc/nginx/sites-available/gin-boilerplate
server {
    listen 80;
    server_name your-domain.com;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req zone=api burst=20 nodelay;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Static files (if any)
    location /static/ {
        alias /var/www/gin-boilerplate/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

### SSL Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL certificates
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Rest of configuration...
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

## 📊 Monitoring & Logging

### Application Logs

```bash
# View logs
journalctl -u gin-boilerplate -f

# Log rotation
sudo tee /etc/logrotate.d/gin-boilerplate > /dev/null <<EOF
/var/log/gin-boilerplate/*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 ubuntu ubuntu
    postrotate
        systemctl reload gin-boilerplate
    endscript
}
EOF
```

### Health Checks

```bash
# Simple health check script
#!/bin/bash
HEALTH_URL="http://localhost:8080/api/v1/health/check"
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $HEALTH_URL)

if [ $RESPONSE -eq 200 ]; then
    echo "✅ Application is healthy"
    exit 0
else
    echo "❌ Application is unhealthy (HTTP $RESPONSE)"
    exit 1
fi
```

## 🔒 Security Considerations

### Environment Variables

```bash
# Never commit these to version control
export GIN_BOILERPLATE_POSTGRES_PASSWORD="secure_database_password"
export GIN_BOILERPLATE_SECRET_KEY="super-secret-jwt-key-at-least-32-chars"
export GIN_BOILERPLATE_REDIS_PASSWORD="redis_password"
```

### Firewall Configuration

```bash
# UFW (Ubuntu)
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw deny 8080  # Don't expose app port directly
sudo ufw enable
```

### Database Security

```sql
-- Create dedicated user with limited privileges
CREATE USER gin_app WITH PASSWORD 'secure_password';
GRANT CONNECT ON DATABASE gin_boilerplate TO gin_app;
GRANT USAGE ON SCHEMA public TO gin_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO gin_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO gin_app;
```

## 🚨 Troubleshooting

### Common Issues

1. **Port binding error**:
```bash
# Check if port is in use
sudo netstat -tlnp | grep :8080
# Kill process if needed
sudo kill -9 <PID>
```

2. **Database connection failed**:
```bash
# Test database connection
psql -h localhost -U gin_user -d gin_boilerplate
```

3. **Permission denied**:
```bash
# Fix file permissions
chmod +x gin-boilerplate
chown ubuntu:ubuntu gin-boilerplate
```

### Performance Tuning

```yaml
# config.prod.yaml
postgres:
  maxOpenConns: 25
  maxIdleConns: 5
  connMaxLifetime: "5m"

security:
  maxBodySize: 5242880  # 5MB for production
  maxFormFields: 50     # Reduced for production
```

## 📈 Scaling

### Horizontal Scaling

```bash
# Load balancer configuration (Nginx)
upstream gin_backend {
    server 127.0.0.1:8080;
    server 127.0.0.1:8081;
    server 127.0.0.1:8082;
}

server {
    location / {
        proxy_pass http://gin_backend;
    }
}
```

### Database Scaling

- Use read replicas for read-heavy workloads
- Implement connection pooling
- Consider database sharding for very large datasets

This deployment guide should help you get your Gin boilerplate running in production! 🚀
