@echo off
setlocal enabledelayedexpansion

echo 🚀 Setting up Gin Boilerplate for Git...
echo.

REM Check if git is installed
git --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Git is not installed. Please install Git first.
    pause
    exit /b 1
)

echo ✅ Git is installed

REM Check if we're already in a git repository
if exist ".git" (
    echo ⚠️ Already in a Git repository
    set /p continue="Do you want to continue? This will modify the existing repository. (y/N): "
    if /i not "!continue!"=="y" (
        echo ℹ️ Setup cancelled
        pause
        exit /b 0
    )
) else (
    echo ℹ️ Initializing Git repository...
    git init
    echo ✅ Git repository initialized
)

REM Create .gitignore if it doesn't exist
if not exist ".gitignore" (
    echo ℹ️ Creating .gitignore...
    (
        echo # Binaries for programs and plugins
        echo *.exe
        echo *.exe~
        echo *.dll
        echo *.so
        echo *.dylib
        echo.
        echo # Test binary, built with `go test -c`
        echo *.test
        echo.
        echo # Output of the go coverage tool
        echo *.out
        echo.
        echo # Go workspace file
        echo go.work
        echo.
        echo # Build output
        echo gin-boilerplate
        echo gin-boilerplate.exe
        echo gin-boilerplate_unix
        echo.
        echo # Log files
        echo *.log
        echo gin-boilerplate.log
        echo.
        echo # Environment configuration files ^(keep examples^)
        echo config.local.yaml
        echo config.dev.yaml
        echo config.uat.yaml
        echo config.prod.yaml
        echo config.test.yaml
        echo.
        echo # Keep example files
        echo !config.local.yaml.example
        echo !config.dev.yaml.example
        echo !config.prod.yaml.example
        echo !config.test.yaml.example
        echo.
        echo # IDE files
        echo .vscode/
        echo .idea/
        echo *.swp
        echo *.swo
        echo *~
        echo.
        echo # OS generated files
        echo .DS_Store
        echo .DS_Store?
        echo ._*
        echo .Spotlight-V100
        echo .Trashes
        echo ehthumbs.db
        echo Thumbs.db
        echo.
        echo # Coverage reports
        echo coverage.out
        echo coverage.html
        echo.
        echo # Temporary files
        echo tmp/
        echo temp/
        echo.
        echo # Database files
        echo *.db
        echo *.sqlite
        echo *.sqlite3
        echo.
        echo # Redis dump
        echo dump.rdb
        echo.
        echo # Air live reload
        echo tmp/
        echo.
        echo # Environment variables
        echo .env
        echo .env.local
        echo .env.development
        echo .env.test
        echo .env.production
        echo.
        echo # Backup files
        echo *.bak
        echo *.backup
        echo.
        echo # Debug files
        echo debug
        echo debug.test
        echo.
        echo # Profiling files
        echo *.prof
        echo *.pprof
    ) > .gitignore
    echo ✅ .gitignore created
) else (
    echo ✅ .gitignore already exists
)

REM Copy the Git README to replace the current README
if exist "GIT_README.md" (
    echo ℹ️ Updating README.md for Git repository...
    copy GIT_README.md README.md >nul
    echo ✅ README.md updated
)

REM Add all files to git
echo ℹ️ Adding files to Git...
git add .

REM Check if there are any changes to commit
git diff --staged --quiet
if errorlevel 1 (
    echo ℹ️ Creating initial commit...
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
    echo ✅ Initial commit created
) else (
    echo ⚠️ No changes to commit
)

echo.
echo ✅ Git repository setup complete!
echo.
echo Next steps:
echo 1. Create a new repository on GitHub/GitLab/etc.
echo 2. Add the remote repository:
echo    git remote add origin https://github.com/yourusername/your-repo.git
echo 3. Push to remote repository:
echo    git branch -M main
echo    git push -u origin main
echo.

set /p add_remote="Do you want to add a remote repository now? (y/N): "
if /i "!add_remote!"=="y" (
    set /p remote_url="Enter the remote repository URL: "
    if not "!remote_url!"=="" (
        git remote add origin "!remote_url!"
        echo ✅ Remote repository added: !remote_url!
        
        set /p push_now="Do you want to push to the remote repository now? (y/N): "
        if /i "!push_now!"=="y" (
            git branch -M main
            git push -u origin main
            echo ✅ Repository pushed to remote
        )
    )
)

echo.
echo ✅ Setup complete! Your Gin boilerplate is ready for development.
echo.
echo Quick start commands:
echo   make deps     # Install dependencies
echo   make run      # Run the application  
echo   make test     # Run tests
echo.
echo Documentation:
echo   📖 README.md - Project overview
echo   📋 USAGE_GUIDE.md - Detailed usage guide
echo   🔒 SECURITY_EXAMPLES.md - Security features examples
echo   🤝 CONTRIBUTING.md - Contributing guidelines
echo.
pause
