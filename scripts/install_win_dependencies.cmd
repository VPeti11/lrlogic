@echo off
SETLOCAL

where choco >nul 2>nul
IF %ERRORLEVEL% NEQ 0 (
    Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
    IF %ERRORLEVEL% NEQ 0 (
        echo Failed to install Chocolatey. Exiting...
        exit /b 1
    )
)

choco install rsvg-convert -y
IF %ERRORLEVEL% NEQ 0 (
    echo Failed to install rsvg-convert. Exiting...
    exit /b 1
)

choco install python -y
IF %ERRORLEVEL% NEQ 0 (
    echo Failed to install Python. Exiting...
    exit /b 1
)

python -m ensurepip --upgrade
IF %ERRORLEVEL% NEQ 0 (
    echo Failed to upgrade pip. Exiting...
    exit /b 1
)

pip install svg.path
IF %ERRORLEVEL% NEQ 0 (
    echo Failed to install svg.path package. Exiting...
    exit /b 1
)

echo Do you want to install Go? (Y/N)
set /p userInput= 
IF /I "%userInput%"=="Y" (
    choco install golang -y
    IF %ERRORLEVEL% NEQ 0 (
        echo Failed to install Go. Exiting...
        exit /b 1
    )
    echo Go installed successfully!
) ELSE (
    echo Skipping Go installation.
)

echo All installations completed successfully!
ENDLOCAL
