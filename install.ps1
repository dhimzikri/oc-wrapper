# OpenCode Wrapper Installer
# Usage: irm https://raw.githubusercontent.com/dhimzikri/oc-wrapper/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$installDir = "$env:USERPROFILE\.config\opencode\wrapper"
$repoUrl = "https://github.com/dhimzikri/oc-wrapper"

Write-Host "Installing OpenCode Wrapper..." -ForegroundColor Cyan

# Create install directory
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

# Check if Go is installed
$hasGo = Get-Command go -ErrorAction SilentlyContinue

if ($hasGo) {
    Write-Host "Building from source..." -ForegroundColor Yellow
    
    $tempDir = "$env:TEMP\oc-wrapper-install"
    if (Test-Path $tempDir) { Remove-Item -Recurse -Force $tempDir }
    
    git clone --depth 1 $repoUrl $tempDir
    Push-Location $tempDir
    go build -o "$installDir\opencode.exe" main.go
    Pop-Location
    Remove-Item -Recurse -Force $tempDir
} else {
    Write-Host "Go not found. Downloading pre-built binary..." -ForegroundColor Yellow
    $binaryUrl = "$repoUrl/releases/latest/download/opencode-wrapper-windows-amd64.exe"
    Invoke-WebRequest -Uri $binaryUrl -OutFile "$installDir\opencode.exe"
}

# Add to PowerShell profile
$profileContent = Get-Content $PROFILE -Raw -ErrorAction SilentlyContinue
$functionDef = 'function opencode { & "$env:USERPROFILE\.config\opencode\wrapper\opencode.exe" @args }'

if ($profileContent -notmatch "oc-wrapper|opencode-wrapper") {
    Write-Host "Adding to PowerShell profile..." -ForegroundColor Yellow
    Add-Content -Path $PROFILE -Value "`n# OpenCode Wrapper`n$functionDef"
} else {
    Write-Host "Profile already configured, skipping..." -ForegroundColor Yellow
}

Write-Host "`nInstallation complete!" -ForegroundColor Green
Write-Host "Restart PowerShell or run: . `$PROFILE" -ForegroundColor Cyan
Write-Host "Then use: opencode" -ForegroundColor Cyan
