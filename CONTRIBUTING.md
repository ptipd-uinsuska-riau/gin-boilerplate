# Contributing to Gin Boilerplate

Thank you for your interest in contributing to this project! This document provides guidelines and information for contributors.

## 🤝 How to Contribute

### 1. Fork & Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/your-username/gin-boilerplate.git
cd gin-boilerplate

# Add upstream remote
git remote add upstream https://github.com/original-owner/gin-boilerplate.git
```

### 2. Set Up Development Environment

```bash
# Install dependencies
make deps

# Set up database
createdb gin_boilerplate_dev

# Copy and configure environment
cp config.local.yaml.example config.local.yaml
# Edit config.local.yaml with your database credentials

# Run tests to ensure everything works
make test
```

### 3. Create Feature Branch

```bash
# Create and switch to feature branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/bug-description
```

## 📝 Development Guidelines

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting: `make fmt`
- Run linter: `make lint`
- Write meaningful commit messages

### Testing

- Write tests for new features
- Ensure all tests pass: `make test`
- Maintain or improve test coverage: `make test-coverage`
- Test security features if modifying security code

### Security Considerations

When contributing to security features:

1. **Test thoroughly** - Security bugs can be critical
2. **Document changes** - Update security documentation
3. **Consider edge cases** - Think about bypass attempts
4. **Performance impact** - Security shouldn't break performance

### Documentation

- Update README.md if adding new features
- Add examples to SECURITY_EXAMPLES.md for security features
- Update API documentation comments for Swagger
- Write clear commit messages

## 🔧 Development Commands

```bash
# Development workflow
make run            # Run application
make test           # Run tests
make test-coverage  # Test with coverage
make fmt            # Format code
make lint           # Lint code
make swagger        # Generate API docs

# Build commands
make build          # Build for current OS
make build-prod     # Build for production
make clean          # Clean build artifacts
```

## 🧪 Testing Guidelines

### Unit Tests

```bash
# Run specific package tests
go test ./infrastructure/sanitizer -v
go test ./infrastructure/middleware -v

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Run integration tests
make test

# Test specific functionality
go test ./modules/auth -v
go test ./modules/user -v
```

### Security Tests

```bash
# Test security features
go test ./infrastructure/sanitizer ./infrastructure/middleware -v

# Manual security testing
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "<script>alert(\"xss\")</script>", "email": "test@example.com"}'
```

## 📋 Pull Request Process

### Before Submitting

1. **Sync with upstream**:
```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

2. **Rebase your feature branch**:
```bash
git checkout feature/your-feature-name
git rebase main
```

3. **Run all checks**:
```bash
make test
make fmt
make lint
```

### PR Requirements

- [ ] All tests pass
- [ ] Code is formatted (`make fmt`)
- [ ] No linting errors (`make lint`)
- [ ] Documentation updated if needed
- [ ] Security implications considered
- [ ] Commit messages are clear

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update
- [ ] Security enhancement

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing completed

## Security Impact
- [ ] No security implications
- [ ] Security feature enhancement
- [ ] Potential security impact (explain)

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added for new functionality
```

## 🐛 Bug Reports

### Before Reporting

1. Check existing issues
2. Ensure you're using the latest version
3. Test with minimal reproduction case

### Bug Report Template

```markdown
**Bug Description**
Clear description of the bug

**To Reproduce**
Steps to reproduce:
1. Go to '...'
2. Click on '....'
3. See error

**Expected Behavior**
What you expected to happen

**Environment**
- OS: [e.g. Ubuntu 20.04]
- Go version: [e.g. 1.23.1]
- Application version: [e.g. v1.0.0]

**Additional Context**
Any other context about the problem
```

## 💡 Feature Requests

### Feature Request Template

```markdown
**Feature Description**
Clear description of the feature

**Use Case**
Why is this feature needed?

**Proposed Solution**
How should this feature work?

**Alternatives Considered**
Other solutions you've considered

**Additional Context**
Any other context or screenshots
```

## 🔒 Security Issues

**DO NOT** open public issues for security vulnerabilities.

Instead:
1. Email security issues to: [security@yourproject.com]
2. Include detailed description
3. Provide reproduction steps
4. Allow time for fix before disclosure

## 📚 Resources

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)

## 🏷️ Commit Message Guidelines

### Format

```
type(scope): subject

body

footer
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```bash
feat(auth): add password reset functionality

Add password reset endpoint with email verification.
Includes rate limiting and security validations.

Closes #123

fix(security): prevent XSS in user input

Improve input sanitization to handle edge cases
in script tag detection.

docs(readme): update installation instructions

Add Docker setup instructions and troubleshooting guide.
```

## 🎉 Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes for significant contributions
- Special thanks for security improvements

Thank you for contributing! 🙏
