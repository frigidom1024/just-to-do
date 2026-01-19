@echo off
REM Docker 镜像构建脚本 - Windows 版本
REM 用法: build.bat [version]
REM 示例: build.bat 1.0.0

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

REM VERSION: 命令行参数 > .env 文件 > latest
if "%~1"=="" (
    REM 没有命令行参数，使用 .env 中的值或默认值
    if "%VERSION%"=="" set VERSION=latest
) else (
    REM 有命令行参数，使用它
    set VERSION=%~1
)

set "REGISTRY=%REGISTRY%"

echo ========================================
echo Docker 镜像构建脚本
echo ========================================
echo 镜像名称: %IMAGE_NAME%
echo 版本标签: %VERSION%
echo ========================================
echo.

REM 检查 Docker 是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker 未运行，请先启动 Docker Desktop
    exit /b 1
)

REM 构建镜像
echo [INFO] 开始构建 Docker 镜像...
echo.

docker build -t %IMAGE_NAME%:%VERSION% .

if errorlevel 1 (
    echo.
    echo [ERROR] Docker 镜像构建失败
    exit /b 1
)

REM 如果版本不是 latest，额外打一个 latest 标签
if not "%VERSION%"=="latest" (
    echo [INFO] 添加 latest 标签...
    docker tag %IMAGE_NAME%:%VERSION% %IMAGE_NAME%:latest
)

REM 推送到镜像仓库（如果配置了）
if not "%REGISTRY%"=="" (
    echo [INFO] 推送到镜像仓库: %REGISTRY%
    docker tag %IMAGE_NAME%:%VERSION% %REGISTRY%/%IMAGE_NAME%:%VERSION%
    docker push %REGISTRY%/%IMAGE_NAME%:%VERSION%
)

echo.
echo ========================================
echo 构建成功！
echo ========================================
echo 镜像: %IMAGE_NAME%:%VERSION%
echo.

REM 显示镜像信息
echo [INFO] 镜像详情:
docker images %IMAGE_NAME%:%VERSION%

echo.
echo 运行容器:
echo   docker run -p %HOST_PORT:8080%:8080 %IMAGE_NAME%:%VERSION%
echo.
echo 或者使用运行脚本:
echo   run.bat
echo.
echo 或者使用 docker-compose:
echo   docker-compose up -d
echo ========================================

ENDLOCAL
