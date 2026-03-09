Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
$env:GOTELEMETRY = "off"
if (Get-Variable PSNativeCommandUseErrorActionPreference -ErrorAction SilentlyContinue) {
    $PSNativeCommandUseErrorActionPreference = $false
}

function Test-Command {
    param(
        [Parameter(Mandatory = $true)][string]$Name,
        [Parameter(Mandatory = $true)][string]$Command
    )

    try {
        $output = Invoke-Expression $Command 2>&1 | Out-String
        return [pscustomobject]@{
            Name = $Name
            Ok = $true
            Details = $output.Trim()
        }
    } catch {
        return [pscustomobject]@{
            Name = $Name
            Ok = $false
            Details = $_.Exception.Message
        }
    }
}

$checks = @()

# Go can be missing from PATH right after installation, so try explicit path fallback.
$goPath = "C:\Program Files\Go\bin\go.exe"
if (Get-Command go -ErrorAction SilentlyContinue) {
    $checks += Test-Command -Name "go" -Command "go version"
} elseif (Test-Path $goPath) {
    $checks += Test-Command -Name "go" -Command "& '$goPath' version"
} else {
    $checks += [pscustomobject]@{
        Name = "go"
        Ok = $false
        Details = "Go is not installed or unavailable."
    }
}

$checks += Test-Command -Name "node" -Command "node -v"
$checks += Test-Command -Name "npm" -Command "npm.cmd -v"
try {
    $sshOutput = (cmd /c "ssh -V 2>&1" | Out-String).Trim()
    $sshOk = $sshOutput -match "OpenSSH"
    $checks += [pscustomobject]@{
        Name = "ssh"
        Ok = $sshOk
        Details = $sshOutput
    }
} catch {
    $checks += [pscustomobject]@{
        Name = "ssh"
        Ok = $false
        Details = $_.Exception.Message
    }
}

# Optional tool for this project.
if (Get-Command docker -ErrorAction SilentlyContinue) {
    $checks += Test-Command -Name "docker(optional)" -Command "docker --version"
} else {
    $checks += [pscustomobject]@{
        Name = "docker(optional)"
        Ok = $true
        Details = "Not installed (allowed for current MVP dev flow)."
    }
}

$requiredFailed = $false

Write-Host "=== Dev Environment Check ==="
foreach ($check in $checks) {
    if ($check.Ok) {
        Write-Host "[OK]  $($check.Name): $($check.Details)"
    } else {
        Write-Host "[ERR] $($check.Name): $($check.Details)"
    }
}

foreach ($check in $checks) {
    if (-not $check.Ok -and $check.Name -notlike "docker(optional)*") {
        $requiredFailed = $true
    }
}

if ($requiredFailed) {
    Write-Host ""
    Write-Host "Result: FAILED (required tools are missing)."
    exit 1
}

Write-Host ""
Write-Host "Result: PASSED (required tools are available)."
