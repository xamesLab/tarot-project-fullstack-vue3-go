# Telegram Mini App Integration Plan
## Nebula Tarot

**Status:** planning only, implementation not started  
**Scope:** one shared web application that works both in a browser and inside Telegram Mini App  
**Primary goal:** make the existing frontend available as a Telegram Mini App and synchronize Telegram authorization with browser authorization through the same backend user account.

---

## 1. Goal

Build Telegram Mini App support without creating a separate product or a forked frontend.

The Mini App should:

- open the same `Vue 3` SPA that is already used by the browser site;
- use the same Go/Gin backend API and MongoDB data;
- preserve existing reading, profile, history, deck, theme, language and admin behavior;
- authenticate Telegram users through Telegram-provided signed data;
- allow the same person to use one Nebula Tarot account from both browser and Telegram.

Important clarification: browser and Telegram WebView storage should not be treated as shared. Synchronization must happen on the server account level: browser and Mini App receive separate app JWTs that point to the same `users._id`.

---

## 2. Current Project Baseline

### 2.1 Frontend

Current stack:

- `Vue 3`
- `Vite`
- `Pinia`
- `Vue Router`
- `Axios`
- `Vue I18n`
- `SCSS`
- `VeeValidate`

Relevant current files:

- `frontend/src/main.js` initializes Vue, Pinia, router and i18n.
- `frontend/src/App.vue` hydrates app state and auth before rendering the shell.
- `frontend/src/router/index.js` protects authenticated routes and redirects guests to `/login`.
- `frontend/src/services/api.js` creates the shared Axios client.
- `frontend/src/stores/auth.js` stores the app JWT in `localStorage` under `tarot_daily_token`.
- `frontend/src/components/layout/AppHeader.vue` owns browser navigation, theme/language controls and logout.

Current auth frontend behavior:

- browser login/register calls `/api/auth/login` and `/api/auth/register`;
- JWT is stored in `localStorage`;
- API requests use `Authorization: Bearer <token>`;
- `/me` hydrates the user after reload.

### 2.2 Backend

Current stack:

- `Go 1.23`
- `Gin`
- `MongoDB`
- `mongo-go-driver`
- `JWT`
- `bcrypt`
- `validator`
- `godotenv`

Relevant current files:

- `backend/main.go` wires repositories, services, handlers, routes and middleware.
- `backend/internal/handlers/auth.go` exposes `POST /api/auth/register` and `POST /api/auth/login`.
- `backend/internal/services/auth.go` owns email/password registration and login.
- `backend/internal/models/user.go` defines `User`.
- `backend/internal/models/auth.go` defines register/login requests and JWT claims.
- `backend/internal/utils/jwt.go` generates and verifies app JWTs.
- `backend/internal/middleware/auth.go` protects app API routes through app JWT verification.
- `backend/internal/middleware/cors.go` currently allows one configured frontend origin plus local aliases.

Current auth backend behavior:

- user identity is email/password based;
- `users.email` and `users.username` are unique;
- no Telegram identity fields exist yet;
- app JWT contains `user_id`, `role`, `timezone` and standard registered claims.

---

## 3. External Telegram Facts To Account For

Reference sources are listed at the end of this document.

### 3.1 Mini App Runtime

Telegram Mini Apps are web apps launched inside Telegram. Telegram documents several launch surfaces: main Mini App profile button, keyboard button, inline button, bot menu button, inline mode, direct link and attachment menu.

For Nebula Tarot the practical launch targets should be:

- Main Mini App button on the bot profile;
- bot menu button;
- direct link like `https://t.me/<botusername>?startapp`;
- optionally an inline button later, if the bot starts sending messages.

Mini Apps need an HTTPS URL. For local device testing this usually means a public HTTPS tunnel to the local Vite/backend setup, or a staging deployment.

### 3.2 Telegram WebApp Object

Telegram exposes Mini App runtime data through `window.Telegram.WebApp` after loading:

