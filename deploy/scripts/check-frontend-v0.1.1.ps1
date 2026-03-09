Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$projectRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
$frontendDir = Join-Path $projectRoot "frontend"

function Assert-Contains([string]$Path, [string]$Pattern, [string]$Message) {
  $hit = Select-String -Path $Path -Pattern $Pattern -SimpleMatch -ErrorAction SilentlyContinue
  if ($null -eq $hit) {
    throw "Check failed: $Message"
  }
}

Write-Host "== Check: Frontend v0.1.1 =="

$uiStorePath = Join-Path $frontendDir "src\stores\ui-preferences.store.ts"
$modeSwitcherPath = Join-Path $frontendDir "src\components\ModeSwitcher.vue"
$cardViewPath = Join-Path $frontendDir "src\components\CardView.vue"
$profilePagePath = Join-Path $frontendDir "src\pages\ProfilePage.vue"
$profileApiPath = Join-Path $frontendDir "src\services\profile.api.ts"

Write-Host "[1/4] locale/theme persistence checks"
Assert-Contains $uiStorePath 'const LOCALE_KEY = "tarot_daily_locale"' "LOCALE_KEY missing"
Assert-Contains $uiStorePath 'const THEME_KEY = "tarot_daily_theme"' "THEME_KEY missing"
Assert-Contains $uiStorePath 'localStorage.setItem(LOCALE_KEY, nextLocale);' "locale persistence missing"
Assert-Contains $uiStorePath 'localStorage.setItem(THEME_KEY, nextTheme);' "theme persistence missing"
Assert-Contains $uiStorePath 'document.documentElement.setAttribute("data-theme", theme);' "theme apply missing"

Write-Host "[2/4] mode switcher options checks"
Assert-Contains $modeSwitcherPath '<option value="ru">' "ru option missing"
Assert-Contains $modeSwitcherPath '<option value="en">' "en option missing"
Assert-Contains $modeSwitcherPath '<option value="light">' "light option missing"
Assert-Contains $modeSwitcherPath '<option value="dark">' "dark option missing"

Write-Host "[3/4] profile page checks"
Assert-Contains $profilePagePath 'profileApi.getProfile(authStore.token)' "profile get call missing"
Assert-Contains $profilePagePath 'profileApi.getAvailableDecks(authStore.token)' "available decks call missing"
Assert-Contains $profileApiPath '"/api/me/profile"' "profile API endpoint missing"
Assert-Contains $profileApiPath '"/api/decks/available"' "available decks API endpoint missing"

Write-Host "[4/4] card flip checks"
Assert-Contains $cardViewPath 'const revealed = ref(false);' "initial closed state missing"
Assert-Contains $cardViewPath 'return revealed.value ? props.card.imageUrl : props.card.backImageUrl;' "front/back switch missing"
Assert-Contains $cardViewPath '@click="toggleReveal"' "flip click handler missing"
Assert-Contains $cardViewPath 'v-if="!revealed"' "hidden state block missing"

Write-Host "Frontend v0.1.1 checks passed"
