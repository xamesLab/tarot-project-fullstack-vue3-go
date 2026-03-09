# Implementation Plan: Frontend (Vue 3 + Pinia)

## 1. Objective
Build a mobile-first frontend for Tarot Daily MVP that supports:
- auth (register/login/logout)
- personal account with daily card + 2 spreads
- user settings (timezone, deck theme)
- role-based access to admin page

## 2. Scope
- In scope: user-facing UI, admin UI, API integration, route guards, state management.
- Out of scope: production deployment, advanced analytics, full CMS.

## 3. Assumptions
- Backend provides REST endpoints defined in `specs/main/PRD.md`.
- JWT-based auth is used.
- MongoDB runs on remote VPS and backend connects via SSH tunnel; frontend talks only to backend API.

## 4. Architecture and Structure
- Framework: Vue 3 (Composition API)
- State: Pinia
- Routing: Vue Router
- Suggested folders:
  - `src/pages` (Login, Register, Dashboard, Settings, Admin)
  - `src/components` (CardView, SpreadView, Header, ThemePicker)
  - `src/stores` (`auth`, `reading`, `settings`)
  - `src/services` (`apiClient`, `authApi`, `readingApi`, `adminApi`)
  - `src/router` (route map + guards)

## 5. Delivery Phases

### Phase 1: Project Setup
- Initialize Vue 3 app.
- Add Pinia and Vue Router.
- Configure environment variables:
  - `VITE_API_BASE_URL`
- Add basic layout and global styles (mobile-first).

### Phase 2: Auth Flow
- Implement pages:
  - `/login`
  - `/register`
- Add auth store:
  - token handling
  - current user profile (`/api/me`)
- Persist auth state (localStorage/sessionStorage).
- Add logout action.

### Phase 3: Routing and Access Control
- Public routes: login/register.
- Private routes: dashboard/settings.
- Admin-only route: `/admin`.
- Implement route guards:
  - redirect unauthenticated users to login
  - block non-admin users from `/admin` (redirect or 403 page)

### Phase 4: Dashboard (Daily Reading)
- Implement dashboard page with sections:
  - Card of the Day
  - Spread A (3 cards)
  - Spread B (5 cards)
- Integrate `GET /api/reading/today`.
- Render card fields:
  - image
  - name
  - state (`open`/`closed`)
  - description
  - interpretation
- Handle loading, empty state, and API errors.

### Phase 5: User Settings
- Implement settings page:
  - timezone selector (default `UTC+00:00`)
  - deck theme selector
- Integrate:
  - `GET /api/deck-themes`
  - `PATCH /api/me/settings`
- Update local state after settings save.

### Phase 6: Admin Page (MVP)
- Implement `/admin` page:
  - list available deck themes
  - add new theme form (MVP-level fields)
- Integrate:
  - `GET /api/admin/deck-themes`
  - `POST /api/admin/deck-themes`
- Add clear feedback for success/error and permission failures.

### Phase 7: UX and Resilience
- Improve responsive behavior for 360px+ width screens.
- Ensure tap targets >= 44px.
- Add skeleton/loading placeholders.
- Standardize API error messaging.

### Phase 8: Testing and QA
- Unit tests (store actions, guards, key components).
- Integration smoke tests for:
  - login -> dashboard
  - settings update
  - admin access control
- Manual mobile checks (Chrome DevTools devices).

## 6. API Contract Checklist
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/me`
- `PATCH /api/me/settings`
- `GET /api/reading/today`
- `GET /api/deck-themes`
- `GET /api/admin/deck-themes` (admin)
- `POST /api/admin/deck-themes` (admin)

## 7. Definition of Done (Frontend)
- All required pages implemented and routed.
- Auth + route guards working.
- Dashboard renders card states and interpretations correctly.
- Settings and deck theme selection persist.
- Admin page restricted to `admin` users and fully functional for MVP operations.
- UI is usable on mobile and desktop.

