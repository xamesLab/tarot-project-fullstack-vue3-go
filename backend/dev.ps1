param(
  [string]$EnvFile = ".env"
)

$ErrorActionPreference = "Stop"

$backendDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Push-Location $backendDir

try {
  if (-not (Test-Path $EnvFile)) {
    Write-Host "warning: '$EnvFile' not found in backend folder; app will rely on process env vars"
  }

  $goCmd = Get-Command go -ErrorAction SilentlyContinue
  $goBin = if ($goCmd) { $goCmd.Source } else { "C:\Program Files\Go\bin\go.exe" }

  if (-not (Test-Path $goBin) -and $goBin -ne "go") {
    throw "go binary not found. Install Go or add it to PATH."
  }

  $env:GO111MODULE = "on"
  & $goBin run ./cmd/api
}
finally {
  Pop-Location
}
