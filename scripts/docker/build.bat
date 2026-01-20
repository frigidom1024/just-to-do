@echo off
REM Docker 镜像构建脚本
REM 用法: build.bat [version]
REM 示例: build.bat 1.0.0

SETLOCAL EnableDelayedExpansion

REM 切换到项目根目录
cd /d %~dp0..\..

REM 配置变量
if "%IMAGE_NAME%"=="" set IMAGE_NAME=todolist
if "%~1"=="" (
    set "VERSION=latest"
) else (
    set "VERSION=%~1"
)

REM 检查 Dockerfile 是否存在
if not exist "Dockerfile" (
    echo 错误: 未找到 Dockerfile
    exit /b 1
)

REM 构建镜像
echo 开始构建 Docker 镜像...
echo 镜像名称: %IMAGE_NAME%
echo 版本标签: %VERSION%
echo.

docker build -t %IMAGE_NAME%:%VERSION% .

if errorlevel 1 (
    echo.
    echo 错误: Docker 镜像构建失败
    exit /b 1
)

REM 如果版本不是 latest，额外打一个 latest 标签
if not "%VERSION%"=="latest" (
    docker tag %IMAGE_NAME%:%VERSION% %IMAGE_NAME%:latest
)

echo.
echo 构建成功！
echo 镜像: %IMAGE_NAME%:%VERSION%

ENDLOCAL
