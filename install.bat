@echo off
setlocal

:: Function to install the application
:install_app
echo Installing dockermi...
move dockermi.exe C:\Program Files\dockermi\ >nul 2>&1
if errorlevel 1 (
    echo Installation failed! Creating the directory...
    mkdir "C:\Program Files\dockermi"
    move dockermi.exe C:\Program Files\dockermi\
)
echo dockermi installed! You can run it by typing 'C:\Program Files\dockermi\dockermi.exe' in your command prompt.
exit /b 0

:: Start the build and install process
call :install_app
