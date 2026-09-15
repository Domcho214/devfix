# devfix installer


Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "         Installing devfix...            " -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

# Set up paths
$InstallDir = "$env:LOCALAPPDATA\devfix\bin"
$ExePath = "$InstallDir\devfix.exe"
$RepoUrl = "https://raw.githubusercontent.com/Domcho214/devfix/main/devfix.exe"

# Create installation directory if it doesn't exist
if (-not (Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

Write-Host "`nDownloading devfix..."
try {
    Invoke-WebRequest -Uri $RepoUrl -OutFile $ExePath -UseBasicParsing
    Write-Host "Download complete!" -ForegroundColor Green
} catch {
    Write-Host "Failed to download devfix. Check github link" -ForegroundColor Red
    Write-Host $_.Exception.Message
    exit
}

# Add to PATH if not already there
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "Adding devfix..."
    $NewPath = $UserPath + ";$InstallDir"
    [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
    $env:PATH = $env:PATH + ";$InstallDir"
}

Write-Host "`n=========================================" -ForegroundColor Green
Write-Host "  devfix is installed successfully!        " -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Green
Write-Host "`nYou can now run it with the command:" -ForegroundColor White
Write-Host "  devfix" -ForegroundColor Yellow
Write-Host ""
