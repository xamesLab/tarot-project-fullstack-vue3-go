param(
  [string]$ApiBaseUrl = "http://127.0.0.1:8080",
  [string]$UserPassword = "testpass123",
  [string]$CallbackHmacSecret = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function New-JsonHeaders(
  [string]$Token = "",
  [hashtable]$Extra = $null
) {
  $headers = @{ "Content-Type" = "application/json" }
  if ($Token -ne "") {
    $headers["Authorization"] = "Bearer $Token"
  }
  if ($null -ne $Extra) {
    foreach ($key in $Extra.Keys) {
      $headers[$key] = $Extra[$key]
    }
  }
  return $headers
}

function Invoke-Api(
  [string]$Method,
  [string]$Path,
  [object]$Body = $null,
  [string]$Token = "",
  [hashtable]$ExtraHeaders = $null,
  [string]$RawBody = ""
) {
  $params = @{
    Uri         = "$ApiBaseUrl$Path"
    Method      = $Method
    Headers     = New-JsonHeaders -Token $Token -Extra $ExtraHeaders
    ErrorAction = "Stop"
  }
  if (-not [string]::IsNullOrWhiteSpace($RawBody)) {
    $params["Body"] = $RawBody
  } elseif ($Body -ne $null) {
    $params["Body"] = ($Body | ConvertTo-Json -Depth 10)
  }
  return Invoke-RestMethod @params
}

function New-HmacSignature(
  [string]$Body,
  [string]$Secret
) {
  $hmac = [System.Security.Cryptography.HMACSHA256]::new([System.Text.Encoding]::UTF8.GetBytes($Secret))
  try {
    $hash = $hmac.ComputeHash([System.Text.Encoding]::UTF8.GetBytes($Body))
    return [System.BitConverter]::ToString($hash).Replace("-", "").ToLowerInvariant()
  } finally {
    $hmac.Dispose()
  }
}

function Invoke-GatewayCallback(
  [object]$Payload
) {
  $rawJson = $Payload | ConvertTo-Json -Depth 10 -Compress
  $headers = @{}
  if (-not [string]::IsNullOrWhiteSpace($CallbackHmacSecret)) {
    $signature = New-HmacSignature -Body $rawJson -Secret $CallbackHmacSecret
    $headers["X-Signature"] = "sha256=$signature"
  }
  return Invoke-Api -Method "POST" -Path "/api/gateway/callback" -RawBody $rawJson -ExtraHeaders $headers
}

Write-Host "== Smoke: API v0.1.2 ask-flow =="
Write-Host "API: $ApiBaseUrl"

$random = [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$userLogin = "v012ask$random"
$successText = "success-result-$random"
$errorCode = "agent_timeout"
$errorMessage = "agent timeout in smoke $random"

Write-Host "[1/7] Register user: $userLogin"
$register = Invoke-Api "POST" "/api/auth/register" @{
  login    = $userLogin
  password = $UserPassword
}
$userToken = $register.token
if ([string]::IsNullOrWhiteSpace($userToken)) {
  throw "Smoke failed: token is empty after register"
}

Write-Host "[2/7] Create ask request (success-flow)"
$askSuccess = Invoke-Api "POST" "/api/ask" @{
  text = "ask-success-$random"
} $userToken
if ([string]::IsNullOrWhiteSpace($askSuccess.request_id)) {
  throw "Smoke failed: request_id is empty for success-flow"
}

Write-Host "[3/7] Send callback success"
$callbackSuccess = Invoke-GatewayCallback @{
  request_id = $askSuccess.request_id
  status     = "success"
  text       = $successText
}
if ($callbackSuccess.status -ne "ok") {
  throw "Smoke failed: callback success was not acknowledged"
}

Write-Host "[4/7] Validate success-flow status"
$askSuccessState = Invoke-Api "GET" "/api/ask/$($askSuccess.request_id)" $null $userToken
if ($askSuccessState.status -ne "success") {
  throw "Smoke failed: expected status=success, got $($askSuccessState.status)"
}
if ($askSuccessState.result_text -ne $successText) {
  throw "Smoke failed: success result_text mismatch"
}

Write-Host "[5/7] Create ask request (error-flow)"
$askError = Invoke-Api "POST" "/api/ask" @{
  text = "ask-error-$random"
} $userToken
if ([string]::IsNullOrWhiteSpace($askError.request_id)) {
  throw "Smoke failed: request_id is empty for error-flow"
}

Write-Host "[6/7] Send callback error"
$callbackError = Invoke-GatewayCallback @{
  request_id = $askError.request_id
  status     = "error"
  error      = $errorCode
  message    = $errorMessage
}
if ($callbackError.status -ne "ok") {
  throw "Smoke failed: callback error was not acknowledged"
}

Write-Host "[7/7] Validate error-flow status"
$askErrorState = Invoke-Api "GET" "/api/ask/$($askError.request_id)" $null $userToken
if ($askErrorState.status -ne "error") {
  throw "Smoke failed: expected status=error, got $($askErrorState.status)"
}
if ($askErrorState.error_code -ne $errorCode) {
  throw "Smoke failed: error_code mismatch"
}
if ($askErrorState.error_message -ne $errorMessage) {
  throw "Smoke failed: error_message mismatch"
}

Write-Host "Smoke v0.1.2 ask-flow passed"