```html
<script src="https://telegram.org/js/telegram-web-app.js"></script>
```

Important fields/methods for this project:

- `initData`: raw signed launch data, must be sent to backend for validation;
- `initDataUnsafe`: convenient parsed data, but not trusted for auth;
- `themeParams` and `colorScheme`: Telegram theme colors;
- `viewportHeight` and `viewportStableHeight`: current and stable Mini App viewport height;
- `safeAreaInset` and `contentSafeAreaInset`: device/control safe areas;
- `ready()`: signal that the app is ready;
- `expand()`: request full available height;
- `BackButton`: integrate Telegram-native back behavior with Vue Router;
- `HapticFeedback`: optional interaction polish.

### 3.3 Init Data Validation

Telegram explicitly warns that `initDataUnsafe` must not be trusted. The backend must validate the raw `initData`.

Validation algorithm for Mini App init data:

1. Parse init data as query string.
2. Exclude and remember the `hash` field.
3. Sort all remaining key-value pairs alphabetically.
4. Join them with `\n` into a data-check string.
5. Build secret key as `HMAC-SHA256(bot_token, "WebAppData")`.
6. Build calculated hash as `HMAC-SHA256(data_check_string, secret_key)` in hex.
7. Compare calculated hash with provided `hash` using constant-time comparison.
8. Check `auth_date` freshness to reduce replay risk.

Only after this validation can the backend trust Telegram user ID and other user fields.

### 3.4 Browser Telegram Login

For browser-side Telegram login/linking, Telegram now documents Telegram Login with OpenID Connect.

Recommended browser path:

- register allowed website origins and redirect URIs in BotFather;
- use Telegram Login OIDC with Authorization Code Flow + PKCE;
- validate `id_token` server-side through Telegram JWKS;
- verify `iss`, `aud`, `exp`, nonce/state and signature;
- then issue Nebula Tarot app JWT.

The legacy Telegram Login Widget still exists, but the OIDC flow is the better default for new implementation. The legacy widget can stay as fallback only if OIDC setup or browser support becomes a blocker.

---

## 4. Core Architecture Decision

Use one frontend app and one backend.

Do not create:

- a separate `miniapp/` frontend project;
- a second backend API;
- duplicated route/page implementations for Telegram.

Add a small Telegram runtime layer to the existing SPA:

- detect Telegram runtime at startup through `window.Telegram?.WebApp` and/or launch params;
- auto-authenticate through Telegram only when Telegram runtime is present and no app JWT exists;
- keep existing routes, stores and services after the app JWT is obtained;
- make Telegram-specific UI adaptations conditional and shallow.

Resulting model:

```text
Browser
  -> same Vue SPA
  -> email/password or Telegram Login OIDC
  -> app JWT
  -> same Go API
  -> same MongoDB user

Telegram Mini App
  -> same Vue SPA inside Telegram WebView
  -> Telegram initData token exchange
  -> app JWT
  -> same Go API
  -> same MongoDB user
```

---

## 5. Authorization And Account Synchronization

### 5.1 Account Identity Rule

The canonical app identity remains the MongoDB user record.

Telegram identity should be attached to that user record, not stored as a separate user system.

Minimum user identity additions:

- `telegram_id`: Telegram numeric user ID, stored as string or int64 consistently;
- `telegram_username`;
- `telegram_first_name`;
- `telegram_last_name`;
- `telegram_photo_url`;
- `telegram_linked_at`;
- `telegram_last_login_at`.

Recommended structural direction:

- allow users to have multiple auth providers;
- keep email/password accounts working;
- allow Telegram-only accounts without fake email addresses;
- use partial unique index for `email` if email becomes optional;
- use unique index for `telegram_id` when present.

Avoid generating fake emails like `tg_<id>@telegram.local` unless a short-term migration constraint forces it. Fake emails leak implementation detail into product logic and make account merge/support harder.

### 5.2 Mini App Login Flow

Target endpoint:

