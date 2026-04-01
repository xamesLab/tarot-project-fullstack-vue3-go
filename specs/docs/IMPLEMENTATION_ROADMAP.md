# Implementation Roadmap
## Tarot Daily - MVP v1.0

**Document Version:** 1.0  
**Last Updated:** April 1, 2026  
**Status:** Working Baseline

---

## Table of Contents

1. [Overview](#1-overview)
2. [Development Phases](#2-development-phases)
3. [Phase 0: Project Setup](#3-phase-0-project-setup)
4. [Phase 1: Core Backend](#4-phase-1-core-backend)
5. [Phase 2: Core Frontend](#5-phase-2-core-frontend)
6. [Phase 3: Reading Features](#6-phase-3-reading-features)
7. [Phase 4: Admin Panel](#7-phase-4-admin-panel)
8. [Phase 5: Polish & Testing](#8-phase-5-polish--testing)
9. [Phase 6: Deployment](#9-phase-6-deployment)
10. [Risk Assessment](#10-risk-assessment)
11. [Testing Strategy](#11-testing-strategy)

---

## 1. Overview

### 1.1 Project Timeline

**Estimated Duration:** 6-8 weeks for MVP

**Phases:**
- **Phase 0:** Project Setup (3-5 days)
- **Phase 1:** Core Backend (7-10 days)
- **Phase 2:** Core Frontend (7-10 days)
- **Phase 3:** Reading Features (7-10 days)
- **Phase 4:** Admin Panel (5-7 days)
- **Phase 5:** Polish & Testing (5-7 days)
- **Phase 6:** Deployment (2-3 days)

### 1.2 Development Approach

**Methodology:** Iterative development with continuous integration

**Principles:**
- Build vertically (full stack features one at a time)
- Test early and often (manual validation)
- UI Kit first for all frontend work
- Frontend implementation must reconcile with `specs/docs/DESIGN_SYSTEM.md`
- Document as you build
- Keep it simple (MVP scope only)

**Current note:** preparatory frontend design-system work and `phase-2-core-frontend` baseline have already been completed. The approved project-level reference now lives in `specs/docs/DESIGN_SYSTEM.md`, and the next implementation focus is `phase-3-reading-features`.

### 1.3 Definition of Done

**Feature Complete:**
- ✅ Code written and reviewed
- ✅ Manual testing passed
- ✅ Documentation updated
- ✅ Works in both EN and RU
- ✅ Responsive on mobile
- ✅ Light/Dark theme applied

---

## 2. Development Phases

### High-Level Phase Dependency

```
Phase 0: Project Setup
    ↓
Phase 1: Core Backend ←──────┐
    ↓                        │
Phase 2: Core Frontend       │ (parallel)
    ↓                        │
Phase 3: Reading Features ←──┘
    ↓
Phase 4: Admin Panel
    ↓
Phase 5: Polish & Testing
    ↓
Phase 6: Deployment
```

---

## 3. Phase 0: Project Setup

**Duration:** 3-5 days  
**Goal:** Establish development environment and project structure

### 3.0.1 Environment Setup

**Tasks:**
- [ ] Install Go 1.22+
- [ ] Install Node.js 20 LTS
- [ ] Setup MongoDB SSH tunnel
- [ ] Create MongoDB database `tarot_daily`
- [ ] Verify MongoDB connection

**Deliverable:** Working local development environment

---

### 3.0.2 Backend Project Initialization

**Tasks:**
- [ ] Create `backend/` directory
- [ ] Initialize Go module: `go mod init tarot-daily`
- [ ] Create directory structure:
  ```
  backend/
  ├── main.go
  ├── config/
  ├── internal/
  │   ├── middleware/
  │   ├── handlers/
  │   ├── services/
  │   ├── repositories/
  │   ├── models/
  │   └── utils/
  └── uploads/
  ```
- [ ] Install dependencies:
  - Gin
  - MongoDB driver
  - JWT library
  - Validator
  - godotenv
- [ ] Create `.env` file with development config
- [ ] Create `.gitignore` (exclude `uploads/`, `.env`)
- [ ] Setup basic Gin server (Hello World)
- [ ] Verify server runs on `localhost:8080`

**Deliverable:** Backend scaffold running

---

### 3.0.3 Frontend Project Initialization

**Tasks:**
- [ ] Create `frontend/` directory
- [ ] Initialize Vite + Vue 3 project: `npm create vite@latest`
- [ ] Create directory structure:
  ```
  frontend/src/
  ├── main.js
  ├── App.vue
  ├── router/
  ├── stores/
  ├── services/
  ├── pages/
  ├── components/
  ├── composables/
  ├── locales/
  ├── styles/
  └── utils/
  ```
- [ ] Install dependencies:
  - Vue Router
  - Pinia
  - Axios
  - Vue I18n
  - VeeValidate
  - date-fns
  - SCSS
- [ ] Create `.env` file with `VITE_API_URL`
- [ ] Create `.gitignore` (exclude `dist/`, `node_modules/`)
- [ ] Setup Vite dev server
- [ ] Verify frontend runs on `localhost:5173`

**Deliverable:** Frontend scaffold running

---

### 3.0.4 UI Kit Page Setup

**Tasks:**
- [ ] Create `frontend/src/pages/dev/UIKit.vue`
- [ ] Setup route: `/ui-kit` (dev mode only)
- [ ] Create empty sections:
  - Typography
  - Colors
  - Buttons
  - Inputs
  - Cards
  - Modals
- [ ] Create SCSS theme files:
  - `styles/themes/_light.scss`
  - `styles/themes/_dark.scss`
  - `styles/variables.scss`
- [ ] Setup CSS variable system
- [ ] Implement theme toggle component
- [ ] Test theme switching

**Deliverable:** Empty UI Kit page accessible at `/ui-kit` in dev mode

**Note:** UI Kit will be populated with components as they're built (separate instruction later)

---

### 3.0.5 Database Schema Setup

**Tasks:**
- [ ] Create seed script: `deploy/seed/main.go`
- [ ] Load card data from `tarot-cards.json`
- [ ] Transform JSON structure to match MongoDB schema
- [ ] Insert 78 cards into `cards` collection
- [ ] Create indexes:
  - `cards`: `{ key: 1 }` (unique)
  - `cards`: `{ arcana: 1 }`
- [ ] Create default role policies in `role_policies`:
  - seeker: 3 spreads/day, 1 day retention
  - adept: 5 spreads/day, 7 days retention
  - master: 10 spreads/day, 15 days retention
  - keeper: unlimited, 30 days retention
- [ ] Create admin settings document in `admin_settings`:
  - maintenance_mode: false
  - default_role: "seeker"
  - registration_enabled: true
- [ ] Create default deck document in `decks`:
  - id: "default"
  - name: { en: "Tarot Classic", ru: "Таро Классика" }
  - status: "done"
- [ ] Run seed script and verify data

**Deliverable:** Database seeded with all static data

---

### 3.0.6 i18n Setup

**Tasks:**
- [ ] Create `frontend/src/locales/en.json` (empty structure)
- [ ] Create `frontend/src/locales/ru.json` (empty structure)
- [ ] Setup Vue I18n plugin in `main.js`
- [ ] Create language switcher component
- [ ] Test switching between EN/RU

**Deliverable:** i18n system working, translations added incrementally

---

**Phase 0 Exit Criteria:**
- ✅ Backend running on `localhost:8080`
- ✅ Frontend running on `localhost:5173`
- ✅ MongoDB connected and seeded
- ✅ UI Kit page accessible
- ✅ Theme system working
- ✅ i18n system working

**Estimated Time:** 3-5 days

---

## 4. Phase 1: Core Backend

**Duration:** 7-10 days  
**Goal:** Build authentication, user management, and core API infrastructure

### 4.1.1 Authentication System

**Tasks:**

**Models:**
- [ ] Create `models/user.go`:
  - User struct with all fields
  - Validation tags
- [ ] Create `models/auth.go`:
  - RegisterRequest
  - LoginRequest
  - JWTClaims

**Utilities:**
- [ ] Create `utils/jwt.go`:
  - GenerateToken(userID, role, timezone)
  - VerifyToken(tokenString)
  - ExtractClaims(token)
- [ ] Create `utils/hash.go`:
  - HashPassword(password)
  - VerifyPassword(hash, password)
- [ ] Create `utils/validator.go`:
  - Custom validators (password strength, timezone)

**Repository:**
- [ ] Create `repositories/user.go`:
  - Create(user)
  - FindByEmail(email)
  - FindByID(id)
  - Update(user)
  - FindByUsername(username)

**Service:**
- [ ] Create `services/auth.go`:
  - Register(registerReq) → (user, token, error)
  - Login(loginReq) → (user, token, error)
  - LoadUser(userID) → (user, error)

**Handlers:**
- [ ] Create `handlers/auth.go`:
  - POST /api/auth/register
  - POST /api/auth/login

**Middleware:**
- [ ] Create `middleware/auth.go`:
  - AuthMiddleware() - verify JWT
  - Attach user info to context
- [ ] Create `middleware/role.go`:
  - RequireRole(...roles) - check user role
- [ ] Create `middleware/cors.go`:
  - CORS configuration

**Testing:**
- [ ] Test registration with valid data
- [ ] Test registration with duplicate email
- [ ] Test registration with weak password
- [ ] Test login with valid credentials
- [ ] Test login with invalid credentials
- [ ] Test JWT generation and verification
- [ ] Test middleware auth rejection

**Deliverable:** Working auth endpoints with JWT

---

### 4.1.2 User Profile Management

**Tasks:**

**Repository:**
- [ ] Extend `repositories/user.go`:
  - UpdateProfile(userID, updates)
  - UpdateSettings(userID, settings)

**Service:**
- [ ] Create `services/user.go`:
  - GetUserProfile(userID)
  - UpdateProfile(userID, updates)
  - UpdateSettings(userID, settings)

**Handlers:**
- [ ] Create `handlers/user.go`:
  - GET /api/me
  - PATCH /api/me/profile
  - PATCH /api/me/settings

**Testing:**
- [ ] Test GET /api/me (authenticated)
- [ ] Test GET /api/me (unauthenticated → 401)
- [ ] Test profile update (username, timezone)
- [ ] Test settings update (language, theme, deck, back)
- [ ] Test username uniqueness validation

**Deliverable:** User profile endpoints working

---

### 4.1.3 Card Repository

**Tasks:**

**Repository:**
- [ ] Create `repositories/card.go`:
  - FindByKey(key)
  - FindAll()
  - FindByArcana(arcana)
  - GetRandomCards(count, excludeKeys)

**Utilities:**
- [ ] Create `utils/random.go`:
  - DeterministicRandom(seed) → random generator
  - SelectRandomCards(seed, count) → []cardKey

**Testing:**
- [ ] Test FindByKey for all 78 cards
- [ ] Test GetRandomCards returns unique cards
- [ ] Test deterministic random (same seed → same cards)

**Deliverable:** Card data access layer

---

### 4.1.4 Time & Timezone Utilities

**Tasks:**

**Utilities:**
- [ ] Create `utils/time.go`:
  - GetAccountingDate(timezone) → string (YYYY-MM-DD)
  - GetNextResetTime(timezone) → time.Time
  - ValidateTimezone(tz) → bool

**Logic:**
```go
func GetAccountingDate(timezone string) (string, error) {
    loc, _ := time.LoadLocation(timezone)
    now := time.Now().In(loc)
    
    // If before 3:00 AM, it's still "yesterday"
    if now.Hour() < 3 {
        now = now.Add(-24 * time.Hour)
    }
    
    return now.Format("2006-01-02"), nil
}
```

**Testing:**
- [ ] Test GetAccountingDate at 02:59 AM → returns yesterday
- [ ] Test GetAccountingDate at 03:00 AM → returns today
- [ ] Test GetAccountingDate at 11:59 PM → returns today
- [ ] Test GetNextResetTime calculation
- [ ] Test with various timezones (UTC, Europe/Amsterdam, America/New_York)

**Deliverable:** Timezone utilities working correctly

---

**Phase 1 Exit Criteria:**
- ✅ User registration working
- ✅ User login working
- ✅ JWT auth middleware working
- ✅ Profile/settings update working
- ✅ Card repository functional
- ✅ Timezone calculations correct
- ✅ All Phase 1 tests passing

**Estimated Time:** 7-10 days

---

## 5. Phase 2: Core Frontend

**Duration:** 7-10 days  
**Goal:** Build UI foundation, authentication flows, and layout

**Frontend implementation baseline:**
- use `specs/docs/DESIGN_SYSTEM.md` as the primary design and component reference
- treat `frontend/src/pages/dev/UIKit.vue` as the visual showcase/reference implementation
- do not introduce a second parallel UI language outside the design-system layer

### 5.2.1 API Service Layer

**Tasks:**

**API Client:**
- [ ] Create `services/api.js`:
  - Axios instance with base URL
  - Request interceptor (add JWT token)
  - Response interceptor (handle 401, 403)
  - Error handling utilities

**Auth Service:**
- [ ] Create `services/auth.js`:
  - register(data)
  - login(credentials)
  - logout()
  - getCurrentUser()

**User Service:**
- [ ] Create `services/user.js`:
  - getProfile()
  - updateProfile(data)
  - updateSettings(settings)

**Testing:**
- [ ] Test API client creation
- [ ] Test token injection
- [ ] Test 401 handling (logout + redirect)

**Deliverable:** API service layer ready

---

### 5.2.2 Pinia Stores

**Tasks:**

**Auth Store:**
- [ ] Create `stores/auth.js`:
  - State: user, token, isAuthenticated
  - Actions: register, login, logout, loadUser
  - Getters: userRole, userTimezone

**App Store:**
- [ ] Create `stores/app.js`:
  - State: theme, language, loading
  - Actions: setTheme, setLanguage, toggleTheme
  - Persist theme/language to localStorage

**Reading Store (placeholder):**
- [ ] Create `stores/readings.js`:
  - State: cardOfDay, spreads, history
  - Actions: (to be implemented in Phase 3)

**Testing:**
- [ ] Test auth store registration flow
- [ ] Test auth store login flow
- [ ] Test auth store logout
- [ ] Test app store theme switching
- [ ] Test app store language switching
- [ ] Test localStorage persistence

**Deliverable:** Pinia stores working

---

### 5.2.3 Router Setup

**Tasks:**

**Routes:**
- [ ] Create `router/index.js`:
  - Public routes: `/login`, `/register`
  - Protected routes: `/`, `/profile`, `/history`
  - Admin routes: `/admin/*`
  - Dev route: `/ui-kit` (dev only)

**Guards:**
- [ ] Create auth guard:
  - Check if authenticated
  - Redirect to `/login` if not
- [ ] Create admin guard:
  - Check if role === 'keeper'
  - Redirect to `/` if not

**Testing:**
- [ ] Test unauthenticated redirect to `/login`
- [ ] Test authenticated access to `/`
- [ ] Test admin-only route protection
- [ ] Test `/ui-kit` only in dev mode

**Deliverable:** Router with guards working

---

### 5.2.4 UI Components (Basic Set)

**Tasks:**

**Base Components:**
- [ ] Create `components/base/BaseButton.vue`:
  - variants from design system (`primary`, `secondary`, `ghost`, `outlined`, `text`)
  - size props and disabled/loading states
  - add to UI Kit showcase

- [ ] Create `components/base/BaseInput.vue`:
  - `modelValue`, `type`, `placeholder`, `error`, `disabled`
  - support field-level error state
  - add to UI Kit showcase

- [ ] Create `components/base/BaseSelect.vue`:
  - `options`, `modelValue`, `placeholder`, `error`
  - add to UI Kit showcase

- [ ] Create `components/base/BaseSurface.vue`:
  - reusable surface wrapper for cards, auth blocks, dashboard sections
  - align borders/backgrounds with design system

- [ ] Create `components/base/BaseSegmentedToggle.vue`:
  - use for `EN/RU` and `Light/Dark`
  - adopt final toggle pattern from UI Kit

- [ ] Create `components/base/BaseSwitchToggle.vue`:
  - use for binary settings controls
  - add to UI Kit showcase

- [ ] Create `components/base/BaseTabs.vue`:
  - variants: `underline`, `pills`, `boxed`
  - add to UI Kit showcase

- [ ] Create `components/base/BaseModal.vue`:
  - props: `isOpen`, `title`
  - slots: content, footer
  - add to UI Kit showcase

**Layout Components:**
- [ ] Create `components/layout/AppHeader.vue`:
  - Brand
  - segmented theme/language controls
  - user/auth actions

- [ ] Create `components/layout/AppFooter.vue`:
  - ornamental divider
  - supporting copy / meta text

- [ ] Create `components/layout/PageContainer.vue`
- [ ] Create `components/layout/PageSection.vue`

**Domain Components:**
- [ ] Create `components/domain/UserChip.vue`
- [ ] Create `components/domain/TarotCardPreview.vue`
- [ ] Create `components/domain/ReadingCard.vue`

**Composables:**
- [ ] Create `composables/useTheme.js`:
  - Get current theme
  - Toggle theme
  - Apply theme to DOM

- [ ] Create `composables/useLocale.js`:
  - Get current locale
  - Toggle EN/RU
  - Persist locale

- [ ] Create `composables/useAuth.js`:
  - Get current user
  - Check authentication
  - Logout

**Deliverable:** Base design-system components extracted from UI Kit and ready for page reuse

---

### 5.2.5 Authentication Pages

**Tasks:**

**Login Page:**
- [ ] Create `pages/Login.vue`:
  - Email input (with validation)
  - Password input
  - Submit button
  - Link to Register
  - Error message display
  - Loading state
  - Redirect to `/` on success

**Register Page:**
- [ ] Create `pages/Register.vue`:
  - Email input (with validation)
  - Username input (with validation)
  - Password input (with validation)
  - Timezone selector (auto-detected, editable)
  - Submit button
  - Link to Login
  - Error message display
  - Loading state
  - Redirect to `/` on success

**Validation:**
- [ ] Setup VeeValidate schemas
- [ ] Real-time validation
- [ ] Display field errors
- [ ] Disable submit if form invalid

**i18n:**
- [ ] Add auth translations to `locales/en.json`:
  - auth.login, auth.register, auth.email, auth.password, etc.
- [ ] Add auth translations to `locales/ru.json`

**Testing:**
- [ ] Test registration flow (valid data)
- [ ] Test registration with duplicate email
- [ ] Test login flow (valid credentials)
- [ ] Test login with invalid credentials
- [ ] Test validation errors display
- [ ] Test redirect after success
- [ ] Test both EN and RU

**Deliverable:** Working login/register pages

---

### 5.2.6 Dashboard Page (Skeleton)

**Tasks:**
- [ ] Create `pages/Dashboard.vue`:
  - Header: "Welcome, {username}"
  - Placeholder for Card of the Day
  - Placeholder for 3-Card Spread
  - Loading state
  - Empty state ("Get your first reading")

**Note:** Full implementation in Phase 3

**Deliverable:** Dashboard skeleton page

---

**Phase 2 Exit Criteria:**
- ✅ Login/Register pages working end-to-end
- ✅ JWT token stored and sent with requests
- ✅ Router guards protecting routes
- ✅ Theme toggle working
- ✅ Language switcher working
- ✅ Base design-system components extracted from UI Kit
- ✅ Dashboard skeleton visible

**Estimated Time:** 7-10 days

---

## 6. Phase 3: Reading Features

**Duration:** 7-10 days  
**Goal:** Implement Card of the Day and 3-Card Spread features

### 6.3.1 Backend: Reading Generation

**Tasks:**

**Models:**
- [ ] Create `models/reading.go`:
  - DailyReading struct
  - CardInReading struct
  - ReadingResponse struct

**Repository:**
- [ ] Create `repositories/reading.go`:
  - Create(reading)
  - FindByUserAndDate(userID, accountingDate, type)
  - FindByUser(userID, filters, pagination)
  - Delete old readings (for TTL)

**Service:**
- [ ] Create `services/reading.go`:
  - GenerateCardOfDay(userID, timezone) → (reading, error)
    - Calculate accounting date
    - Check if already generated
    - If not: generate deterministic seed, select card
    - Create reading record
    - Return reading with card data
  - GenerateThreeCardSpread(userID, timezone) → (reading, error)
    - Calculate accounting date
    - Check daily limit against role policy
    - If under limit: select 3 random cards
    - Create reading record
    - Return reading with card data
  - GetReadingHistory(userID, filters, pagination) → (readings, error)
    - Respect role retention policy
    - Return readings sorted by date desc

**Handlers:**
- [ ] Create `handlers/reading.go`:
  - POST /api/reading/generate/card-of-day
  - POST /api/reading/generate/spread-a (3-card)
  - GET /api/reading/history

**Utilities:**
- [ ] Create deterministic card selection:
  ```go
  func SelectCardOfDay(userID, accountingDate, timezone string) string {
      seed := fmt.Sprintf("%s-%s-%s", userID, accountingDate, timezone)
      hash := sha256.Sum256([]byte(seed))
      rng := rand.New(rand.NewSource(int64(binary.BigEndian.Uint64(hash[:]))))
      cardIndex := rng.Intn(78)
      return allCards[cardIndex].Key
  }
  ```

**Testing:**
- [ ] Test Card of Day generation (first time)
- [ ] Test Card of Day already generated (error)
- [ ] Test deterministic card (same seed → same card)
- [ ] Test 3-card spread generation (unique cards)
- [ ] Test 3-card spread limit enforcement (Seeker: 3/day)
- [ ] Test limit reset after accounting day change
- [ ] Test history retrieval (respects retention)
- [ ] Test timezone edge cases (2:59 AM vs 3:00 AM)

**Deliverable:** Reading generation endpoints working

---

### 6.3.2 Backend: Role Policies

**Tasks:**

**Repository:**
- [ ] Create `repositories/policy.go`:
  - FindByRole(role)
  - UpdatePolicy(role, updates)

**Service:**
- [ ] Create `services/policy.go`:
  - GetUserLimits(userID) → limits
  - CheckSpreadLimit(userID, accountingDate) → (allowed, remaining)

**Integration:**
- [ ] Integrate limit check into GenerateThreeCardSpread
- [ ] Return remaining count in response

**Testing:**
- [ ] Test limit check for each role
- [ ] Test remaining count calculation
- [ ] Test limit enforcement (reject when limit reached)

**Deliverable:** Role-based limits working

---

### 6.3.3 Frontend: Tarot Card Components

**Tasks:**

**TarotCard Component:**
- [ ] Create `components/tarot/TarotCard.vue`:
  - Props: card (object), size (small/medium/large)
  - Display card image
  - Display card name (localized)
  - Flip animation (future enhancement)
  - Responsive sizing
  - Add to UI Kit page

**CardOfDay Component:**
- [ ] Create `components/tarot/CardOfDay.vue`:
  - Display single tarot card (large)
  - Display card interpretation (localized)
  - Display tags/keywords
  - Show "Next card in X hours" countdown
  - "Generate" button (if not generated today)
  - Loading state
  - Error state

**ThreeCardSpread Component:**
- [ ] Create `components/tarot/ThreeCardSpread.vue`:
  - Display 3 cards in row (medium size)
  - Display interpretations below each card
  - Show remaining spreads count
  - "Generate New Spread" button
  - Disable button when limit reached
  - Loading state
  - Error state

**Testing:**
- [ ] Test card display with real card data
- [ ] Test localized names/interpretations
- [ ] Test responsive layout (mobile/desktop)
- [ ] Test loading states
- [ ] Test error states

**Deliverable:** Tarot card components in UI Kit

---

### 6.3.4 Frontend: Reading Services & Store

**Tasks:**

**Reading Service:**
- [ ] Create `services/reading.js`:
  - generateCardOfDay(timezone)
  - generateThreeCardSpread(timezone)
  - getReadingHistory(filters, page)

**Reading Store:**
- [ ] Extend `stores/readings.js`:
  - State:
    - cardOfDay: null
    - cardOfDayNextReset: null
    - threeCardSpreads: []
    - spreadLimits: { count: 0, limit: 0 }
    - history: []
    - loading: false
    - error: null
  - Actions:
    - fetchCardOfDay()
    - generateCardOfDay()
    - generateThreeCardSpread()
    - fetchHistory(filters)
  - Getters:
    - canGenerateCardOfDay
    - canGenerateSpread
    - hasCardOfDayToday

**Testing:**
- [ ] Test store actions
- [ ] Test state updates
- [ ] Test error handling

**Deliverable:** Reading store integrated

---

### 6.3.5 Frontend: Dashboard Implementation

**Tasks:**

**Dashboard Page:**
- [ ] Implement `pages/Dashboard.vue`:
  - Load card of day on mount (if available)
  - Show CardOfDay component
  - Show ThreeCardSpread component
  - Show spread limits
  - Handle loading/error states
  - Handle empty state (first-time user)

**User Flow:**
1. User lands on dashboard
2. If card of day exists today → show it
3. If not → show "Get your Card of the Day" button
4. User clicks → generate card → display card
5. User can generate 3-card spreads (within limits)
6. Limits displayed: "Spreads remaining: 2/5"

**i18n:**
- [ ] Add reading translations:
  - readings.cardOfDay, readings.threeCard, readings.generate, etc.

**Testing:**
- [ ] Test dashboard with card of day available
- [ ] Test dashboard with card of day not available
- [ ] Test card generation flow
- [ ] Test spread generation flow
- [ ] Test limit enforcement (button disabled)
- [ ] Test error handling
- [ ] Test both EN and RU

**Deliverable:** Fully functional Dashboard

---

### 6.3.6 Frontend: History Page

**Tasks:**

**ReadingHistory Component:**
- [ ] Create `components/tarot/ReadingHistory.vue`:
  - List of past readings (grouped by date)
  - Filter by type (all, card_of_day, three_card)
  - Pagination (10 per page)
  - Click reading → expand to show details
  - Show retention info ("Your history is kept for 7 days")

**History Page:**
- [ ] Create `pages/History.vue`:
  - Use ReadingHistory component
  - Load history on mount
  - Handle loading/error/empty states

**i18n:**
- [ ] Add history translations

**Testing:**
- [ ] Test history display
- [ ] Test filtering
- [ ] Test pagination
- [ ] Test expand/collapse
- [ ] Test retention info display

**Deliverable:** History page working

---

**Phase 3 Exit Criteria:**
- ✅ Card of Day generation working (deterministic)
- ✅ 3-Card Spread generation working (random)
- ✅ Role-based limits enforced
- ✅ Dashboard displays readings correctly
- ✅ History page shows past readings
- ✅ All reading flows tested end-to-end
- ✅ Works in both EN and RU

**Estimated Time:** 7-10 days

---

## 7. Phase 4: Admin Panel

**Duration:** 5-7 days  
**Goal:** Build admin panel for user/deck/policy management

### 7.4.1 Backend: Admin Endpoints

**Tasks:**

**Admin Service:**
- [ ] Create `services/admin.go`:
  - GetDashboardMetrics()
  - ListUsers(search, filters, pagination)
  - UpdateUserRole(userID, role)
  - DisableUser(userID)
  - GetAllRolePolicies()
  - UpdateRolePolicy(role, updates)
  - GetGlobalSettings()
  - UpdateGlobalSettings(updates)

**Admin Repository:**
- [ ] Extend repositories for admin operations:
  - User aggregation queries (count by role, status)
  - Reading aggregation (count by type, date range)
  - Deck queries

**Admin Handlers:**
- [ ] Create `handlers/admin.go`:
  - GET /api/admin/access-check
  - GET /api/admin/dashboard
  - GET /api/admin/users
  - PATCH /api/admin/users/:userId/role
  - GET /api/admin/roles
  - PATCH /api/admin/roles/:role
  - GET/PATCH /api/admin/settings

**Middleware:**
- [ ] Apply RequireRole("keeper") to all admin routes

**Testing:**
- [ ] Test admin access with keeper role
- [ ] Test admin access rejection (non-keeper)
- [ ] Test dashboard metrics calculation
- [ ] Test user list with search/filter
- [ ] Test role change
- [ ] Test policy updates
- [ ] Test settings updates

**Deliverable:** Admin API endpoints working

---

### 7.4.2 Backend: Deck Management

**Tasks:**

**Deck Service:**
- [ ] Create `services/deck.go`:
  - CreateDeck(metadata) → deck
  - UploadDeckZIP(deckID, zipFile) → error
  - ProcessDeckUpload(deckID) → error (async)
  - GetDeckDetails(deckID) → deck
  - ListDecks(filters) → decks
  - UpdateDeckMetadata(deckID, updates)
  - DisableDeck(deckID)

**Deck Repository:**
- [ ] Create `repositories/deck.go`:
  - Create(deck)
  - FindByID(deckID)
  - FindAll(filters)
  - Update(deckID, updates)

**File Upload Handlers:**
- [ ] Create `handlers/deck.go`:
  - POST /api/admin/decks (create + upload ZIP)
  - GET /api/admin/decks
  - GET /api/admin/decks/:deckId
  - PATCH /api/admin/decks/:deckId (metadata)

**ZIP Processing:**
- [ ] Create `utils/deck_upload.go`:
  - ValidateZIPStructure(zipPath) → error
  - ExtractDeckZIP(zipPath, destPath) → error
  - ValidateCardImages(dirPath) → error
    - Check 78 cards present
    - Validate file formats (PNG/JPG)
    - Validate file sizes (<2MB)
    - Validate aspect ratio (2:3 ±5%)

**Flow:**
1. Admin uploads ZIP
2. Backend saves to `/uploads/temp/{session}/`
3. Backend validates ZIP structure
4. If valid: extract to `/uploads/decks/{deck_id}/`
5. Update deck status: "processing" → "done"
6. Delete temp files

**Testing:**
- [ ] Test deck creation
- [ ] Test ZIP upload (valid structure)
- [ ] Test ZIP validation (missing cards → error)
- [ ] Test card image validation (wrong format → error)
- [ ] Test deck listing
- [ ] Test deck detail retrieval

**Deliverable:** Deck upload system working

---

### 7.4.3 Frontend: Admin Layout & Navigation

**Tasks:**

**Admin Layout:**
- [ ] Create `components/admin/AdminLayout.vue`:
  - Sidebar navigation:
    - Dashboard
    - Users
    - Decks
    - Role Policies
    - Settings
  - Main content area
  - Breadcrumbs

**Admin Guard:**
- [ ] Update router guards to check role === 'keeper'
- [ ] Show admin link in header (if keeper)

**i18n:**
- [ ] Add admin translations

**Testing:**
- [ ] Test admin layout display
- [ ] Test navigation between admin pages
- [ ] Test access control (non-keeper redirect)

**Deliverable:** Admin layout with navigation

---

### 7.4.4 Frontend: Admin Pages

**Tasks:**

**AdminDashboard Page:**
- [ ] Create `pages/admin/AdminDashboard.vue`:
  - Display metrics:
    - Total users (by role, by status)
    - Readings today/this week
    - Decks (by status)
  - Use cards/stats components
  - Auto-refresh every 30s (optional)

**Users Page:**
- [ ] Create `pages/admin/Users.vue`:
  - User list table:
    - Username, Email, Role, Status, Registered, Actions
  - Search bar (username/email)
  - Filter by role, status
  - Sort by date
  - Pagination
  - Actions:
    - Edit Role (modal with select)
    - Disable Account (confirmation modal)

**RolePolicies Page:**
- [ ] Create `pages/admin/RolePolicies.vue`:
  - List all 4 roles
  - Editable fields:
    - Three-card spread limit
    - History retention days
  - Save button
  - Confirmation: "This will affect all users with this role"

**Settings Page:**
- [ ] Create `pages/admin/Settings.vue`:
  - Global settings form:
    - Maintenance mode (toggle)
    - Maintenance message (textarea, i18n)
    - Default role (select)
    - Registration enabled (toggle)
  - Save button

**Decks Page:**
- [ ] Create `pages/admin/Decks.vue`:
  - Deck list:
    - Name, Status, Cards, Actions
  - "Create New Deck" button
  - Create modal:
    - Name (EN/RU)
    - Description (EN/RU)
    - Author
    - Tags
    - Upload ZIP file
  - Processing indicator
  - Deck detail view (click deck)

**Admin Services:**
- [ ] Create `services/admin.js`:
  - getDashboardMetrics()
  - getUsers(search, filters, page)
  - updateUserRole(userId, role)
  - getRolePolicies()
  - updateRolePolicy(role, updates)
  - getSettings()
  - updateSettings(updates)
  - getDecks()
  - createDeck(formData)
  - getDeckDetails(deckId)

**Admin Store:**
- [ ] Create `stores/admin.js`:
  - State: metrics, users, policies, settings, decks
  - Actions: fetch/update for each

**Testing:**
- [ ] Test all admin pages display correctly
- [ ] Test user search/filter
- [ ] Test role change flow
- [ ] Test policy update flow
- [ ] Test settings update flow
- [ ] Test deck creation flow
- [ ] Test both EN and RU

**Deliverable:** Full admin panel functional

---

**Phase 4 Exit Criteria:**
- ✅ Admin dashboard shows metrics
- ✅ User management working (role change, disable)
- ✅ Role policy editing working
- ✅ Global settings editable
- ✅ Deck upload working (ZIP validation)
- ✅ All admin pages tested
- ✅ Admin access restricted to keeper role

**Estimated Time:** 5-7 days

---

## 8. Phase 5: Polish & Testing

**Duration:** 5-7 days  
**Goal:** UI polish, bug fixes, performance optimization, comprehensive testing

### 8.5.1 UI Kit Completion

**Tasks:**
- [ ] Review all components in UI Kit
- [ ] Ensure consistency:
  - Colors (theme variables)
  - Typography (font sizes, weights)
  - Spacing (margins, paddings)
  - Border radius
  - Shadows
- [ ] Add missing components to UI Kit
- [ ] Document component usage in UI Kit page
- [ ] Test all components in both themes
- [ ] Test all components in both languages

**Deliverable:** Complete UI Kit documentation

---

### 8.5.2 Responsive Design

**Tasks:**
- [ ] Test all pages on mobile (375px width)
- [ ] Test all pages on tablet (768px width)
- [ ] Test all pages on desktop (1280px+ width)
- [ ] Fix layout issues:
  - Dashboard: stack cards vertically on mobile
  - History: single column on mobile
  - Admin tables: horizontal scroll on mobile
  - Modals: full-screen on mobile
- [ ] Test touch interactions (buttons, inputs)
- [ ] Test navigation menu on mobile (hamburger?)

**Devices to Test:**
- iOS Safari (iPhone)
- Android Chrome
- Desktop Chrome
- Desktop Safari

**Deliverable:** Fully responsive on all screen sizes

---

### 8.5.3 Accessibility

**Tasks:**
- [ ] Add ARIA labels to interactive elements
- [ ] Ensure keyboard navigation works:
  - Tab through forms
  - Enter to submit
  - Escape to close modals
- [ ] Test with screen reader (basic)
- [ ] Ensure color contrast meets WCAG AA:
  - Text on backgrounds
  - Button states
- [ ] Add focus indicators
- [ ] Test form validation messages are announced

**Deliverable:** Basic accessibility compliance

---

### 8.5.4 Performance Optimization

**Tasks:**

**Frontend:**
- [ ] Optimize images (compress deck images)
- [ ] Lazy load images (use `loading="lazy"`)
- [ ] Code splitting:
  - Admin pages in separate chunk
  - UI Kit in dev-only chunk
- [ ] Minify CSS/JS in production build
- [ ] Test page load times:
  - Dashboard: <2s
  - Login/Register: <1s

**Backend:**
- [ ] Add database indexes (already done in schema)
- [ ] Optimize queries (avoid N+1)
- [ ] Add request logging (Gin logger)
- [ ] Test API response times:
  - GET /api/me: <200ms
  - POST /api/reading/generate/*: <500ms

**Deliverable:** Performance targets met

---

### 8.5.5 Error Handling & Edge Cases

**Tasks:**

**Frontend:**
- [ ] Test network errors (disconnect during request)
- [ ] Test timeout errors (slow API)
- [ ] Test 500 errors from backend
- [ ] Ensure all error messages are user-friendly
- [ ] Ensure all error messages are localized

**Backend:**
- [ ] Test invalid JWT token
- [ ] Test expired JWT token
- [ ] Test malformed requests (invalid JSON)
- [ ] Test database connection errors
- [ ] Add error logging (log to file/console)

**Edge Cases:**
- [ ] User changes timezone while card of day is active
- [ ] User hits spread limit exactly
- [ ] Admin uploads invalid ZIP (wrong structure)
- [ ] Admin uploads ZIP with wrong file formats
- [ ] User tries to access deleted deck
- [ ] Clock rollback (accounting day calculation)

**Deliverable:** Robust error handling

---

### 8.5.6 Manual Testing Checklist

**User Flows:**
- [ ] **New User Journey:**
  1. Register account
  2. Auto-login to dashboard
  3. Generate first Card of the Day
  4. Generate 3-card spread
  5. View history
  6. Update profile (username, timezone, language, theme)
  7. Logout
  8. Login again
  9. Verify card of day not regenerated
  10. Hit spread limit

- [ ] **Returning User Journey:**
  1. Login
  2. View existing card of day
  3. Generate spreads (within limit)
  4. View history (check retention)
  5. Change deck selection
  6. Generate new spread (verify new deck used)

- [ ] **Admin Journey:**
  1. Login as keeper
  2. View dashboard metrics
  3. Search for user
  4. Change user role
  5. Update role policy (limits)
  6. Upload new deck (ZIP)
  7. Verify deck appears in user deck selection
  8. Update global settings

**Cross-Browser Testing:**
- [ ] Chrome (desktop)
- [ ] Safari (desktop)
- [ ] Safari (iOS)
- [ ] Chrome (Android)

**i18n Testing:**
- [ ] Switch language to RU, navigate all pages
- [ ] Verify all text translated
- [ ] Verify card names/interpretations in RU

**Theme Testing:**
- [ ] Switch to dark theme, navigate all pages
- [ ] Verify all colors use CSS variables
- [ ] Verify readability in both themes

**Deliverable:** All manual tests passed

---

### 8.5.7 Bug Fixing

**Process:**
- [ ] Document all bugs found during testing
- [ ] Prioritize: Critical → High → Medium → Low
- [ ] Fix critical and high-priority bugs
- [ ] Retest after fixes

**Deliverable:** All critical bugs fixed

---

**Phase 5 Exit Criteria:**
- ✅ UI Kit complete and documented
- ✅ Responsive on mobile/tablet/desktop
- ✅ Basic accessibility compliance
- ✅ Performance targets met
- ✅ All critical bugs fixed
- ✅ Manual testing checklist completed
- ✅ Error handling robust

**Estimated Time:** 5-7 days

---

## 9. Phase 6: Deployment

**Duration:** 2-3 days  
**Goal:** Deploy to production, smoke test, and go live

### 9.6.1 Pre-Deployment Checklist

**Tasks:**
- [ ] Verify all Phase 5 exit criteria met
- [ ] Review and finalize documentation:
  - README.md (setup instructions)
  - CONTRIBUTING.md (for future developers)
  - CHANGELOG.md (initial version)
- [ ] Set production environment variables:
  - Backend: `GIN_MODE=release`, prod MongoDB URI, strong JWT secret
  - Frontend: prod API URL
- [ ] Build production artifacts:
  - Backend: `go build -o tarot-daily-api`
  - Frontend: `npm run build`
- [ ] Backup MongoDB (if existing data)
- [ ] Run seed script on production database

**Deliverable:** Production build ready

---

### 9.6.2 Deployment

**Tasks:**
- [ ] Deploy backend binary to server
- [ ] Deploy frontend static files
- [ ] Configure web server (nginx/caddy) to:
  - Serve frontend from `/`
  - Proxy `/api` to backend
  - Serve `/uploads` as static files
- [ ] Setup SSL certificate (Let's Encrypt)
- [ ] Verify MongoDB connection from production server
- [ ] Start backend service
- [ ] Verify frontend accessible

**Deliverable:** Application deployed and accessible

---

### 9.6.3 Smoke Testing

**Tasks:**
- [ ] Test production endpoints:
  - GET /api/gateway/health → should return "stub"
  - POST /api/auth/register → create test user
  - POST /api/auth/login → login test user
  - POST /api/reading/generate/card-of-day → generate card
  - POST /api/reading/generate/spread-a → generate spread
- [ ] Test frontend:
  - Open in browser
  - Register new user
  - Generate readings
  - Verify images load
  - Test theme switch
  - Test language switch
- [ ] Test admin panel:
  - Login as keeper
  - View dashboard
  - Upload deck (small test deck)

**Deliverable:** Production smoke tests passed

---

### 9.6.4 Monitoring Setup

**Tasks:**
- [ ] Setup log collection (view logs from backend)
- [ ] Setup basic uptime monitoring (ping /api/gateway/health)
- [ ] Document monitoring process

**Deliverable:** Basic monitoring in place

---

### 9.6.5 Post-Deployment

**Tasks:**
- [ ] Announce MVP launch (internal/beta users)
- [ ] Monitor error logs for first 24-48 hours
- [ ] Address any critical issues immediately
- [ ] Collect user feedback

**Deliverable:** MVP live and stable

---

**Phase 6 Exit Criteria:**
- ✅ Application deployed to production
- ✅ SSL certificate active
- ✅ Smoke tests passed
- ✅ Monitoring in place
- ✅ No critical errors in first 24 hours
- ✅ MVP officially launched

**Estimated Time:** 2-3 days

---

## 10. Risk Assessment

### 10.1 Technical Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| **MongoDB connection issues** | Medium | High | Test SSH tunnel early, document setup, have fallback plan |
| **Timezone calculation bugs** | Medium | High | Extensive testing with various timezones, edge case handling |
| **JWT token security** | Low | High | Use strong secret, implement token expiration, validate on every request |
| **File upload vulnerabilities** | Medium | Medium | Strict validation (file type, size, structure), sandboxed processing |
| **Performance bottleneck (DB queries)** | Low | Medium | Indexes in place, monitor query performance, optimize if needed |
| **Browser compatibility issues** | Low | Medium | Test on major browsers early, use polyfills if needed |

### 10.2 Schedule Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| **Scope creep** | High | High | Strict adherence to MVP scope, defer all non-essential features |
| **Underestimated complexity** | Medium | Medium | Add buffer time to estimates, prioritize core features first |
| **Blocking dependencies** | Low | Medium | Identify dependencies early, parallelize where possible |
| **Testing takes longer than expected** | Medium | Low | Allocate full phase to testing, don't rush |

### 10.3 Product Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| **Poor user adoption** | Medium | Medium | Focus on UX quality, gather beta feedback early, iterate |
| **Confusing UI/UX** | Low | High | Test with non-technical users, iterate based on feedback |
| **Inaccurate card interpretations** | Low | Medium | Review all card texts before launch, ensure quality |
| **Performance issues at scale** | Low | Low | Optimize early, monitor performance, plan for scaling later |

---

## 11. Testing Strategy

### 11.1 Testing Levels

**Manual Testing (Primary for MVP):**
- Unit: Test individual functions (utils, validators)
- Integration: Test API endpoints with real database
- E2E: Test complete user flows manually

**Automated Testing (Deferred):**
- Backend unit tests: Deferred until explicit request
- Frontend unit tests: Deferred until explicit request
- E2E tests (Playwright): Deferred until explicit request

### 11.2 Testing Checklist

**Backend:**
- [ ] All auth endpoints tested (register, login, token validation)
- [ ] All reading endpoints tested (card of day, spreads, history)
- [ ] All admin endpoints tested (metrics, users, policies, decks)
- [ ] Role-based access control tested
- [ ] Timezone calculations tested (edge cases)
- [ ] File upload validation tested
- [ ] Error handling tested (invalid inputs, database errors)

**Frontend:**
- [ ] All pages rendered correctly
- [ ] All forms validated correctly
- [ ] All user flows tested (registration → reading → history)
- [ ] All admin flows tested (user management, deck upload)
- [ ] Theme switching tested
- [ ] Language switching tested
- [ ] Responsive design tested (mobile, tablet, desktop)
- [ ] Error states tested (network errors, API errors)

**Integration:**
- [ ] Frontend ↔ Backend integration tested for all features
- [ ] Token refresh flow tested (future: if implemented)
- [ ] File upload frontend-to-backend tested

---

## 12. Post-MVP Roadmap

### Future Enhancements (Not in MVP)

**v1.1: Monetization**
- Stripe integration
- Subscription management
- Payment flows
- Upgrade prompts

**v1.2: Gateway Integration**
- Real AI-powered interpretations
- Chat with AI about readings
- Personalized insights

**v1.3: Advanced Features**
- 5-card spreads
- Custom spread creation
- Reversed card interpretations
- Spread position meanings (Past/Present/Future)

**v1.4: User Engagement**
- Email notifications (card of day reminder)
- Push notifications (mobile PWA)
- Streak tracking and rewards
- Journaling with readings

**v1.5: Infrastructure**
- Object storage migration (S3/R2)
- CDN for images
- Advanced analytics (Mixpanel)
- Automated tests (Playwright)

---

**End of Implementation Roadmap**

**Success Criteria:** MVP delivered with all core features functional, tested, and ready for users.

**Next Steps:**
1. Review and approve this roadmap
2. Begin Phase 0: Project Setup
3. Track progress weekly
4. Adjust timeline as needed based on actual progress

**Document Owner:** Engineering Lead  
**Stakeholders:** Product, Design, Business  
**Review Cycle:** Weekly during development

---
