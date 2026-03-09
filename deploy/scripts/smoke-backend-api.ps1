param(
  [string]$ApiBaseUrl = "http://127.0.0.1:8080",
  [string]$AdminLogin = "admin",
  [string]$AdminPassword = "admin12345",
  [string]$UserPassword = "testpass123"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function New-JsonHeaders([string]$token = "") {
  $headers = @{ "Content-Type" = "application/json" }
  if ($token -ne "") {
    $headers["Authorization"] = "Bearer $token"
  }
  return $headers
}

function Invoke-Api(
  [string]$Method,
  [string]$Path,
  [object]$Body = $null,
  [string]$Token = ""
) {
  $params = @{
    Uri         = "$ApiBaseUrl$Path"
    Method      = $Method
    Headers     = New-JsonHeaders $Token
    ErrorAction = "Stop"
  }
  if ($Body -ne $null) {
    $params["Body"] = ($Body | ConvertTo-Json -Depth 8)
  }
  return Invoke-RestMethod @params
}

Write-Host "== Smoke: Backend API =="
Write-Host "API: $ApiBaseUrl"

$random = [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$userLogin = "smoke$random"

Write-Host "[1/7] Register user: $userLogin"
$register = Invoke-Api "POST" "/api/auth/register" @{
  login    = $userLogin
  password = $UserPassword
}
$userToken = $register.token

Write-Host "[2/7] Check /api/me"
$me = Invoke-Api "GET" "/api/me" $null $userToken
if ($me.login -ne $userLogin) {
  throw "Smoke failed: /api/me login mismatch"
}

Write-Host "[3/7] Update settings"
$settings = Invoke-Api "PATCH" "/api/me/settings" @{
  timezone    = "UTC+03:00"
  deckThemeId = ""
} $userToken
if ($settings.timezone -ne "UTC+03:00") {
  throw "Smoke failed: timezone did not update"
}

Write-Host "[4/7] Check stable cardOfDay"
$today1 = Invoke-Api "GET" "/api/reading/today" $null $userToken
$today2 = Invoke-Api "GET" "/api/reading/today" $null $userToken
if ($today1.cardOfDay.id -ne $today2.cardOfDay.id) {
  throw "Smoke failed: cardOfDay is not stable for same day"
}

Write-Host "[5/7] Generate spreads"
$generated = Invoke-Api "POST" "/api/reading/generate" $null $userToken
if ($null -eq $generated.spreadA -or $null -eq $generated.spreadB) {
  throw "Smoke failed: generate response has no spreads"
}

Write-Host "[6/7] Login as admin"
$admin = Invoke-Api "POST" "/api/auth/login" @{
  login    = $AdminLogin
  password = $AdminPassword
}
$adminToken = $admin.token

Write-Host "[7/7] Check admin endpoints"
$access = Invoke-Api "GET" "/api/admin/access-check" $null $adminToken
if ($access.scope -ne "admin") {
  throw "Smoke failed: admin access check failed"
}
$roles = Invoke-Api "GET" "/api/admin/roles" $null $adminToken
if ($null -eq $roles) {
  throw "Smoke failed: admin roles response is empty"
}

Write-Host "Smoke passed"