```text
POST /api/auth/telegram/miniapp
```

Request:

```json
{
  "init_data": "<raw Telegram.WebApp.initData>",
  "timezone": "Europe/Moscow"
}
```

Backend flow:

1. Validate request body.
2. Validate `init_data` signature with `TELEGRAM_BOT_TOKEN`.
3. Validate `auth_date` freshness, for example `<= 24h` for normal login.
4. Extract Telegram user ID and profile fields.
5. If `telegram_id` is linked to an active Nebula user, update Telegram profile snapshot and issue app JWT.
6. If `telegram_id` is linked to a disabled user, return the existing disabled-account error.
7. If no linked user exists, create a new `seeker` user with Telegram provider data, default language/theme/timezone, and issue app JWT.
8. Return the same response shape as existing login/register where possible: `{ success, user, token }`.

Frontend flow:

1. On app startup, detect Telegram runtime.
2. Call `WebApp.ready()` after the shell can render.
3. If no app token exists, send `WebApp.initData` to backend.
4. Persist returned app JWT with the existing auth store path.
5. Continue using existing protected routes and services.

### 5.3 Browser Telegram Login Flow

Target endpoints:

```text
POST /api/auth/telegram/browser
POST /api/me/telegram/link
DELETE /api/me/telegram/link
```

Browser login:

1. User clicks "Continue with Telegram" on `/login`.
2. Browser completes Telegram Login OIDC flow.
3. Frontend sends authorization result or code exchange result to backend, depending on final implementation shape.
4. Backend validates with Telegram OIDC rules.
5. If Telegram ID is linked, issue app JWT for the same user.
6. If Telegram ID is unknown, create Telegram provider user or ask user to register/link depending on product decision.

Browser account linking:

1. Existing email/password user opens `/profile`.
2. User clicks "Link Telegram".
3. Telegram OIDC completes.
4. Backend validates Telegram identity.
5. If Telegram ID is not linked elsewhere, attach it to current user.
6. If Telegram ID is already linked to a different user, return a conflict and do not merge automatically.

### 5.4 Sync Semantics

"Synchronized with browser auth" means:

- same `users._id`;
- same role, limits, history, selected deck/back, theme and language preferences;
- separate app JWTs per environment;
- no dependency on shared cookies or shared `localStorage`;
- deterministic account linking rules.

It does not mean:

- Telegram WebView can read the browser's `localStorage`;
- browser can automatically become logged in just because the Mini App has a token;
- existing email accounts can be merged silently without explicit user action.

---

## 6. Backend Work Plan

### 6.1 Configuration

Add backend env values:

```text
TELEGRAM_BOT_TOKEN=
TELEGRAM_BOT_ID=
TELEGRAM_LOGIN_CLIENT_ID=
TELEGRAM_LOGIN_CLIENT_SECRET=
TELEGRAM_MINIAPP_INITDATA_TTL=24h
TELEGRAM_OIDC_ISSUER=https://oauth.telegram.org
TELEGRAM_OIDC_JWKS_URL=https://oauth.telegram.org/.well-known/jwks.json
FRONTEND_URLS=https://example.com,https://staging.example.com
```

Notes:

- `TELEGRAM_BOT_TOKEN` and client secret must never be sent to frontend or logged.
- `FRONTEND_URLS` may be needed because browser site, staging, and Mini App URL can use more than one origin.
- Existing single `FRONTEND_URL` can be preserved for compatibility and expanded later.

### 6.2 Data Model And Indexes

Update `User` model:

- add Telegram identity fields;
- allow provider-only auth if product accepts Telegram-only signup;
- preserve existing email/password behavior;
- preserve existing response shape, adding Telegram linked status only where useful.

Indexes:

- existing unique `username`;
- partial unique `email` if email becomes optional;
- partial unique `telegram_id`;
- role/status indexes stay.

Migration/backfill:

- existing users have no `telegram_id`;
- existing email/password login must continue working;
- account linking attaches Telegram identity later.

