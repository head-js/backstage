@echo off
rem build-windows.bat — 薄壳脚本，转发到 Makefile
rem 无参数：构建 Windows 版本（make build-windows）
rem 有参数：透传给 make，例如 build-windows.bat clean
cd /d "%~dp0"

if "%~1"=="" (
    make build-windows
) else (
    make %*
)
