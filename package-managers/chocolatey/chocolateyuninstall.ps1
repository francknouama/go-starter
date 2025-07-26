$ErrorActionPreference = 'Stop'

$packageName = 'go-starter'

# Remove the binary from chocolatey bin directory
$binDirectory = Join-Path (Get-ToolsLocation) "chocolatey\bin"
$shimPath = Join-Path $binDirectory "go-starter.exe"

if (Test-Path $shimPath) {
    Remove-Item -Path $shimPath -Force
    Write-Host "Removed go-starter.exe from $binDirectory" -ForegroundColor Green
}

# Remove any leftover files from the package tools directory
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$exeFiles = Get-ChildItem -Path $toolsDir -Filter "go-starter*.exe" -Recurse

foreach ($file in $exeFiles) {
    if (Test-Path $file.FullName) {
        Remove-Item -Path $file.FullName -Force
        Write-Host "Removed $($file.Name)" -ForegroundColor Green
    }
}

Write-Host "go-starter has been uninstalled successfully!" -ForegroundColor Green