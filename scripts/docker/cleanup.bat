@echo off
REM Docker 清理脚本 - 停止并删除容器和镜像
REM 用法: cleanup.bat [all]

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

REM 配置变量（从 .env 读取）
if "%IMAGE_NAME%"=="" set IMAGE_NAME=todolist
if "%CONTAINER_NAME%"=="" set CONTAINER_NAME=todo_server
if "%NETWORK_NAME%"=="" set NETWORK_NAME=todo_network

set CLEAN_ALL=%1

echo ========================================
echo Docker 清理脚本
echo ========================================
echo.

REM 停止并删除容器
echo [INFO] 检查容器...
docker ps -a -q --filter "name=%CONTAINER_NAME%" | findstr . >nul
if not errorlevel 1 (
    echo [INFO] 停止并删除容器: %CONTAINER_NAME%
    docker rm -f %CONTAINER_NAME%
    echo [OK] 容器已删除
) else (
    echo [INFO] 未发现容器: %CONTAINER_NAME%
)

REM 如果指定了 all 参数，同时删除镜像
if "%CLEAN_ALL%"=="all" (
    echo.
    echo [INFO] 检查镜像...
    docker images -q %IMAGE_NAME% | findstr . >nul
    if not errorlevel 1 (
        echo [INFO] 删除镜像: %IMAGE_NAME%
        docker rmi %IMAGE_NAME% 2>nul
        echo [OK] 镜像已删除
    ) else (
        echo [INFO] 未发现镜像: %IMAGE_NAME%
    )
)

echo.
echo ========================================
echo 清理完成！
echo ========================================
echo.

REM 显示当前运行的容器
echo [INFO] 网络 %NETWORK_NAME% 中的容器:
docker network ls --format "{{.Name}}" | findstr "^%NETWORK_NAME%$" >nul
if not errorlevel 1 (
    docker ps --filter "network=%NETWORK_NAME%" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
) else (
    echo [INFO] 网络 %NETWORK_NAME% 不存在
)

echo.

ENDLOCAL
