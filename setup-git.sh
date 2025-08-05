#!/bin/bash

# Gin Boilerplate Git Setup Script
# This script helps you set up the repository for Git

set -e

echo "🚀 Setting up Gin Boilerplate for Git..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

# Check if git is installed
if ! command -v git &> /dev/null; then
    print_error "Git is not installed. Please install Git first."
    exit 1
fi

print_status "Git is installed"

# Check if we're already in a git repository
if [ -d ".git" ]; then
    print_warning "Already in a Git repository"
    read -p "Do you want to continue? This will modify the existing repository. (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Setup cancelled"
        exit 0
    fi
else
    # Initialize git repository
    print_info "Initializing Git repository..."
    git init
    print_status "Git repository initialized"
fi

# Create .gitignore if it doesn't exist
if [ ! -f ".gitignore" ]; then
    print_info "Creating .gitignore..."
    cat > .gitignore << 'EOF'
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# Build output
gin-boilerplate
gin-boilerplate.exe
gin-boilerplate_unix

# Log files
*.log
gin-boilerplate.log

# Environment configuration files (keep examples)
config.local.yaml
config.dev.yaml
config.uat.yaml
config.prod.yaml
config.test.yaml

# Keep example files
!config.local.yaml.example
!config.dev.yaml.example
!config.prod.yaml.example
!config.test.yaml.example

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Coverage reports
coverage.out
coverage.html

# Temporary files
tmp/
temp/

# Database files
*.db
*.sqlite
*.sqlite3

# Redis dump
dump.rdb

# Air live reload
tmp/

# Environment variables
.env
.env.local
.env.development
.env.test
.env.production

# Backup files
*.bak
*.backup

# Debug files
debug
debug.test

# Profiling files
*.prof
*.pprof
EOF
    print_status ".gitignore created"
else
    print_status ".gitignore already exists"
fi

# Copy the Git README to replace the current README
if [ -f "GIT_README.md" ]; then
    print_info "Updating README.md for Git repository..."
    cp GIT_README.md README.md
    print_status "README.md updated"
fi

# Add all files to git
print_info "Adding files to Git..."
git add .

# Check if there are any changes to commit
if git diff --staged --quiet; then
    print_warning "No changes to commit"
else
    # Commit initial files
    print_info "Creating initial commit..."
    git commit -m "feat: initial commit with security-enhanced gin boilerplate

- Add comprehensive input sanitization (XSS & SQL injection protection)
- Add request size limiting middleware
- Add JWT authentication system
- Add user management module
- Add health check endpoints
- Add Swagger documentation
- Add environment-specific configurations
- Add comprehensive test suite
- Add production-ready security features"
    print_status "Initial commit created"
fi

# Ask for remote repository URL
echo
print_info "Git repository setup complete!"
echo
echo "Next steps:"
echo "1. Create a new repository on GitHub/GitLab/etc."
echo "2. Add the remote repository:"
echo "   ${BLUE}git remote add origin https://github.com/yourusername/your-repo.git${NC}"
echo "3. Push to remote repository:"
echo "   ${BLUE}git branch -M main${NC}"
echo "   ${BLUE}git push -u origin main${NC}"
echo

read -p "Do you want to add a remote repository now? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "Enter the remote repository URL: " remote_url
    if [ ! -z "$remote_url" ]; then
        git remote add origin "$remote_url"
        print_status "Remote repository added: $remote_url"
        
        read -p "Do you want to push to the remote repository now? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            git branch -M main
            git push -u origin main
            print_status "Repository pushed to remote"
        fi
    fi
fi

echo
print_status "Setup complete! Your Gin boilerplate is ready for development."
echo
echo "Quick start commands:"
echo "  ${BLUE}make deps${NC}     # Install dependencies"
echo "  ${BLUE}make run${NC}      # Run the application"
echo "  ${BLUE}make test${NC}     # Run tests"
echo
echo "Documentation:"
echo "  📖 README.md - Project overview"
echo "  📋 USAGE_GUIDE.md - Detailed usage guide"
echo "  🔒 SECURITY_EXAMPLES.md - Security features examples"
echo "  🤝 CONTRIBUTING.md - Contributing guidelines"
echo
