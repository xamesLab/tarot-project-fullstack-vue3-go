# Technical Specification
## Nebula Tarot - MVP v1.0

**Document Version:** 1.0  
**Last Updated:** April 3, 2026  
**Status:** Working Baseline

---

## Table of Contents

1. [System Architecture](#1-system-architecture)
2. [Technology Stack](#2-technology-stack)
3. [Database Schema](#3-database-schema)
4. [API Contracts](#4-api-contracts)
5. [Authentication & Authorization](#5-authentication--authorization)
6. [File Management](#6-file-management)
7. [Internationalization Implementation](#7-internationalization-implementation)
8. [Theme System](#8-theme-system)
9. [Error Handling](#9-error-handling)
10. [Security](#10-security)
11. [Development Environment](#11-development-environment)
12. [Deployment Strategy](#12-deployment-strategy)

---

## 1. System Architecture

### 1.1 High-Level Architecture

```
┌─────────────────────────────────────────────────┐
│                   CLIENT                        │
│  ┌──────────────────────────────────────────┐   │
│  │          Vue 3 SPA                       │   │
│  │  ├─ Pages (Dashboard, Profile, Admin)   │   │
│  │  ├─ Components (Card, Spread, Modal)    │   │
│  │  ├─ Stores (Pinia: auth, readings, app) │   │
│  │  ├─ Services (API client)               │   │
│  │  ├─ Router (vue-router)                 │   │
│  │  └─ i18n (vue-i18n)                     │   │
│  └──────────────────────────────────────────┘   │
└─────────────────┬───────────────────────────────┘
                  │ HTTPS
                  │ JSON
                  ↓
┌─────────────────────────────────────────────────┐
│                   SERVER                        │
│  ┌──────────────────────────────────────────┐   │
│  │          Go + Gin API                    │   │
│  │  ├─ Router (Gin)                         │   │
│  │  ├─ Middleware (Auth, CORS, Logging)    │   │
│  │  ├─ Handlers (Controllers)              │   │
│  │  ├─ Services (Business Logic)           │   │
│  │  └─ Repositories (Data Access)          │   │
│  └──────────────────────────────────────────┘   │
└─────────────────┬───────────────────────────────┘
                  │
                  ↓
┌─────────────────────────────────────────────────┐
│               PERSISTENCE                       │
│  ┌──────────────────────┐  ┌─────────────────┐  │
│  │   MongoDB            │  │ Local Filesystem│  │
│  │  (self-hosted)       │  │ /uploads/       │  │
│  │  - users             │  │  ├─ decks/      │  │
│  │  - cards             │  │  └─ backs/      │  │
│  │  - daily_readings    │  └─────────────────┘  │
│  │  - decks             │                       │
│  │  - card_backs        │                       │
│  │  - role_policies     │                       │
│  │  - admin_settings    │                       │
│  │  - ask_requests      │                       │
│  └──────────────────────┘                       │
└─────────────────────────────────────────────────┘
```

### 1.2 Request Flow

**Typical Reading Request:**

```
User clicks "Get Card of the Day"
    ↓
[Frontend] → POST /api/reading/generate/card-of-day
    ↓
[Router] → Auth Middleware (verify JWT)
    ↓
[Handler] → ReadingService.GenerateCardOfDay(userID, timezone)
    ↓
[Service] → Check if card already generated today
    ↓ (if not generated)
[Service] → Calculate seed: hash(userID + accountingDate + timezone)
    ↓
[Service] → Select card using deterministic random
    ↓
[Repository] → CardRepository.FindByKey(cardKey)
    ↓
[Repository] → ReadingRepository.Create(reading)
    ↓
[Service] → Return reading with card data
    ↓
[Handler] → JSON response
    ↓
[Frontend] → Display card
```

### 1.3 Directory Structure

#### Backend (`backend/`)
```
backend/
├── main.go                 # Entry point
├── go.mod                  # Dependencies
├── go.sum
├── config/
│   ├── config.go          # Configuration loader
│   └── env.go             # Environment variables
├── internal/
│   ├── middleware/
│   │   ├── auth.go        # JWT verification
│   │   ├── cors.go        # CORS configuration
│   │   ├── logger.go      # Request logging
│   │   └── role.go        # Role-based access control
│   ├── handlers/
│   │   ├── auth.go        # Auth endpoints
│   │   ├── user.go        # User profile endpoints
│   │   ├── reading.go     # Reading generation endpoints
│   │   ├── admin.go       # Admin endpoints
│   │   ├── catalog.go     # Deck/card back catalog endpoints
│   │   └── gateway.go     # Gateway stub endpoints
│   ├── services/
│   │   ├── auth.go        # Auth business logic
│   │   ├── user.go        # User management
│   │   ├── reading.go     # Reading generation logic
│   │   ├── catalog.go     # Catalog selection logic
│   │   └── admin.go       # Admin operations
│   ├── repositories/
│   │   ├── user.go        # User data access
│   │   ├── card.go        # Card data access
│   │   ├── reading.go     # Reading data access
│   │   ├── deck.go        # Deck data access
│   │   ├── card_back.go   # Card back data access
│   │   ├── policy.go      # Role policy data access
│   │   └── admin_settings.go # Settings data access
│   ├── models/
│   │   ├── user.go        # User model
│   │   ├── card.go        # Card model
│   │   ├── reading.go     # Reading model
│   │   ├── deck.go        # Deck model
│   │   ├── admin.go       # Admin API models
│   │   ├── policy.go      # RolePolicy model
│   │   └── settings.go    # AdminSettings model
│   └── utils/
│       ├── jwt.go         # JWT utilities
│       ├── hash.go        # Hashing utilities
│       ├── random.go      # Deterministic random
│       ├── time.go        # Time/timezone utilities
│       └── validator.go   # Input validation
└── uploads/               # File storage (gitignored)
    ├── decks/
    │   └── {deck_id}/
    │       └── cards/
    │           ├── the_fool.png
    │           └── ...
    ├── card-backs/
    │   └── {card_back_id}/
    │       └── back.png
    └── temp/              # Temporary uploads
```

#### Frontend (`frontend/`)
```
frontend/
├── src/
│   ├── main.js            # Entry point
│   ├── App.vue            # Root component
│   ├── router/
│   │   └── index.js       # Vue Router config
│   ├── stores/
│   │   ├── auth.js        # Auth state (Pinia)
│   │   ├── readings.js    # Reading state
│   │   ├── app.js         # App state (theme, language)
│   │   └── admin.js       # Admin state
│   ├── services/
│   │   ├── api.js         # Axios instance & interceptors
│   │   ├── auth.js        # Auth API calls
│   │   ├── user.js        # User API calls
│   │   ├── reading.js     # Reading API calls
│   │   └── admin.js       # Admin API calls
│   ├── pages/
│   │   ├── Login.vue
│   │   ├── Register.vue
│   │   ├── Dashboard.vue
│   │   ├── History.vue
│   │   ├── Profile.vue
│   │   ├── admin/
│   │   │   ├── AdminDashboard.vue
│   │   │   ├── Users.vue
│   │   │   ├── Decks.vue
│   │   │   ├── RolePolicies.vue
│   │   │   └── Settings.vue
│   │   └── dev/
│   │       └── UIKit.vue  # Dev-only
│   ├── components/
│   │   ├── base/
│   │   │   ├── BaseButton.vue
│   │   │   ├── BaseInput.vue
│   │   │   ├── BaseSelect.vue
│   │   │   ├── BaseSurface.vue
│   │   │   ├── BaseSegmentedToggle.vue
│   │   │   ├── BaseSwitchToggle.vue
│   │   │   ├── BaseTabs.vue
│   │   │   └── BaseModal.vue
│   │   ├── form/
│   │   │   ├── FormField.vue
│   │   │   └── AuthCard.vue
│   │   ├── layout/
│   │   │   ├── AppHeader.vue
│   │   │   ├── AppFooter.vue
│   │   │   ├── PageContainer.vue
│   │   │   └── PageSection.vue
│   │   ├── domain/
│   │   │   ├── UserChip.vue
│   │   │   ├── TarotCardPreview.vue
│   │   │   └── ReadingCard.vue
│   │   └── feedback/
│   │       ├── EmptyStateCard.vue
│   │       └── LoadingCard.vue
│   ├── composables/
│   │   ├── useAuth.js     # Auth composition
│   │   ├── useLocale.js   # Locale composition
│   │   └── useTheme.js    # Theme composition
│   ├── locales/
│   │   ├── en.json        # English translations
│   │   └── ru.json        # Russian translations
│   ├── styles/
│   │   ├── main.scss      # Global styles entry
│   │   ├── tokens.scss    # Design tokens
│   │   ├── themes/
│   │   │   ├── _light.scss
│   │   │   └── _dark.scss
│   │   ├── foundations/
│   │   │   ├── _reset.scss
│   │   │   ├── _typography.scss
│   │   │   └── _motion.scss
│   │   └── components/
│   │       ├── _button.scss
│   │       ├── _field.scss
│   │       ├── _surface.scss
│   │       └── _tabs.scss
│   └── utils/
│       ├── date.js        # Date formatting
│       └── validators.js  # Form validation
├── package.json
└── vite.config.js
```

**Design-system note:** canonical frontend visual and component architecture is defined in `specs/docs/DESIGN_SYSTEM.md`. The UI Kit page at `frontend/src/pages/dev/UIKit.vue` is the approved dev-only showcase/reference implementation for `v0.1.0`.

---

## 2. Technology Stack

### 2.1 Backend

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Language** | Go | 1.22+ | Backend runtime |
| **Framework** | Gin | 1.9+ | HTTP router |
| **Database** | MongoDB | 6.0+ | Primary data store |
| **Database Driver** | mongo-go-driver | 1.13+ | MongoDB client |
| **Authentication** | JWT | via golang-jwt/jwt | Token-based auth |
| **Validation** | go-playground/validator | 10.16+ | Input validation |
| **Environment** | godotenv | 1.5+ | Env variable loading |
| **Hashing** | bcrypt | - | Password hashing |
| **UUID** | google/uuid | 1.5+ | Unique ID generation |

**Key Dependencies:**
```go
// go.mod
module tarot-daily

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/joho/godotenv v1.5.1
    go.mongodb.org/mongo-driver v1.13.1
    github.com/google/uuid v1.5.0
    golang.org/x/crypto v0.17.0
    github.com/go-playground/validator/v10 v10.16.0
)
```

### 2.2 Frontend

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Language** | JavaScript | ES2022 | Frontend runtime |
| **Framework** | Vue 3 | 3.4+ | UI framework |
| **State** | Pinia | 2.1+ | State management |
| **Router** | Vue Router | 4.2+ | Client-side routing |
| **Build Tool** | Vite | 5.0+ | Dev server & bundler |
| **HTTP Client** | Axios | 1.6+ | API requests |
| **i18n** | Vue I18n | 9.8+ | Internationalization |
| **Styles** | SCSS | - | CSS preprocessor |
| **Date** | date-fns | 3.0+ | Date formatting |
| **Validation** | VeeValidate | 4.12+ | Form validation |

**Key Dependencies:**
```json
// package.json
{
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.2.5",
    "pinia": "^2.1.7",
    "axios": "^1.6.2",
    "vue-i18n": "^9.8.0",
    "date-fns": "^3.0.6",
    "vee-validate": "^4.12.4"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0",
    "vite": "^5.0.0",
    "sass": "^1.69.5"
  }
}
```

### 2.3 Infrastructure

| Component | Technology | Notes |
|-----------|-----------|-------|
| **Database** | MongoDB | Self-hosted, SSH tunnel for local dev |
| **File Storage** | Local Filesystem | `/uploads/` directory, no object storage in MVP |
| **Web Server** | Gin (Go) | Serves API + static files |
| **Development** | Local | Go backend + Vite dev server |

---

## 3. Database Schema

### 3.1 MongoDB Collections

#### 3.1.1 `users`

**Purpose:** Store user accounts and settings

```javascript
{
  "_id": ObjectId("..."),
  "username": "alice123",           // String, unique, 3-20 chars
  "email": "alice@example.com",     // String, unique, valid email
  "password_hash": "bcrypt...",     // String, bcrypt hash
  "role": "seeker",                 // String: seeker | adept | master | keeper
  "status": "active",               // String: active | disabled
  "timezone": "Europe/Amsterdam",   // String, IANA timezone
  "language": "en",                 // String: en | ru
  "theme": "light",                 // String: light | dark
  "selected_deck_id": "",           // String, optional, references decks._id
  "selected_back_id": "",           // String, optional, references card_backs._id
  "created_at": ISODate("..."),     // Date
  "updated_at": ISODate("..."),     // Date
  "last_login_at": ISODate("...")   // Date, nullable
}
```

**Indexes:**
```javascript
db.users.createIndex({ "email": 1 }, { unique: true })
db.users.createIndex({ "username": 1 }, { unique: true })
db.users.createIndex({ "role": 1 })
db.users.createIndex({ "status": 1 })
```

#### 3.1.2 `cards`

**Purpose:** Store tarot card catalog with multilingual data

```javascript
{
  "_id": "the_fool",                // String, card key (unique)
  "id": 100,                        // Number, original ID from JSON
  "key": "the_fool",                // String, unique identifier
  "arcana": "major",                // String: major | minor
  "type": {                         // Card type/suit
    "en": "0",                      // For Major: roman numeral
    "ru": "0"                       // For Minor: suit name
  },
  "names": {
    "en": "The Fool",
    "ru": "Шут"
  },
  "interpretations": {
    "en": "New beginnings, innocence...",
    "ru": "Новые начинания, невинность..."
  },
  "tags": {
    "en": "New ideas, inspiration, opportunities, enthusiasm, creative energy.",
    "ru": "Новые идеи, вдохновение, возможности, энтузиазм, творческая энергия."
  },
  "created_at": ISODate("..."),
  "updated_at": ISODate("...")
}
```

**Indexes:**
```javascript
db.cards.createIndex({ "key": 1 }, { unique: true })
db.cards.createIndex({ "arcana": 1 })
```

**Total Documents:** 78 (22 Major + 56 Minor)

#### 3.1.3 `daily_readings`

**Purpose:** Store user reading history

```javascript
{
  "_id": ObjectId("..."),
  "user_id": ObjectId("..."),       // References users._id
  "type": "card_of_day",            // String: card_of_day | three_card
  "accounting_date": "2026-03-31",  // String, YYYY-MM-DD format
  "timezone": "Europe/Amsterdam",   // String, user's timezone at generation
  "cards": [                        // Array of card data
    {
      "key": "the_fool",            // String, references cards._id
      "position": 0,                // Number, position in spread
      "reversed": false,            // Boolean (always false in MVP)
      "viewed": false               // Boolean, persisted reveal state
    }
    // For three_card: 3 cards
    // For card_of_day: 1 card
  ],
  "deck_id": "default",             // String, deck used for this reading
  "back_id": "default",             // String, card back used for this reading
  "generated_at": ISODate("..."),   // Date, when reading was created
  "expires_at": ISODate("..."),     // Date, calculated based on role retention
  "ai": {                           // Optional cached AI interpretation for reuse
    "request_id": ObjectId("..."),
    "provider": "openrouter",
    "model": "provider/model-id",
    "prompt_version": "v0.1.1-compact-100w",
    "language": "ru",
    "status": "completed",
    "content": "Короткая интерпретация...",
    "error_code": null,
    "error_message": null,
    "last_requested_at": ISODate("..."),
    "completed_at": ISODate("..."),
    "updated_at": ISODate("...")
  },
  "metadata": {                     // Optional metadata
    "seed": "hash...",              // For card_of_day: deterministic seed
    "request_ip": "192.168.1.1"     // Optional: user IP
  }
}
```

**Indexes:**
```javascript
db.daily_readings.createIndex({ "user_id": 1, "accounting_date": 1, "type": 1 })
db.daily_readings.createIndex({ "user_id": 1, "generated_at": -1 })
db.daily_readings.createIndex({ "expires_at": 1 }, { expireAfterSeconds: 0 }) // TTL index
```

**TTL Index:**
- MongoDB automatically deletes documents where `expires_at` < current time
- Retention calculated on creation: `expires_at = generated_at + role_retention_days`

#### 3.1.4 `decks`

**Purpose:** Store deck metadata and file references

```javascript
{
  "_id": "default",                 // String, unique deck ID
  "name": {
    "en": "Tarot Classic",
    "ru": "Таро Классика"
  },
  "description": {
    "en": "Traditional tarot imagery...",
    "ru": "Традиционные образы таро..."
  },
  "author": "Pamela Colman Smith",  // String, optional
  "year": 1909,                     // Number, optional
  "tags": ["traditional", "classic"],
  "status": "done",                 // String: processing | done | archived
  "card_count": 78,                 // Number, must be 78
  "card_image_base_path": "/uploads/decks/default/cards/",
  "created_at": ISODate("..."),
  "updated_at": ISODate("..."),
  "created_by": ObjectId("...")     // References users._id (admin who created)
}
```

**Indexes:**
```javascript
db.decks.createIndex({ "status": 1 })
db.decks.createIndex({ "tags": 1 })
```

#### 3.1.5 `card_backs`

**Purpose:** Store card back metadata independently from decks

```javascript
{
  "_id": "default",                 // String, unique card back ID
  "name": {
    "en": "Classic Back",
    "ru": "Классическая рубашка"
  },
  "description": {
    "en": "Default classic card back.",
    "ru": "Базовая классическая рубашка карт."
  },
  "tags": ["traditional", "classic"],
  "status": "done",                 // String: processing | done | archived
  "image_path": "/uploads/card-backs/default/back.png",
  "created_at": ISODate("..."),
  "updated_at": ISODate("...")
}
```

**Indexes:**
```javascript
db.card_backs.createIndex({ "status": 1 })
db.card_backs.createIndex({ "tags": 1 })
```

#### 3.1.6 `role_policies`

**Purpose:** Store daily limits per role

```javascript
{
  "_id": "seeker",                  // String, role name
  "role": "seeker",                 // String: seeker | adept | master | keeper
  "display_name": {
    "en": "Seeker",
    "ru": "Искатель"
  },
  "three_card_limit": 3,            // Number, daily limit for 3-card spreads
  "history_retention_days": 1,      // Number, days to keep history
  "created_at": ISODate("..."),
  "updated_at": ISODate("...")
}
```

**Indexes:**
```javascript
db.role_policies.createIndex({ "role": 1 }, { unique: true })
```

**Default Policies:**
```javascript
[
  { "_id": "seeker", "role": "seeker", "three_card_limit": 3, "history_retention_days": 1 },
  { "_id": "adept", "role": "adept", "three_card_limit": 5, "history_retention_days": 7 },
  { "_id": "master", "role": "master", "three_card_limit": 10, "history_retention_days": 15 },
  { "_id": "keeper", "role": "keeper", "three_card_limit": 999999, "history_retention_days": 30 }
]
```

#### 3.1.7 `admin_settings`

**Purpose:** Global application settings

```javascript
{
  "_id": "global",                  // String, singleton document
  "maintenance_mode": false,        // Boolean, if true, show maintenance page
  "maintenance_message": {
    "en": "Site under maintenance",
    "ru": "Сайт на обслуживании"
  },
  "default_role": "seeker",         // String, role for new users
  "default_deck_id": "default",     // String, global fallback deck
  "default_back_id": "default",     // String, global fallback card back
  "registration_enabled": true,     // Boolean, allow new registrations
  "created_at": ISODate("..."),
  "updated_at": ISODate("...")
}
```

**Note:** Only one document with `_id: "global"`

#### 3.1.8 `ask_requests`

**Purpose:** Persist direct AI interpretation requests and keep compatibility with a future gateway/fallback architecture

```javascript
{
  "_id": ObjectId("..."),
  "user_id": ObjectId("..."),       // References users._id
  "reading_id": ObjectId("..."),    // References daily_readings._id
  "provider": "openrouter",         // String, current primary provider for v0.1.1
  "model": "provider/model-id",     // String, resolved model used for generation
  "prompt_version": "v1",           // String, prompt template version
  "status": "pending",              // String: pending | processing | completed | failed
  "request_payload": {...},         // Object, normalized request snapshot
  "response_data": {...},           // Object, normalized provider response
  "response_excerpt": "Short text", // String, compact AI interpretation shown in UI
  "usage": {                        // Optional provider usage snapshot
    "prompt_tokens": 123,
    "completion_tokens": 45,
    "total_tokens": 168
  },
  "error_code": "provider_timeout", // Nullable
  "error_message": "Timed out",     // Nullable, safe/public form only
  "requested_at": ISODate("..."),
  "created_at": ISODate("..."),
  "updated_at": ISODate("..."),
  "completed_at": ISODate("...")    // Nullable
}
```

**Current note:** `v0.1.0` kept this collection as a stub. `v0.1.1` promotes it to the primary persistence layer for direct AI requests through `OpenRouter`, while preserving room for a future `gateway` fallback model.

**Storage note:** the finalized compact interpretation is additionally cached in `daily_readings.ai`, so reopening the same reading does not require a new provider request.

---

### 3.2 Data Relationships

```
users (1) ──────── (N) daily_readings
  │
  └──── (1) selected_deck ──> decks (1)
  │
  └──── (1) selected_back ──> card_backs (1)
  │
  └──── (1) role ──> role_policies (1)

daily_readings (N) ──────── (N) cards (via cards array)

decks (1) ──────── (78) card_files (filesystem)
card_backs (1) ──────── (1) back_image (filesystem)
```

---

## 4. API Contracts

### 4.1 Authentication Endpoints

#### POST `/api/auth/register`

**Description:** Register new user

**Request:**
```json
{
  "email": "alice@example.com",
  "username": "alice123",
  "password": "SecurePass123",
  "timezone": "Europe/Amsterdam"
}
```

**Validation:**
- `email`: required, valid email, unique
- `username`: optional, 3-20 chars, letters/numbers/dot/underscore/hyphen, unique if provided
- `password`: required, min 8 chars, contains letter + number
- `timezone`: required, valid IANA timezone

**Success Response (201):**
```json
{
  "success": true,
  "user": {
    "id": "uuid",
    "username": "alice123",
    "email": "alice@example.com",
    "role": "seeker",
    "timezone": "Europe/Amsterdam",
    "language": "en",
    "theme": "light"
  },
  "token": "jwt_token_here"
}
```

**Error Responses:**
```json
// 400 Bad Request
{
  "success": false,
  "error": "validation_failed",
  "message": "Email is already registered"
}

// 400 Bad Request
{
  "success": false,
  "error": "validation_failed",
  "message": "Password must be at least 8 characters with letter and number"
}
```

---

#### POST `/api/auth/login`

**Description:** Authenticate user

**Request:**
```json
{
  "email": "alice@example.com",
  "password": "SecurePass123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "user": {
    "id": "uuid",
    "username": "alice123",
    "email": "alice@example.com",
    "role": "adept",
    "timezone": "Europe/Amsterdam",
    "language": "ru",
    "theme": "dark"
  },
  "token": "jwt_token_here"
}
```

**Error Response (401):**
```json
{
  "success": false,
  "error": "invalid_credentials",
  "message": "Email or password incorrect"
}
```

---

### 4.2 User Endpoints

#### GET `/api/me`

**Description:** Get current user info

**Headers:** `Authorization: Bearer {token}`

**Response (200):**
```json
{
  "success": true,
  "user": {
    "id": "uuid",
    "username": "alice123",
    "email": "alice@example.com",
    "role": "adept",
    "status": "active",
    "timezone": "Europe/Amsterdam",
    "language": "en",
    "theme": "light",
    "created_at": "2026-03-15T10:00:00Z",
    "last_login_at": "2026-03-31T14:23:00Z"
  }
}
```

---

#### PATCH `/api/me/profile`

**Description:** Update user profile

**Headers:** `Authorization: Bearer {token}`

**Request:**
```json
{
  "username": "alice_tarot",
  "timezone": "America/New_York"
}
```

**Response (200):**
```json
{
  "success": true,
  "user": {
    "id": "uuid",
    "username": "alice_tarot",
    "timezone": "America/New_York",
    ...
  }
}
```

---

#### PATCH `/api/me/settings`

**Description:** Update user settings

**Notes:**
- `selected_deck_id` and `selected_back_id` are optional user overrides
- Empty string resets the value back to the admin default
- Custom deck/back selection is available only for roles `adept`, `master`, and `keeper`

**Request:**
```json
{
  "language": "ru",
  "theme": "dark",
  "selected_deck_id": "rider-waite-modern",
  "selected_back_id": "celestial"
}
```

**Response (200):**
```json
{
  "success": true,
  "user": {
    "id": "uuid",
    "username": "alice123",
    "email": "alice@example.com",
    "role": "adept",
    "status": "active",
    "language": "ru",
    "theme": "dark",
    "selected_deck_id": "rider-waite-modern",
    "selected_back_id": "celestial"
  }
}
```

---

### 4.3 Reading Endpoints

#### POST `/api/reading/generate/card-of-day`

**Description:** Generate Card of the Day

**Headers:** `Authorization: Bearer {token}`

**Request:**
```json
{
  "timezone": "Europe/Amsterdam"
}
```

**Success Response (200) - New Card:**
```json
{
  "success": true,
  "reading": {
    "id": "uuid",
    "type": "card_of_day",
    "accounting_date": "2026-03-31",
    "generated_at": "2026-03-31T15:30:00Z",
    "cards": [
      {
        "key": "the_fool",
        "name": "The Fool",
        "name_ru": "Шут",
        "interpretation": "New beginnings, innocence, spontaneity...",
        "interpretation_ru": "Новые начинания, невинность, спонтанность...",
        "tags": "New ideas, inspiration, opportunities...",
        "tags_ru": "Новые идеи, вдохновение, возможности...",
        "image_url": "/uploads/decks/default/cards/the_fool.png",
        "back_image_url": "/uploads/card-backs/default/back.png",
        "position": 0,
        "reversed": false,
        "viewed": false
      }
    ],
    "deck_id": "default",
    "back_id": "default",
    "next_available_at": "2026-04-01T03:00:00+02:00"
  }
}
```

**Error Response (400) - Already Generated:**
```json
{
  "success": false,
  "error": "card_of_day_already_generated",
  "message": "You already received your Card of the Day",
  "existing_reading": {
    "id": "uuid",
    "type": "card_of_day",
    "accounting_date": "2026-03-31",
    "generated_at": "2026-03-31T10:15:00Z",
    "cards": [...]
  },
  "next_available_at": "2026-04-01T03:00:00+02:00"
}
```

---

#### POST `/api/reading/generate/spread-a`

**Description:** Generate 3-card spread (renamed from "spread-a" to "three-card" for clarity)

**Headers:** `Authorization: Bearer {token}`

**Request:**
```json
{
  "timezone": "Europe/Amsterdam"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "reading": {
    "id": "uuid",
    "type": "three_card",
    "accounting_date": "2026-03-31",
    "generated_at": "2026-03-31T16:45:00Z",
    "cards": [
      {
        "key": "the_magician",
        "name": "The Magician",
        "name_ru": "Маг",
        "interpretation": "...",
        "interpretation_ru": "...",
        "tags": "...",
        "tags_ru": "...",
        "image_url": "/uploads/decks/default/cards/the_magician.png",
        "back_image_url": "/uploads/card-backs/default/back.png",
        "position": 0,
        "reversed": false,
        "viewed": false
      },
      {
        "key": "ace_of_cups",
        "name": "Ace of Cups",
        "name_ru": "Туз Кубков",
        "interpretation": "...",
        "interpretation_ru": "...",
        "tags": "...",
        "tags_ru": "...",
        "image_url": "/uploads/decks/default/cards/ace_of_cups.png",
        "back_image_url": "/uploads/card-backs/default/back.png",
        "position": 1,
        "reversed": false,
        "viewed": false
      },
      {
        "key": "strength",
        "name": "Strength",
        "name_ru": "Сила",
        "interpretation": "...",
        "interpretation_ru": "...",
        "tags": "...",
        "tags_ru": "...",
        "image_url": "/uploads/decks/default/cards/strength.png",
        "back_image_url": "/uploads/card-backs/default/back.png",
        "position": 2,
        "reversed": false,
        "viewed": false
      }
    ],
    "deck_id": "default",
    "back_id": "default",
    "remaining_today": {
      "count": 2,
      "limit": 5
    }
  }
}
```

**Error Response (429) - Limit Reached:**
```json
{
  "success": false,
  "error": "daily_limit_reached",
  "message": "You've reached your daily limit of 5 spreads",
  "next_reset_at": "2026-04-01T03:00:00+02:00",
  "current_limit": {
    "count": 5,
    "limit": 5
  }
}
```

---

#### GET `/api/reading/history`

**Description:** Get user's reading history

**Headers:** `Authorization: Bearer {token}`

**Query Parameters:**
- `type` (optional): `card_of_day` | `three_card` | `all` (default: `all`)
- `page` (optional): number (default: 1)
- `limit` (optional): number (default: 10, max: 50)

**Response (200):**
```json
{
  "success": true,
  "readings": [
    {
      "id": "uuid",
      "type": "card_of_day",
      "accounting_date": "2026-03-31",
      "generated_at": "2026-03-31T10:15:00Z",
      "cards": [...]
    },
    {
      "id": "uuid",
      "type": "three_card",
      "accounting_date": "2026-03-31",
      "generated_at": "2026-03-31T14:30:00Z",
      "cards": [...]
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 23,
    "total_pages": 3
  },
  "retention_info": {
    "role": "adept",
    "retention_days": 7,
    "oldest_date": "2026-03-24"
  }
}
```

---

#### POST `/api/reading/:readingId/cards/:position/reveal`

**Description:** Persistently reveal one card inside an existing reading

**Headers:** `Authorization: Bearer {token}`

**Success Response (200):**
```json
{
  "success": true,
  "reading": {
    "id": "uuid",
    "type": "three_card",
    "cards": [
      {
        "key": "the_magician",
        "position": 0,
        "viewed": true
      },
      {
        "key": "ace_of_cups",
        "position": 1,
        "viewed": false
      }
    ]
  }
}
```

**Behavior Notes:**
- reveal is idempotent for an already opened card
- viewed state is stored in `daily_readings.cards[].viewed`
- reveal affects only the targeted card and only for the owner of the reading

---

### 4.4 Deck Endpoints

#### GET `/api/decks/available`

**Description:** Get available decks for user

**Headers:** `Authorization: Bearer {token}`

**Response (200):**
```json
{
  "success": true,
  "default_deck_id": "default",
  "decks": [
    {
      "id": "default",
      "name": "Tarot Classic",
      "name_ru": "Таро Классика",
      "description": "Traditional tarot imagery...",
      "description_ru": "Традиционные образы таро...",
      "author": "Pamela Colman Smith",
      "tags": ["traditional", "classic"],
      "preview_images": [
        "/uploads/decks/default/cards/the_fool.png",
        "/uploads/decks/default/cards/the_magician.png",
        "/uploads/decks/default/cards/the_empress.png"
      ]
    },
    {
      "id": "rider-waite-modern",
      "name": "Rider-Waite Modern",
      "name_ru": "Райдер-Уэйт Модерн",
      ...
    }
  ]
}
```

**Note:** Only returns decks with `status: "done"`

**Role Note:** custom deck selection is available only to `adept`, `master`, and `keeper`. Other roles continue using the global default deck.

---

#### GET `/api/card-backs/available`

**Description:** Get available card backs for user

**Headers:** `Authorization: Bearer {token}`

**Response (200):**
```json
{
  "success": true,
  "default_back_id": "default",
  "card_backs": [
    {
      "id": "default",
      "name": "Classic Back",
      "name_ru": "Классическая рубашка",
      "image_url": "/uploads/card-backs/default/back.png"
    },
    {
      "id": "celestial",
      "name": "Celestial",
      "name_ru": "Небесная",
      "image_url": "/uploads/card-backs/celestial/back.png"
    }
  ]
}
```

**Note:** card back selection is independent from deck selection and is available only to `adept`, `master`, and `keeper`. Other roles continue using the global default card back.

---

### 4.5 Admin Endpoints

**All admin endpoints require `role: "keeper"`**

#### GET `/api/admin/access-check`

**Description:** Verify admin access

**Headers:** `Authorization: Bearer {token}`

**Response (200):**
```json
{
  "success": true,
  "has_access": true
}
```

**Response (403):**
```json
{
  "success": false,
  "error": "forbidden",
  "message": "Admin access required"
}
```

---

#### GET `/api/admin/dashboard`

**Description:** Get admin dashboard metrics

**Response (200):**
```json
{
  "success": true,
  "metrics": {
    "users": {
      "total": 1247,
      "by_role": {
        "seeker": 1100,
        "adept": 120,
        "master": 25,
        "keeper": 2
      },
      "by_status": {
        "active": 1240,
        "disabled": 7
      }
    },
    "readings": {
      "today": {
        "total": 3842,
        "card_of_day": 1247,
        "three_card": 2595
      },
      "this_week": {
        "total": 18234,
        "card_of_day": 7421,
        "three_card": 10813
      }
    },
    "decks": {
      "total": 3,
      "by_status": {
        "done": 2,
        "processing": 1,
        "archived": 1
      }
    }
  }
}
```

---

#### GET `/api/admin/users`

**Description:** List users with search/filter

**Query Parameters:**
- `search` (optional): search username or email
- `role` (optional): filter by role
- `status` (optional): filter by status
- `page` (optional): page number (default: 1)
- `limit` (optional): items per page (default: 20, max: 100)

**Response (200):**
```json
{
  "success": true,
  "users": [
    {
      "id": "uuid",
      "username": "alice123",
      "email": "alice@example.com",
      "role": "adept",
      "status": "active",
      "created_at": "2026-03-15T10:00:00Z",
      "last_login_at": "2026-03-31T14:23:00Z"
    },
    ...
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1247,
    "total_pages": 63
  }
}
```

---

#### PATCH `/api/admin/users/:userId/role`

**Description:** Change user role

**Request:**
```json
{
  "role": "master"
}
```

**Response (200):**
```json
{
  "success": true,
  "user": {
    "id": "uuid",
    "username": "alice123",
    "role": "master"
  }
}
```

---

#### GET `/api/admin/roles`

**Description:** Get all role policies

**Response (200):**
```json
{
  "success": true,
  "policies": [
    {
      "role": "seeker",
      "display_name": { "en": "Seeker", "ru": "Искатель" },
      "three_card_limit": 3,
      "history_retention_days": 1
    },
    {
      "role": "adept",
      "display_name": { "en": "Adept", "ru": "Адепт" },
      "three_card_limit": 5,
      "history_retention_days": 7
    },
    {
      "role": "master",
      "display_name": { "en": "Master", "ru": "Мастер" },
      "three_card_limit": 10,
      "history_retention_days": 15
    },
    {
      "role": "keeper",
      "display_name": { "en": "Keeper", "ru": "Смотритель" },
      "three_card_limit": 999999,
      "history_retention_days": 30
    }
  ]
}
```

---

#### PATCH `/api/admin/roles/:role`

**Description:** Update role policy

**Request:**
```json
{
  "three_card_limit": 7,
  "history_retention_days": 14
}
```

**Response (200):**
```json
{
  "success": true,
  "policy": {
    "role": "adept",
    "three_card_limit": 7,
    "history_retention_days": 14,
    "updated_at": "2026-03-31T16:00:00Z"
  }
}
```

---

#### POST `/api/admin/decks`

**Description:** Create new deck from ZIP archive only

**Request (multipart/form-data):**
```
name_en: "Ethereal Dreams"
name_ru: "Эфирные Мечты"
description_en: "A mystical deck..."
description_ru: "Мистическая колода..."
author: "Jane Smith"
tags: "mystical,modern,ethereal"
cards_zip: [file upload] // ZIP with 78 card images named by tarot key
```

**Response (201):**
```json
{
  "success": true,
  "deck": {
    "id": "uuid",
    "name": { "en": "Ethereal Dreams", "ru": "Эфирные Мечты" },
    "status": "done",
    "card_count": 78
  },
  "message": "Deck created and processed successfully"
}
```

**Processing:**
- Backend extracts ZIP
- Accepts optional `0_meta.json`
- Validates 78 cards
- Validates filenames against known tarot `key` values
- Saves to `/uploads/decks/{deck_id}/`
- Updates `status: "done"`

---

#### POST `/api/admin/card-backs`

**Description:** Create new card back

**Request (multipart/form-data):**
```
name_en: "Moonlit Back"
name_ru: "Лунная рубашка"
description_en: "A celestial design..."
description_ru: "Небесный дизайн..."
tags: "moon,gold,minimal"
back_image: [file upload]
```

**Response (201):**
```json
{
  "success": true,
  "card_back": {
    "id": "uuid",
    "name": { "en": "Moonlit Back", "ru": "Лунная рубашка" },
    "status": "done",
    "image_path": "/uploads/card-backs/uuid/back.png"
  },
  "message": "Card back created successfully"
}
```

---

#### GET `/api/admin/card-backs`

**Description:** List uploaded card backs

**Query Parameters:**
- `scope` (optional): `active` | `archived` | `all` (default: `active`)

**Response (200):**
```json
{
  "success": true,
  "card_backs": [
    {
      "id": "default",
      "name": { "en": "Classic Back", "ru": "Классическая рубашка" },
      "status": "done",
      "image_path": "/uploads/card-backs/default/back.png"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

---

#### GET `/api/admin/decks/:deckId`

**Description:** Get deck details

**Response (200):**
```json
{
  "success": true,
  "deck": {
    "id": "uuid",
    "name": { "en": "Ethereal Dreams", "ru": "Эфирные Мечты" },
    "description": { ... },
    "author": "Jane Smith",
    "tags": ["mystical", "modern"],
    "status": "done",
    "card_count": 78,
    "cards": [
      {
        "key": "the_fool",
        "image_path": "/uploads/decks/uuid/cards/the_fool.png"
      },
      ...
    ]
  }
}
```

---

#### GET `/api/admin/decks`

**Description:** List uploaded decks

**Query Parameters:**
- `scope` (optional): `active` | `archived` | `all` (default: `active`)

---

#### POST `/api/admin/decks/:deckId/archive`

**Description:** Move deck to archive

**Behavior Notes:**
- archived decks are removed from active selection lists
- current default deck cannot be archived

---

#### POST `/api/admin/decks/:deckId/restore`

**Description:** Restore archived deck to active catalog

---

#### DELETE `/api/admin/decks/:deckId`

**Description:** Permanently delete archived deck and its filesystem assets

---

#### POST `/api/admin/card-backs/:cardBackId/archive`

**Description:** Move card back to archive

---

#### POST `/api/admin/card-backs/:cardBackId/restore`

**Description:** Restore archived card back to active catalog

---

#### DELETE `/api/admin/card-backs/:cardBackId`

**Description:** Permanently delete archived card back and its filesystem assets

---

#### GET/PATCH `/api/admin/settings`

**GET Response (200):**
```json
{
  "success": true,
  "settings": {
    "maintenance_mode": false,
    "maintenance_message": {
      "en": "Site under maintenance",
      "ru": "Сайт на обслуживании"
    },
    "default_role": "seeker",
    "registration_enabled": true,
    "default_deck_id": "default",
    "default_back_id": "default"
  }
}
```

**Behavior Notes:**
- `default_deck_id` and `default_back_id` are configured independently by admin
- They apply to all users who have not chosen custom values
- User-level overrides are allowed only for roles `adept`, `master`, and `keeper`
- Archived decks and card backs are excluded from available selection lists

**PATCH Request:**
```json
{
  "maintenance_mode": true,
  "registration_enabled": false
}
```

**PATCH Response (200):**
```json
{
  "success": true,
  "settings": {
    "maintenance_mode": true,
    "registration_enabled": false,
    ...
  }
}
```

---

### 4.6 AI Interpretation Endpoints

#### POST `/api/ask`

**Description:** Create or reuse an AI interpretation request for a reading. In `v0.1.1` this is expected to use direct `OpenRouter` integration as the primary provider path.

**Behavior note for `v0.1.1`:**
- allowed only for the owner of the reading
- request payload must be built from the reading type and revealed cards
- prompt must require a compact interpretation only
- intended response length is capped at `100` words
- frontend renders the resulting text inside the existing `Meaning` / `Значение расклада` block
- backend should reuse the cached `daily_readings.ai` value when the same reading was already interpreted successfully

**Legacy `v0.1.0` stub response (501):**
```json
{
  "success": false,
  "error": "not_implemented",
  "message": "Gateway integration is not available in this version"
}
```

---

#### GET `/api/ask/:request_id`

**Description:** Fetch AI interpretation request status and completed result for the current user.

**Planned role in `v0.1.1`:**
- returns persisted status lifecycle
- returns compact finalized text when ready
- never exposes provider secrets

---

#### GET `/api/gateway/health`

**Description:** Legacy gateway health check kept for compatibility while direct AI becomes primary in `v0.1.1`.

**Response (200):**
```json
{
  "status": "stub",
  "message": "Gateway is not implemented"
}
```

---

## 5. Authentication & Authorization

### 5.1 JWT Token Structure

**Payload:**
```json
{
  "user_id": "uuid",
  "role": "adept",
  "timezone": "Europe/Amsterdam",
  "iat": 1711896000,        // Issued at
  "exp": 1714488000         // Expires (30 days)
}
```

**Secret:** Loaded from environment variable `JWT_SECRET`

**Expiry:** 30 days

**Storage:** Frontend stores in `localStorage` as `auth_token`

### 5.2 Authentication Flow

```
1. User submits login credentials
   ↓
2. Backend validates email/password
   ↓
3. Backend generates JWT with user claims
   ↓
4. Frontend stores token in localStorage
   ↓
5. Frontend includes token in all API requests:
   Header: "Authorization: Bearer {token}"
   ↓
6. Backend middleware verifies token on each request
   ↓
7. If valid: extract user_id, role → proceed
   If invalid: return 401 Unauthorized
```

### 5.3 Middleware

#### Auth Middleware

**Purpose:** Verify JWT and attach user to request context

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := extractToken(c)
        if tokenString == "" {
            c.JSON(401, gin.H{"error": "missing_token"})
            c.Abort()
            return
        }
        
        claims, err := verifyToken(tokenString)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid_token"})
            c.Abort()
            return
        }
        
        // Attach user info to context
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Set("user_timezone", claims.Timezone)
        
        c.Next()
    }
}
```

#### Role Middleware

**Purpose:** Require specific role(s)

```go
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("user_role")
        
        for _, role := range allowedRoles {
            if userRole == role {
                c.Next()
                return
            }
        }
        
        c.JSON(403, gin.H{"error": "forbidden"})
        c.Abort()
    }
}
```

**Usage:**
```go
// Admin-only route
admin.GET("/dashboard", RequireRole("keeper"), handlers.AdminDashboard)

// Multiple roles allowed
api.GET("/profile", RequireRole("seeker", "adept", "master", "keeper"), handlers.GetProfile)
```

---

## 6. File Management

### 6.1 Upload Directory Structure

```
/uploads/
├── decks/
│   ├── default/
│   │   └── cards/
│   │       ├── the_fool.png
│   │       ├── the_magician.png
│   │       └── ... (78 cards)
│   ├── {deck_uuid_1}/
│   │   └── cards/
│   │       └── ...
│   └── {deck_uuid_2}/
│       └── ...
├── card-backs/
│   ├── default/
│   │   └── back.png
│   └── {card_back_uuid_1}/
│       └── back.png
└── temp/
    └── {upload_session_id}/
        └── ... (temporary uploads, cleaned up after processing)
```

### 6.2 File Upload Flow

**Deck Upload (ZIP):**

```
1. Admin uploads ZIP via POST /api/admin/decks
   ↓
2. Backend saves ZIP to /uploads/temp/{session_id}/
   ↓
3. Backend validates ZIP structure:
   - Contains exactly 78 card images
   - File names match known tarot `key` values
   - Optional `0_meta.json` may exist for frontend autofill
   - All files are PNG/JPG and under allowed size
   ↓
4. Backend extracts ZIP to /uploads/decks/{deck_id}/
   ↓
5. Backend updates deck status: "processing" → "done"
   ↓
6. Backend deletes temp files
```

**Card Back Upload:**

```
1. Admin uploads image via POST /api/admin/card-backs
   ↓
2. Backend saves image to /uploads/temp/{session_id}/
   ↓
3. Backend validates format, size, and aspect ratio
   ↓
4. Backend saves image to /uploads/card-backs/{card_back_id}/back.png
   ↓
5. Backend updates card back status: "processing" → "done"
```

**Validation:**
```go
func ValidateDeckZIP(zipPath string) error {
    // Extract to temp
    // Ignore 0_meta.json if present
    // Check for exactly 78 cards
    // Check that file names map to known tarot keys
    // Validate each file:
    //   - Format: PNG or JPG
    //   - Size: <= 2MB
    // Return error if any validation fails
}
```

### 6.3 Static File Serving

**Backend serves uploaded files:**

```go
// In main.go
router.Static("/uploads", "./uploads")
```

**Frontend accesses files:**
```
GET /uploads/decks/default/cards/the_fool.png
→ Returns image file
```

**CORS:** Allow frontend domain to access `/uploads/*`

---

## 7. Internationalization Implementation

### 7.1 Frontend (Vue I18n)

**Setup:**
```javascript
// src/i18n/index.js
import { createI18n } from 'vue-i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'

const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('language') || navigator.language.split('-')[0] || 'en',
  fallbackLocale: 'en',
  messages: {
    en,
    ru
  }
})

export default i18n
```

**Usage in Components:**
```vue
<template>
  <h1>{{ $t('dashboard.welcome') }}</h1>
  <button>{{ $t('readings.generate') }}</button>
  <p>{{ $t('readings.remainingToday', { count: 2, limit: 5 }) }}</p>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

function changeLanguage(lang) {
  locale.value = lang
  localStorage.setItem('language', lang)
  // Also update user profile via API
}
</script>
```

### 7.2 Backend (Localized Responses)

**Accept-Language Header:**

Frontend sends:
```
Accept-Language: ru
```

Backend responds with localized card data:
```go
func GetLocalizedCard(cardKey string, lang string) CardResponse {
    card := repository.FindCardByKey(cardKey)
    
    return CardResponse{
        Key:            card.Key,
        Name:           card.Names[lang],
        Interpretation: card.Interpretations[lang],
        Tags:           card.Tags[lang],
        ImageURL:       card.GetImageURL(),
    }
}
```

### 7.3 Date/Time Localization

**Frontend:**
```javascript
import { format } from 'date-fns'
import { enUS, ru } from 'date-fns/locale'

const formatDate = (date, language) => {
  const locale = language === 'ru' ? ru : enUS
  return format(date, 'PPPp', { locale })
}

// English: "March 31, 2026 at 3:00 PM"
// Russian: "31 марта 2026 г. в 15:00"
```

---

## 8. Theme System

**Reference:** theme tokens, semantic roles, component states, and reusable UI primitives must follow `specs/docs/DESIGN_SYSTEM.md`.

### 8.1 CSS Variables

**Base Structure:**
```scss
// src/styles/themes/_light.scss
:root {
  --ds-color-bg-page: #f5f0e6;
  --ds-color-bg-surface: #fffaf2;
  --ds-color-text-primary: #2f241d;
  --ds-color-text-secondary: #6e5a4e;
  --ds-color-accent-primary: #9f5c36;
  --ds-color-border-subtle: rgba(106, 74, 55, 0.16);
  --ds-motion-base: 300ms ease;
}

// src/styles/themes/_dark.scss
[data-theme="dark"] {
  --ds-color-bg-page: #000000;
  --ds-color-bg-page-alt: #0d0d0d;
  --ds-color-bg-surface: #1a1a1a;
  --ds-color-bg-elevated: #111111;
  --ds-color-text-primary: #ffffff;
  --ds-color-text-secondary: #c4c4c4;
  --ds-color-text-muted: #6b6b6b;
  --ds-color-accent-primary: #d4af37;
  --ds-color-accent-light: #f4e5b0;
  --ds-color-accent-dark: #9b7e1f;
  --ds-color-border-subtle: rgba(212, 175, 55, 0.18);
  --ds-color-border-accent: rgba(212, 175, 55, 0.45);
  --ds-motion-base: 300ms ease;
}
```

### 8.2 Theme Toggle

**Implementation:**
```javascript
// src/composables/useTheme.js
import { ref, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'

export function useTheme() {
  const authStore = useAuthStore()
  
  const theme = ref(
    localStorage.getItem('theme') || 
    authStore.user?.theme || 
    'light'
  )
  
  const setTheme = (newTheme) => {
    theme.value = newTheme
    document.documentElement.setAttribute('data-theme', newTheme)
    localStorage.setItem('theme', newTheme)
    
    // Update user profile if logged in
    if (authStore.isAuthenticated) {
      authStore.updateSettings({ theme: newTheme })
    }
  }
  
  // Apply theme on mount
  setTheme(theme.value)
  
  return {
    theme,
    setTheme,
    toggleTheme: () => setTheme(theme.value === 'light' ? 'dark' : 'light')
  }
}
```

**Usage in Component:**
```vue
<template>
  <button @click="toggleTheme">
    <IconSun v-if="theme === 'light'" />
    <IconMoon v-else />
  </button>
</template>

<script setup>
import { useTheme } from '@/composables/useTheme'

const { theme, toggleTheme } = useTheme()
</script>
```

---

## 9. Error Handling

### 9.1 Error Response Format

**Standard Error Response:**
```json
{
  "success": false,
  "error": "error_code",
  "message": "Human-readable error message",
  "details": {
    "field": "email",
    "reason": "Email is already registered"
  }
}
```

### 9.2 HTTP Status Codes

| Code | Usage |
|------|-------|
| **200** | Success |
| **201** | Resource created |
| **400** | Bad request / validation error |
| **401** | Unauthorized / invalid token |
| **403** | Forbidden / insufficient permissions |
| **404** | Resource not found |
| **429** | Rate limit exceeded |
| **500** | Internal server error |
| **501** | Not implemented (Gateway stub) |

### 9.3 Frontend Error Handling

**Axios Interceptor:**
```javascript
// src/services/api.js
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8888/api',
  timeout: 10000
})

// Request interceptor: add auth token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor: handle errors
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      const authStore = useAuthStore()
      authStore.logout()
      router.push('/login')
    }
    
    return Promise.reject(error)
  }
)

export default api
```

**Component-Level Error Handling:**
```javascript
// In component
import { ref } from 'vue'
import api from '@/services/api'

const error = ref(null)
const loading = ref(false)

async function generateReading() {
  error.value = null
  loading.value = true
  
  try {
    const response = await api.post('/reading/generate/card-of-day', {
      timezone: 'Europe/Amsterdam'
    })
    // Handle success
  } catch (err) {
    if (err.response?.data?.error === 'card_of_day_already_generated') {
      error.value = 'You already received your Card of the Day today'
    } else if (err.response?.data?.error === 'daily_limit_reached') {
      error.value = 'Daily limit reached. Try again tomorrow!'
    } else {
      error.value = 'Something went wrong. Please try again.'
    }
  } finally {
    loading.value = false
  }
}
```

---

## 10. Security

### 10.1 Password Security

**Hashing:** bcrypt with cost factor 12

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(hash), err
}

func VerifyPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 10.2 CORS Configuration

**Backend:**
```go
import "github.com/gin-contrib/cors"

func main() {
    router := gin.Default()
    
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
        AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept-Language"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
}
```

### 10.3 Input Validation

**Backend Validation:**
```go
import "github.com/go-playground/validator/v10"

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
    Password string `json:"password" validate:"required,min=8,password"`
    Timezone string `json:"timezone" validate:"required,timezone"`
}

// Custom validator for password
func passwordValidator(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    return hasLetter && hasNumber
}
```

**Frontend Validation (VeeValidate):**
```vue
<script setup>
import { useForm, useField } from 'vee-validate'
import * as yup from 'yup'

const schema = yup.object({
  email: yup.string().required().email(),
  username: yup.string().required().min(3).max(20).matches(/^[a-zA-Z0-9]+$/),
  password: yup.string()
    .required()
    .min(8)
    .matches(/[a-zA-Z]/, 'Must contain a letter')
    .matches(/[0-9]/, 'Must contain a number')
})

const { handleSubmit, errors } = useForm({ validationSchema: schema })
</script>
```

### 10.4 Rate Limiting

**TODO for future versions:**
- Implement rate limiting middleware (e.g., 100 requests/minute per IP)
- Use Redis for distributed rate limiting

---

## 11. Development Environment

### 11.1 Backend Setup

**Requirements:**
- Go 1.22+
- MongoDB 6.0+ (self-hosted, SSH tunnel)

**Environment Variables (`.env`):**
```bash
# Server
PORT=8888
GIN_MODE=debug

# Database
MONGO_URI=mongodb://localhost:27017
MONGO_DB=tarot_daily

# Auth
JWT_SECRET=your_jwt_secret_here_min_32_chars

# Frontend URL (for CORS)
FRONTEND_URL=http://localhost:5173

# File Upload
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE=2097152  # 2MB in bytes
```

**Run:**
```bash
cd backend
go mod download
go run main.go
```

### 11.2 Frontend Setup

**Requirements:**
- Node.js 18 LTS or 20
- npm 9+

**Environment Variables (`.env`):**
```bash
VITE_API_URL=http://localhost:8888/api
```

**Run:**
```bash
cd frontend
npm install
npm run dev
```

### 11.3 MongoDB Setup

**SSH Tunnel:**
```bash
ssh -L 27017:localhost:27017 user@your-mongo-server
```

**Database Initialization:**
```bash
# Connect to MongoDB
mongosh

# Create database
use tarot_daily

# Create collections (MongoDB creates them automatically on first insert)
# Seed data will be loaded via Go script
```

**Seed Script (Go):**
```go
// deploy/seed/main.go
func SeedDatabase() {
    // Load cards from JSON
    // Insert into cards collection
    // Create default deck
    // Create role policies
    // Create admin user
    // Create admin settings
}
```

---

## 12. Deployment Strategy

### 12.1 Build Process

**Backend:**
```bash
cd backend
go build -o tarot-daily-api main.go
```

**Frontend:**
```bash
cd frontend
npm run build
# Output: dist/
```

### 12.2 Production Environment Variables

**Backend:**
```bash
GIN_MODE=release
PORT=8888
MONGO_URI=mongodb://prod-host:27017
MONGO_DB=tarot_daily
JWT_SECRET=production_secret_min_32_chars
FRONTEND_URL=https://nebula-taro.com
```

**Frontend:**
```bash
VITE_API_URL=https://test-api.xmlab.ru/api
```

### 12.3 Deployment Checklist

**Pre-Deployment:**
- [ ] Run backend tests: `go test ./...`
- [ ] Run frontend build: `npm run build`
- [ ] Verify environment variables
- [ ] Backup MongoDB
- [ ] Test on staging environment

**Deployment:**
- [ ] Deploy backend binary
- [ ] Deploy frontend static files
- [ ] Run database migrations (if any)
- [ ] Verify API health: `GET /api/gateway/health`
- [ ] Smoke test critical flows (login, card generation)

**Post-Deployment:**
- [ ] Monitor error logs
- [ ] Check API response times
- [ ] Verify MongoDB connections
- [ ] Test file uploads

---

**End of Technical Specification**

**Next Document:** Implementation Roadmap (phased development plan)
