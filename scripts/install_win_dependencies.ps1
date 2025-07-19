if (!(Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Chocolatey..."
    Set-ExecutionPolicy Bypass -Scope Process -Force
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

    if (!(Get-Command choco -ErrorAction SilentlyContinue)) {
        Write-Error "Chocolatey installation failed. Exiting..."
        exit 1
    }
}

Import-Module "$env:ChocolateyInstall\helpers\chocolateyProfile.psm1" -ErrorAction SilentlyContinue
refreshenv

Write-Host "Installing rsvg-convert..."
choco install rsvg-convert -y
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to install rsvg-convert. Exiting..."
    exit 1
}

Write-Host "Installing Python..."
choco install python -y
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to install Python. Exiting..."
    exit 1
}

refreshenv

Write-Host "Upgrading pip..."
python -m ensurepip --upgrade
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to upgrade pip. Exiting..."
    exit 1
}

Write-Host "Installing svg.path..."
pip install svg.path
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to install svg.path. Exiting..."
    exit 1
}

$installGo = Read-Host "Do you want to install Go? (Y/N)"
if ($installGo -match '^[Yy]$') {
    Write-Host "Installing Go..."
    choco install golang -y
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to install Go. Exiting..."
        exit 1
    }
    refreshenv
    Write-Host "Go installed successfully!"
}
else {
    Write-Host "Skipping Go installation."
}

Write-Host "All installations completed successfully!"