### 6.3 Telegram Init Data Utility

Add a backend utility, likely under `backend/internal/utils/telegram.go`, that:

- accepts raw `init_data`;
- parses query string safely;
- builds data-check string from decoded fields;
- validates HMAC with bot token;
- uses constant-time comparison;
- parses Telegram user JSON;
- validates `auth_date`;
- returns a typed struct to the service layer.

Unit tests must cover:

- valid init data fixture;
- invalid hash;
- missing hash;
- missing user;
- expired `auth_date`;
- malformed user JSON;
- field ordering;
- special characters and URL encoding.

### 6.4 Auth Service Changes

Add a Telegram auth service path that:

- finds user by `telegram_id`;
- creates Telegram-only user when appropriate;
- links Telegram identity to current authenticated user;
- updates Telegram profile snapshot on login/link;
- issues the same app JWT through existing `utils.GenerateToken`;
- rejects disabled users consistently with existing login.

Conflict rules:

- if Telegram ID is already linked to another user, return `409 Conflict`;
- if current user is disabled, reject linking;
- if username derived from Telegram username is unavailable, generate a unique fallback;
- never overwrite an existing user's email/password when linking Telegram.

### 6.5 Handlers And Routes

Public routes:

```text
POST /api/auth/telegram/miniapp
POST /api/auth/telegram/browser
```

Protected routes:

```text
POST /api/me/telegram/link
DELETE /api/me/telegram/link
```

Optional admin visibility:

- show whether a user has Telegram linked;
- show Telegram username and last Telegram login;
- avoid exposing Telegram IDs broadly unless needed for support.

### 6.6 Security

Security requirements:

- validate Mini App `initData` only on backend;
- never trust `initDataUnsafe` for auth;
- verify `auth_date` freshness;
- use constant-time hash comparison;
- redact `init_data`, bot token, OIDC tokens and app JWTs from logs;
- verify OIDC tokens server-side;
- verify OIDC nonce/state to prevent replay/CSRF;
- return generic auth errors where possible;
- add tests for forged Telegram payloads;
- keep existing role middleware as the source of authorization.

Rate limiting is not currently present in the project. For this task, add it only if the project decides to introduce rate limiting globally; otherwise document as a deployment/API-gateway concern.

---

## 7. Frontend Work Plan

### 7.1 Telegram Runtime Adapter

Add a small adapter/composable, for example:

```text
frontend/src/composables/useTelegramMiniApp.js
frontend/src/services/telegramAuth.js
```

Responsibilities:

- detect Telegram runtime;
- expose `isTelegramMiniApp`;
- expose raw `initData`;
- call `WebApp.ready()` and `WebApp.expand()`;
- expose platform/version if useful for diagnostics;
- connect Telegram BackButton to Vue Router where appropriate;
- map Telegram color scheme/theme params only where it improves native feel.

Keep this adapter optional. Browser mode must work if `window.Telegram` is absent.

### 7.2 Startup Integration

Current `App.vue` already hydrates the auth store before rendering. Extend this flow later so that:

- in browser mode: current `hydrateSession()` behavior remains;
- in Telegram mode with stored app token: current `hydrateSession()` behavior remains;
- in Telegram mode without stored app token: exchange `initData` before redirecting to `/login`;
- failed Telegram auth shows a compact recoverable error instead of looping through `/login`.

The router guard should avoid forcing Mini App users to `/login` before Telegram auto-auth has had a chance to run.

### 7.3 Auth Store

Add actions without changing existing login/register contract:

- `loginWithTelegramMiniApp(initData, timezone)`;
- `loginWithTelegramBrowser(payload)` or equivalent OIDC callback handler;
- `linkTelegram(payload)`;
- `unlinkTelegram()`.

Keep token persistence centralized in the existing `persistToken()`.

### 7.4 UI Adaptation

The Mini App should show the same product, but not blindly duplicate browser chrome where Telegram already provides it.

