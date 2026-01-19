@echo off
REM Docker 构建并运行脚本 - Windows 版本
REM 用法: build-and-run.bat [version] [port]

REM 设置 UTF-8 编码以避免中文乱码
chcp 65001 >nul 2>&1

SETLOCAL EnableDelayedExpansion

set VERSION=%1
set PORT=%2

if "%VERSION%"=="" set VERSION=latest
if "%PORT%"=="" set PORT=8080

echo ========================================
echo Docker 构建并运行脚本
echo ========================================
echo 版本: %VERSION%
echo 端口: %PORT%
echo ========================================

REM 1. 构建镜像
echo.
echo [步骤 1/2] 构建镜像...
echo.
call "%~dp0build.bat" %VERSION%

if errorlevel 1 (
    echo [ERROR] 构建失败
    exit /b 1
)

REM 2. 运行容器
echo.
echo [步骤 2/2] 运行容器...
echo.
call "%~dp0run.bat" %VERSION% %PORT%

if errorlevel 1 (
    echo [ERROR] 运行失败
    exit /b 1
)

ENDLOCAL
