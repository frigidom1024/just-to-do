@echo off
REM Todo Project Startup Script for Windows

setlocal enabledelayedexpansion

echo =========================================
echo   Todo Project - Development Setup
echo =========================================
echo.

REM Check if Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker is not running. Please start Docker and try again.
    exit /b 1
)
echo [OK] Docker is running

REM Check configuration file
if not exist "config\config.yaml" (
    echo [WARN] Configuration file not found. Creating from example...
    if not exist "config" mkdir config
    copy config\config.example.yaml config\config.yaml >nul
    echo [OK] Created config\config.yaml
    echo [WARN] Please review and modify the configuration if needed
) else (
    echo [OK] Configuration file exists
)

REM Start database
echo.
echo Starting MySQL database...
docker-compose up -d mysql

echo [WAIT] Waiting for database to be ready...
set /a count=0
:wait_loop
if !count! geq 30 (
    echo [ERROR] Database failed to start
    exit /b 1
)

docker-compose exec -T mysql mysqladmin ping -h localhost -u root -prootpassword --silent >nul 2>&1
if errorlevel 1 (
    echo|set /p="."
    set /a count+=1
    timeout /t 1 >nul
    goto wait_loop
)

echo.
echo [OK] Database is ready

REM Show database information
echo.
echo =========================================
echo   Database Information
echo =========================================
echo Host: 127.0.0.1
echo Port: 3306
echo Database: test
echo User: root
echo Password: 123456
echo.
echo Test Accounts:
echo   Admin:     admin / 123456
echo   Test User: test_user / 123456
echo =========================================

echo.
echo =========================================
echo   Setup Complete!
echo =========================================
echo.
echo You can now:
echo   1. Run the application:
echo      go run cmd/server/main.go
echo.
echo   2. View database logs:
echo      docker-compose logs -f mysql
echo.
echo   3. Connect to database:
echo      docker-compose exec mysql mysql -u root -p123456 test
echo.
echo   4. Stop database:
echo      docker-compose down
echo.

endlocal