Potential conditional adjustments:

- keep app navigation, but check whether sticky header consumes too much vertical space in Telegram WebView;
- hide or simplify `/login` and `/register` in Telegram mode because auth is automatic;
- keep profile, history, dashboard and admin flows unchanged;
- respect `safeAreaInset` and `contentSafeAreaInset`;
- avoid bottom-fixed controls that ignore Telegram bottom bar and safe area;
- use Telegram BackButton for nested routes/admin pages when it improves navigation;
- keep current theme selector, but consider initial theme from Telegram `colorScheme`;
- keep language selector, but use Telegram `language_code` as first-run default for Telegram-only users.

Do not introduce a separate Telegram-only design system. Any UI adjustment should be a conditional variant of the existing components.

### 7.5 Browser Telegram Login UI

Add browser-only Telegram login/link entry points:

- `/login`: "Continue with Telegram";
- `/profile`: "Link Telegram" / "Unlink Telegram";
- optional account state text: linked Telegram username.

Do not show browser Telegram login controls inside Mini App if Mini App auto-auth is already active.

---

## 8. Bot And Deployment Setup

### 8.1 BotFather Setup

Required setup:

- create or reuse a Telegram bot that represents Nebula Tarot;
- set bot name/logo consistent with the product;
- configure Main Mini App URL;
- configure menu button URL;
- register allowed URLs for Telegram Login OIDC;
- keep production and staging URLs explicit.

Recommended:

- separate dev/staging bot from production bot;
- separate production and staging allowed URLs;
- document bot username and non-secret IDs, but not token/secret.

### 8.2 HTTPS Requirements

Telegram Mini App URL must be HTTPS in real Telegram clients.

Local options:

- staging deployment;
- `ngrok`;
- `cloudflared tunnel`;
- another HTTPS reverse tunnel.

The backend public asset URLs must also be reachable through HTTPS. Current uploaded deck/card images are served from `backend/uploads`; production Mini App must not point to `localhost` asset URLs.

### 8.3 CORS

Current backend CORS accepts one frontend URL. Mini App/browser deployment will likely require:

- production browser origin;
- production Mini App origin if different;
- staging origin;
- local dev origin.

Plan to move from one `FRONTEND_URL` to a list of allowed origins, while keeping local defaults.

---

## 9. Implementation Phases

### Phase 0: Planning And Architecture

Status for this file: completed.

Deliverables:

- current project auth/frontend/backend baseline reviewed;
- Telegram Mini Apps and Telegram Login docs reviewed;
- integration plan documented;
- no implementation changes made.

### Phase 1: Backend Telegram Identity Foundation

Goal: backend can store, link and resolve Telegram identities.

Tasks:

- add config values;
- extend user model and indexes;
- add repository methods for `telegram_id`;
- add Telegram init data validation utility;
- add unit tests for validation and account resolution;
- keep email/password auth fully compatible.

Ready when:

- existing backend tests pass;
- forged/expired Telegram init data tests fail correctly;
- existing users remain readable after model/index changes.

### Phase 2: Mini App Token Exchange

Goal: Mini App can exchange Telegram `initData` for Nebula app JWT.

Tasks:

- add `POST /api/auth/telegram/miniapp`;
- issue existing app JWT for linked or newly-created Telegram user;
- update frontend auth store with Mini App login action;
- add Telegram runtime adapter;
- integrate startup auto-auth without breaking browser flow.

Ready when:

- Telegram Mini App opens dashboard without showing email login;
- same app API calls work with returned JWT;
- browser email/password login still works unchanged.

### Phase 3: Browser Telegram Login And Linking

Goal: browser auth can resolve to the same Telegram-linked account.

Tasks:

- implement Telegram Login OIDC server validation;
- add browser login button;
- add profile linking/unlinking flow;
- handle identity conflicts clearly;
- add backend and frontend tests/manual checks for linking.

Ready when:

