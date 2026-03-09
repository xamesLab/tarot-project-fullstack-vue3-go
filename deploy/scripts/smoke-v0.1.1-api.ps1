param(
  [string]$ApiBaseUrl = "http://127.0.0.1:8080",
  [string]$AdminLogin = "admin",
  [string]$AdminPassword = "admin12345",
  [string]$UserPassword = "testpass123",
  [string]$BackImagePath = "specs/temp/card-image/back.jpg",
  [string]$CardImagePath = "specs/temp/card-image/22.jpg"
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
    $params["Body"] = ($Body | ConvertTo-Json -Depth 10)
  }
  return Invoke-RestMethod @params
}

function Invoke-ApiMultipart(
  [string]$Path,
  [string]$Token,
  [string]$FilePath
) {
  if (-not (Test-Path $FilePath)) {
    throw "File not found: $FilePath"
  }

  $resolved = (Resolve-Path $FilePath).Path
  $raw = & curl.exe -sS -w "`n%{http_code}" -X POST "$ApiBaseUrl$Path" -H "Authorization: Bearer $Token" -F "file=@$resolved"
  if ($LASTEXITCODE -ne 0) {
    throw "curl failed for $Path"
  }

  $parts = $raw -split "`n"
  if ($parts.Count -lt 2) {
    throw "Unexpected multipart response for $Path"
  }

  $status = [int]$parts[$parts.Count - 1]
  $body = ($parts[0..($parts.Count - 2)] -join "`n").Trim()
  if ($status -lt 200 -or $status -ge 300) {
    throw "HTTP $status for $Path`n$body"
  }
  if ([string]::IsNullOrWhiteSpace($body)) {
    return $null
  }
  return $body | ConvertFrom-Json
}

Write-Host "== Smoke: API v0.1.1 =="
Write-Host "API: $ApiBaseUrl"

