# Build script for open-workbench-cli binaries
Write-Host "Building open-workbench-cli binaries for all platforms..." -ForegroundColor Green

# For macOS (Intel)
Write-Host "Building for macOS (Intel)..." -ForegroundColor Yellow
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o open-workbench-cli-darwin-amd64 main.go tui.go
if ($LASTEXITCODE -eq 0) { Write-Host "✅ macOS (Intel) build successful" -ForegroundColor Green } else { Write-Host "❌ macOS (Intel) build failed" -ForegroundColor Red }

# For macOS (Apple Silicon)
Write-Host "Building for macOS (Apple Silicon)..." -ForegroundColor Yellow
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o open-workbench-cli-darwin-arm64 main.go tui.go
if ($LASTEXITCODE -eq 0) { Write-Host "✅ macOS (Apple Silicon) build successful" -ForegroundColor Green } else { Write-Host "❌ macOS (Apple Silicon) build failed" -ForegroundColor Red }

# For Linux
Write-Host "Building for Linux..." -ForegroundColor Yellow
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o open-workbench-cli-linux-amd64 main.go tui.go
if ($LASTEXITCODE -eq 0) { Write-Host "✅ Linux build successful" -ForegroundColor Green } else { Write-Host "❌ Linux build failed" -ForegroundColor Red }

# For Windows (AMD64)
Write-Host "Building for Windows (AMD64)..." -ForegroundColor Yellow
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o open-workbench-cli-windows-amd64.exe main.go tui.go
if ($LASTEXITCODE -eq 0) { Write-Host "✅ Windows (AMD64) build successful" -ForegroundColor Green } else { Write-Host "❌ Windows (AMD64) build failed" -ForegroundColor Red }

# For Windows (ARM64)
Write-Host "Building for Windows (ARM64)..." -ForegroundColor Yellow
$env:GOOS="windows"; $env:GOARCH="arm64"; go build -o open-workbench-cli-windows-arm64.exe main.go tui.go
if ($LASTEXITCODE -eq 0) { Write-Host "✅ Windows (ARM64) build successful" -ForegroundColor Green } else { Write-Host "❌ Windows (ARM64) build failed" -ForegroundColor Red }

Write-Host "Build process completed!" -ForegroundColor Green 