- existing browser user can link Telegram from profile;
- after linking, Mini App opens the same account and shows the same history/settings;
- browser "Continue with Telegram" logs into the same linked user.

### Phase 4: Mini App UX Polish

Goal: same frontend feels correct inside Telegram WebView.

Tasks:

- tune safe-area and viewport CSS;
- check sticky header and mobile navigation in Telegram;
- wire Telegram BackButton for nested navigation if needed;
- use Telegram theme/language defaults for first-run Telegram users;
- verify card reveal, history, profile, admin and AI interpretation flows inside Telegram.

Ready when:

- no content is hidden behind Telegram controls;
- iOS, Android, Telegram Desktop and Telegram Web behave acceptably;
- current browser layout remains unchanged.

### Phase 5: Bot Setup, Rollout And Acceptance

Goal: make the feature operational.

Tasks:

- configure BotFather Main Mini App and menu button;
- configure Telegram Login allowed URLs;
- document non-secret bot identifiers and required env values;
- validate staging through real Telegram clients;
- run full browser regression;
- update project docs after implementation.

Ready when:

- Mini App launches from bot profile, menu button and direct link;
- production/staging HTTPS URLs work;
- account sync acceptance scenarios pass;
- secrets are not committed or logged.

---

## 10. Thin Points And Risks

### 10.1 Storage Is Not Sync

Do not rely on shared `localStorage` or cookies between browser and Telegram WebView. Treat all clients as separate sessions that authenticate to the same backend user.

### 10.2 Telegram Does Not Provide Email By Default

Telegram Mini App user data does not solve email identity matching. Existing email accounts need explicit linking. Silent merge by username is unsafe.

### 10.3 Two Telegram Auth Algorithms

Mini App `initData` validation and browser Telegram Login/OIDC validation are different flows.

- Mini App: HMAC over `initData` using bot token-derived key.
- Browser OIDC: validate signed `id_token` through Telegram JWKS and claims.

Keep them as separate utilities to avoid security mistakes.

### 10.4 `initDataUnsafe` Is Untrusted

Frontend can use `initDataUnsafe` for display hints, never for authorization decisions. Backend must accept raw `initData`.

### 10.5 Replay And Token Lifetime

`auth_date` should be checked. App JWT TTL remains controlled by existing `JWT_TTL`. If a Telegram token exchange is replayed after TTL, backend should reject it.

### 10.6 Existing Unique Email Index

If Telegram-only accounts are allowed, the current required unique email model needs careful migration. Use partial unique indexes rather than fake placeholder emails where possible.

### 10.7 Existing CORS Shape

Current CORS is single-origin oriented. Multi-origin deployment should be handled intentionally before staging/prod Telegram testing.

### 10.8 Browser Login Popup Policy

Telegram Login library depends on popup communication. Avoid strict `Cross-Origin-Opener-Policy: same-origin`; if COOP is used, it needs `same-origin-allow-popups` or equivalent.

### 10.9 WebView Layout

Telegram WebView height, safe area, keyboard and fullscreen behavior vary across clients. Test on real iOS and Android, not only browser emulator.

### 10.10 Uploaded Assets

Card images served from backend must resolve through public HTTPS in Mini App. Local `localhost` asset URLs will fail on real devices.

### 10.11 Account Conflict UX

If Telegram ID is already linked to user A and browser user B tries to link it, do not merge automatically. Return conflict and require manual support/product decision.

---

## 11. Criteria Of Readiness

Feature is ready only when all criteria below are met.

### Functional

- Existing browser email/password auth still works.
- Telegram Mini App auto-authenticates using validated `initData`.
- Browser Telegram Login authenticates through server-validated OIDC.
- Existing email user can link Telegram from profile.
- Linked user sees the same profile, settings, selected deck/back, readings and limits in browser and Mini App.
- Telegram-only new user can complete first reading without email/password registration.
- Disabled users cannot access through Telegram auth.
- Keeper-only admin routes remain protected.