$random = [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$userLogin = "v011$random"
$displayName = "v011-user-$random"
$deckName = "v011-deck-$random"

Write-Host "[1/14] Register user: $userLogin"
$register = Invoke-Api "POST" "/api/auth/register" @{
  login    = $userLogin
  password = $UserPassword
}
$userToken = $register.token

Write-Host "[2/14] Check profile endpoint"
$meProfile = Invoke-Api "GET" "/api/me/profile" $null $userToken
if ($meProfile.login -ne $userLogin) {
  throw "Smoke failed: /api/me/profile login mismatch"
}

Write-Host "[3/14] Patch profile (name/timezone, clear preferredDeckId)"
$updatedProfile = Invoke-Api "PATCH" "/api/me/profile" @{
  displayName     = $displayName
  timezone        = "UTC+03:00"
  preferredDeckId = ""
} $userToken
if ($updatedProfile.timezone -ne "UTC+03:00") {
  throw "Smoke failed: profile timezone mismatch"
}

Write-Host "[4/14] Check stable card of day + image/back image payload"
$today1 = Invoke-Api "GET" "/api/reading/today" $null $userToken
$today2 = Invoke-Api "GET" "/api/reading/today" $null $userToken
if ($today1.cardOfDay.id -ne $today2.cardOfDay.id) {
  throw "Smoke failed: cardOfDay must be stable within day"
}
if ([string]::IsNullOrWhiteSpace($today1.cardOfDay.imageUrl) -or [string]::IsNullOrWhiteSpace($today1.cardOfDay.backImageUrl)) {
  throw "Smoke failed: reading card payload must include imageUrl and backImageUrl"
}

Write-Host "[5/14] Login as admin"
$admin = Invoke-Api "POST" "/api/auth/login" @{
  login    = $AdminLogin
  password = $AdminPassword
}
$adminToken = $admin.token

Write-Host "[6/14] Check admin users exact search + filters"
$usersSearch = Invoke-Api "GET" "/api/admin/users?page=1&pageSize=20&name=$([uri]::EscapeDataString($displayName))&status=active&sort=registeredAt:desc" $null $adminToken
if ($usersSearch.total -lt 1) {
  throw "Smoke failed: admin users search did not find the created user"
}

Write-Host "[7/14] Create admin deck"
$createdDeck = Invoke-Api "POST" "/api/admin/decks" @{ name = $deckName } $adminToken
if ($createdDeck.status -ne "process") {
  throw "Smoke failed: new deck status must be process"
}

Write-Host "[8/14] Upload deck back image"
$deckAfterBack = Invoke-ApiMultipart "/api/admin/decks/$($createdDeck.id)/assets/back" $adminToken $BackImagePath
if ([string]::IsNullOrWhiteSpace($deckAfterBack.backImageUrl)) {
  throw "Smoke failed: back image URL was not set"
}

Write-Host "[9/14] Upload one deck card image"
$deckAfterCard = Invoke-ApiMultipart "/api/admin/decks/$($createdDeck.id)/assets/cards/22" $adminToken $CardImagePath
if ($deckAfterCard.cardsUploadedCount -lt 1) {
  throw "Smoke failed: cardsUploadedCount did not increase"
}

Write-Host "[10/14] Check deck list/details"
$decksList = Invoke-Api "GET" "/api/admin/decks" $null $adminToken
$createdFromList = $decksList | Where-Object { $_.id -eq $createdDeck.id } | Select-Object -First 1
if ($null -eq $createdFromList) {
  throw "Smoke failed: created deck not found in list"
}
$deckDetails = Invoke-Api "GET" "/api/admin/decks/$($createdDeck.id)" $null $adminToken
if ($deckDetails.cardsUploadedCount -lt 1) {
  throw "Smoke failed: deck details did not return uploaded card state"
}

Write-Host "[11/14] Ensure /api/decks/available returns only done decks"
$availableDecks = Invoke-Api "GET" "/api/decks/available" $null $userToken
$processDeckInAvailable = $availableDecks | Where-Object { $_.id -eq $createdDeck.id } | Select-Object -First 1
if ($null -ne $processDeckInAvailable) {
  throw "Smoke failed: process deck must not appear in available decks"
}

Write-Host "[12/14] Ensure user cannot set preferredDeckId to process deck"
try {
  $null = Invoke-Api "PATCH" "/api/me/profile" @{
    preferredDeckId = $createdDeck.id
  } $userToken
  throw "Smoke failed: setting process deck as preferredDeckId should return error"
} catch {
  $message = $_.Exception.Message
  if ($message -notlike "*(400)*") {
    throw "Smoke failed: expected 400 for invalid preferredDeckId, got: $message"
  }
}

Write-Host "[13/14] Upload remaining 77 cards and verify status=done"
$allCardKeys = @()
for ($i = 22; $i -le 77; $i++) { $allCardKeys += [string]$i }
for ($i = 100; $i -le 121; $i++) { $allCardKeys += [string]$i }
foreach ($cardKey in $allCardKeys) {
  if ($cardKey -eq "22") { continue }
  $cardPath = "specs/temp/card-image/$cardKey.jpg"
  $deckAfterCard = Invoke-ApiMultipart "/api/admin/decks/$($createdDeck.id)/assets/cards/$cardKey" $adminToken $cardPath
}
if ($deckAfterCard.status -ne "done") {
  throw "Smoke failed: deck status should become done after full upload"
}
if ($deckAfterCard.cardsUploadedCount -ne 78) {
  throw "Smoke failed: expected cardsUploadedCount=78 after full upload"
}

Write-Host "[14/14] Ensure done deck is available and can be selected in profile"
$availableAfterDone = Invoke-Api "GET" "/api/decks/available" $null $userToken
$doneDeck = $availableAfterDone | Where-Object { $_.id -eq $createdDeck.id } | Select-Object -First 1
if ($null -eq $doneDeck) {
  throw "Smoke failed: done deck must appear in /api/decks/available"
}
$profileWithDeck = Invoke-Api "PATCH" "/api/me/profile" @{
  preferredDeckId = $createdDeck.id
} $userToken
if ($profileWithDeck.preferredDeckId -ne $createdDeck.id) {
  throw "Smoke failed: preferredDeckId was not set to done deck"
}

Write-Host "Smoke v0.1.1 passed"
