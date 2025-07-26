$ErrorActionPreference = 'Stop'

$packageName = 'go-starter'
$version = '1.4.0'
$url64 = 'https://github.com/francknouama/go-starter/releases/download/v1.4.0/go-starter-v1.4.0-windows-amd64.zip'
$checksum64 = 'placeholder_checksum_windows_amd64' # Will be updated by release automation
$checksumType = 'sha256'

$packageArgs = @{
    packageName   = $packageName
    unzipLocation = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
    url64bit      = $url64
    checksum64    = $checksum64
    checksumType64= $checksumType
}

# Download and extract
Install-ChocolateyZipPackage @packageArgs

# Find the extracted binary
$binPath = Get-ChildItem -Path $packageArgs.unzipLocation -Filter "go-starter*.exe" -Recurse | Select-Object -First 1 -ExpandProperty FullName

if (-not $binPath) {
    throw "Could not find go-starter.exe in the extracted files"
}

# Create a shim in the chocolatey bin directory
$binDirectory = Join-Path (Get-ToolsLocation) "chocolatey\bin"
if (-not (Test-Path $binDirectory)) {
    New-Item -ItemType Directory -Path $binDirectory -Force | Out-Null
}

$shimPath = Join-Path $binDirectory "go-starter.exe"

# Copy the binary to the bin directory
Copy-Item -Path $binPath -Destination $shimPath -Force

Write-Host "go-starter has been installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Quick start:" -ForegroundColor Yellow
Write-Host "  go-starter new my-project --type=cli" -ForegroundColor White
Write-Host ""
Write-Host "Interactive mode:" -ForegroundColor Yellow  
Write-Host "  go-starter new" -ForegroundColor White
Write-Host ""
Write-Host "Advanced mode:" -ForegroundColor Yellow
Write-Host "  go-starter new --advanced" -ForegroundColor White
Write-Host ""
Write-Host "For help: go-starter --help" -ForegroundColor Cyan
Write-Host "Documentation: https://github.com/francknouama/go-starter" -ForegroundColor Cyan