### Security

- Bot token and Telegram client secret are backend-only.
- `initData` validation uses constant-time hash comparison.
- `auth_date` freshness is enforced.
- OIDC `id_token` signature and claims are validated.
- Forged, expired and malformed Telegram auth payloads are rejected.
- JWTs and Telegram auth payloads are not logged.

### UX

- Mini App opens from bot profile, menu button and direct link.
- Dashboard, history, profile and admin are usable in Telegram WebView.
- Mobile layout respects Telegram viewport and safe areas.
- Theme/language behavior remains predictable.
- Browser site remains visually unchanged unless user uses Telegram login controls.

### Operations

- Required env values are documented without secrets.
- BotFather setup is documented.
- Staging and production URLs are HTTPS.
- CORS allows only intended origins.
- Manual test checklist is updated after implementation.

---

## 12. Suggested Manual Acceptance Scenarios

1. New Telegram user opens Mini App and reaches dashboard.
2. Same Telegram user closes and reopens Mini App; session hydrates.
3. Telegram-only user opens browser site and logs in via Telegram; same user data appears.
4. Existing email user links Telegram from browser profile; Mini App opens same account.
5. Existing email user unlinks Telegram; Mini App no longer logs into that account.
6. Attempt to link Telegram ID already attached to another account returns conflict.
7. Disabled linked user cannot log in through Mini App or browser Telegram login.
8. Forged `initData` is rejected.
9. Expired `initData` is rejected.
10. Browser email/password login/register still pass.
11. Role limits and history retention are identical in browser and Mini App.
12. Card of the Day remains deterministic by user/timezone/accounting day.
13. 3-card spread reveal state persists in Mini App.
14. AI interpretation flow works inside Telegram if `v0.1.1` AI is enabled.
15. iOS Telegram, Android Telegram, Telegram Desktop and normal browser are checked.

---

## 13. Recommended Stack Additions

Backend:

- Go standard `crypto/hmac`, `crypto/sha256`, `crypto/subtle`, `net/url`;
- JWT/JWKS verification for Telegram OIDC, either through existing JWT library or a small focused dependency if needed;
- no frontend exposure of Telegram secrets.

Frontend:

- official Telegram Web App script: `https://telegram.org/js/telegram-web-app.js`;
- Telegram Login/OIDC integration for browser login;
- optional `@telegram-apps/sdk-vue` only if the project wants a wrapper around launch params, mocks and Vue composables. It is not required for a minimal implementation.

Do not introduce a separate frontend framework or bot framework for this task. Bot setup can be handled through BotFather and one-time Bot API calls unless future notification features require a real bot service.

---

## 14. Open Product Decisions

These should be answered before implementation:

1. Should Mini App first-time Telegram users be auto-created immediately, or should they see a short onboarding confirmation first?
2. Should existing email users be required to manually link Telegram before Mini App can access the same account? Recommended: yes.
3. Should browser Telegram Login create Telegram-only accounts automatically, or only log in already-linked users? Recommended: allow creation if Mini App also allows it.
4. Should unlinking Telegram be allowed when it is the user's only auth provider? Recommended: no, unless the user first adds email/password.
5. Should production and staging use separate Telegram bots? Recommended: yes.
6. Should admin UI expose Telegram identity fields? Recommended: limited visibility for support/debugging.

---

## 15. References

- Telegram Mini Apps overview and WebApp runtime: https://core.telegram.org/bots/webapps
- Telegram Bot API `WebAppInfo`: https://core.telegram.org/bots/api#webappinfo
- Telegram Bot API `setChatMenuButton`: https://core.telegram.org/bots/api#setchatmenubutton
- Telegram Mini Apps init data validation guide: https://docs.telegram-mini-apps.com/platform/init-data
- Telegram Login / OpenID Connect: https://core.telegram.org/bots/telegram-login
- Legacy Telegram Login Widget: https://core.telegram.org/widgets/login-legacy
