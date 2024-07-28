@echo off
setlocal

:: Function to build the application for Windows
:build_windows
echo Building for Windows...
set GOOS=windows
set GOARCH=amd64
go build -o dockermi.exe dockermi.go
if errorlevel 1 (
    echo Build failed!
    exit /b 1
)
echo Build completed! The executable is dockermi.exe.
exit /b 0

:: Start the build and install process
call :build_windows
