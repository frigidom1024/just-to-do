@echo off
REM Docker 容器运行脚本 - Windows 版本
REM 用法: run.bat [version] [port]
REM 示例: run.bat latest 8080

REM 设置 UTF-8 编码以避免中文乱码
chcp 65001 >nul 2>&1

SETLOCAL EnableDelayedExpansion

REM 切换到项目根目录
cd /d %~dp0..\..\..

REM 加载 .env 文件（如果存在）
if exist .env (
    for /f "tokens=1,2 delims==" %%a in ('type .env ^| findstr /v /c:"#" ^| findstr /v /c:"^$"') do (
        set "%%a=%%b"
    )
)

REM 配置变量（从 .env 读取，命令行参数优先）
if "%IMAGE_NAME%"=="" set IMAGE_NAME=todolist
if "%CONTAINER_NAME%"=="" set CONTAINER_NAME=todo_server

REM VERSION: 命令行参数 > .env 文件 > latest
if "%~1"=="" (
    if "%VERSION%"=="" set VERSION=latest
) else (
    set VERSION=%~1
)

REM HOST_PORT: 命令行参数 > .env 文件 > 8080
if "%~2"=="" (
    if "%HOST_PORT%"=="" set HOST_PORT=8080
) else (
    set HOST_PORT=%~2
)

if "%SERVER_PORT%"=="" set SERVER_PORT=8080
if "%NETWORK_NAME%"=="" set NETWORK_NAME=todo_network

REM 数据库配置
if "%DB_HOST%"=="" set DB_HOST=mysql
if "%DB_PORT%"=="" set DB_PORT=3306
if "%DB_USER%"=="" set DB_USER=root
if "%DB_PASSWORD%"=="" set DB_PASSWORD=123456
if "%DB_NAME%"=="" set DB_NAME=test

echo ========================================
echo Docker 容器运行脚本
echo ========================================
echo 镜像: %IMAGE_NAME%:%VERSION%
echo 容器: %CONTAINER_NAME%
echo 端口: %HOST_PORT% -^> %SERVER_PORT%
echo 网络: %NETWORK_NAME%
echo ========================================
echo.

REM 检查镜像是否存在
docker images %IMAGE_NAME%:%VERSION% --format "{{.Repository}}:{{.Tag}}" | findstr "%IMAGE_NAME%:%VERSION%" >nul
if errorlevel 1 (
    echo [ERROR] 镜像 %IMAGE_NAME%:%VERSION% 不存在
    echo 请先运行: build.bat
    exit /b 1
)

REM 检查容器是否已运行
docker ps -q -f name=%CONTAINER_NAME% | findstr . >nul
if not errorlevel 1 (
    echo [ERROR] 容器 %CONTAINER_NAME% 已在运行
    echo 如需重启，请先运行: docker stop %CONTAINER_NAME%
    exit /b 1
)

REM 检查网络是否存在
docker network ls --format "{{.Name}}" | findstr "^%NETWORK_NAME%$" >nul
if errorlevel 1 (
    echo [INFO] 创建网络: %NETWORK_NAME%
    docker network create %NETWORK_NAME%
)

REM 停止并删除旧容器（如果存在）
docker ps -a -q -f name=%CONTAINER_NAME% | findstr . >nul
if not errorlevel 1 (
    echo [INFO] 删除旧容器: %CONTAINER_NAME%
    docker rm -f %CONTAINER_NAME% >nul 2>&1
)

REM 运行容器
echo [INFO] 启动容器...
docker run -d ^
    --name %CONTAINER_NAME% ^
    --network %NETWORK_NAME% ^
    -p %HOST_PORT%:%SERVER_PORT% ^
    -e DB_HOST=%DB_HOST% ^
    -e DB_PORT=%DB_PORT% ^
    -e DB_USER=%DB_USER% ^
    -e DB_PASSWORD=%DB_PASSWORD% ^
    -e DB_NAME=%DB_NAME% ^
    -e JWT_SECRET=%JWT_SECRET% ^
    -e JWT_EXPIRE=%JWT_EXPIRE% ^
    --restart unless-stopped ^
    %IMAGE_NAME%:%VERSION%

if errorlevel 1 (
    echo [ERROR] 容器启动失败
    exit /b 1
)

REM 等待容器启动
timeout /t 2 /nobreak >nul

echo.
echo ========================================
echo 容器启动成功！
echo ========================================
echo 容器名称: %CONTAINER_NAME%
echo.

REM 显示容器日志
echo [INFO] 容器日志:
docker logs %CONTAINER_NAME% 2^>^&1

echo.
echo 查看实时日志:
echo   docker logs -f %CONTAINER_NAME%
echo.
echo 访问服务:
echo   健康检查: http://localhost:%HOST_PORT%/health
echo   API: http://localhost:%HOST_PORT%/api/v1/...
echo.
echo 停止容器:
echo   docker stop %CONTAINER_NAME%
echo.
echo 删除容器:
echo   docker rm -f %CONTAINER_NAME%
echo ========================================

ENDLOCAL
