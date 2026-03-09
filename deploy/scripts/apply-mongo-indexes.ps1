Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$projectRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$backendEnvPath = Join-Path $projectRoot "backend\.env"
$indexScriptPath = Join-Path $projectRoot "deploy\mongo\create-indexes.js"

if (-not (Test-Path $backendEnvPath)) {
    throw "backend/.env not found: $backendEnvPath"
}

if (-not (Get-Command mongosh -ErrorAction SilentlyContinue)) {
    throw "mongosh is not installed or not available in PATH."
}

$envMap = @{}
Get-Content $backendEnvPath | ForEach-Object {
    $line = $_.Trim()
    if (-not $line -or $line.StartsWith("#")) {
        return
    }
    $parts = $line -split "=", 2
    if ($parts.Count -eq 2) {
        $key = $parts[0].Trim()
        $value = $parts[1].Trim()
        $envMap[$key] = $value
    }
}

if (-not $envMap.ContainsKey("MONGO_URI")) {
    throw "MONGO_URI is missing in backend/.env"
}

$mongoUri = $envMap["MONGO_URI"]
$mongoDbName = if ($envMap.ContainsKey("MONGO_DB_NAME")) { $envMap["MONGO_DB_NAME"] } else { "tarot_daily" }

Write-Host "Applying indexes to database: $mongoDbName"
Write-Host "Using MONGO_URI: $mongoUri"

& mongosh "$mongoUri" --eval "const DB_NAME='$mongoDbName'" --file "$indexScriptPath"

Write-Host ""
Write-Host "Mongo indexes applied successfully."
