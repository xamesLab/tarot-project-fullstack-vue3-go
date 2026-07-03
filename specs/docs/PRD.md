# Product Requirements Document (PRD)
## Nebula Tarot — product baseline with v0.2.0 frontend overlay

**Document Version:** 1.1

**Last Updated:** July 3, 2026

**Status:** Working Baseline

> Version note: sections describing authenticated readings, history, catalogs and administration remain the product contract implemented by the `v0.1.x` Go + Vue stack. For `v0.2.0`, the active client moved to `frontend/` (React + TypeScript), while the previous Vue client moved to `frontend-legacy/`. The new client currently implements the public Card of the Day experience through mock APIs; see [`specs/v0.2.0`](../v0.2.0/README.md).

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Product Vision & Goals](#2-product-vision--goals)
3. [Target Audience](#3-target-audience)
4. [Feature Specifications](#4-feature-specifications)
5. [User Roles & Permissions](#5-user-roles--permissions)
6. [User Flows & Journeys](#6-user-flows--journeys)
7. [Internationalization](#7-internationalization)
8. [Success Metrics](#8-success-metrics)
9. [Out of Scope](#9-out-of-scope)
10. [Appendix](#10-appendix)

---

## 1. Executive Summary

### 1.1 Product Overview

**Nebula Tarot** is a web-based tarot reading platform designed for casual users, primarily women, seeking daily spiritual guidance through tarot card interpretations. The platform provides deterministic "card of the day" readings and random 3-card spreads with role-based access controls.

### 1.2 Business Model

- **Subscription-based** access tiers (post-MVP)
- **Role hierarchy** controls feature access and usage limits
- **Freemium approach** with progressive value unlock

### 1.3 Key Differentiators

- **Deterministic Card of the Day** with timezone-aware daily reset (3:00 AM local time)
- **Role-based progression system** (Seeker → Adept → Master → Keeper)
- **Bilingual support** (English/Russian) from day one
- **Clean, aesthetic UI** with light/dark theme support
- **Future AI integration** capability (Gateway stub in place)

### 1.4 MVP Scope

**In Scope:**
- User authentication and profile management
- Card of the Day (once per accounting day, 3:00 AM reset)
- 3-card spread (random, role-limited)
- Reading history (role-based retention)
- Deck and card back selection
- Admin panel for user/deck/policy management
- Bilingual interface (EN/RU)
- Light/Dark theme

**Out of Scope (Post-MVP):**
- Monetization/payment integration
- Email notifications
- AI-powered interpretations (Gateway stub only)
- Push notifications
- Social features
- 5-card spreads
- Custom spread creation
- Object storage (using local storage)

**Approved Frontend Reference:**
- MVP frontend implementation should follow the project design system in `specs/docs/DESIGN_SYSTEM.md`
- `frontend/src/pages/ui-kit/UiKitPage.tsx` and `/ui-kit` are the approved visual/component reference for `v0.2.0`
- the former `v0.1.x` reference is preserved at `frontend-legacy/src/pages/dev/UIKit.vue`

### 1.5 v0.2.0 Frontend Scope

`v0.2.0` introduces a public, low-friction entry point before the authenticated product journey:

- Home is available without authentication and presents a localized Card of the Day ritual.
- Card flip, interpretation spheres, Message of the Day and oracle chat are the primary engagement loop.
- Spreads, Horoscope and Learn are visible as gated navigation and open a registration modal.
- Registration and oracle calls use a typed mock API; they do not create a backend session.
- Theme and language preferences persist locally.
- Real auth, profile, history, catalogs and admin remain available only in the legacy client until migrated.

---

## 2. Product Vision & Goals

### 2.1 Vision Statement

> "To provide accessible, beautiful, and meaningful daily tarot guidance that fits seamlessly into modern life, empowering users to reflect, grow, and connect with their intuition."

### 2.2 Product Goals

**Primary Goals:**
1. Deliver a delightful, low-friction daily ritual for tarot enthusiasts
2. Build a scalable foundation for subscription-based growth
3. Establish a design-first approach with strong aesthetic appeal
4. Create a technically robust platform for future AI integration

**Secondary Goals:**
1. Support bilingual users without requiring language switching friction
2. Provide admin tools for content and user management
3. Maintain reading history for reflection and pattern recognition

### 2.3 Success Criteria (MVP)

- **User Activation:** 80% of registered users complete their first reading within 24 hours
- **Retention:** 40% of users return for their Card of the Day within 7 days
- **Technical:** Zero critical bugs, <2s page load time, mobile-responsive on iOS/Android
- **Content:** 78 cards with bilingual interpretations loaded and functional

---

## 3. Target Audience

### 3.1 Primary Persona: **The Casual Seeker**

**Demographics:**
- **Gender:** Female (85%)
- **Age:** 22-45
- **Location:** Urban/suburban, tech-savvy
- **Language:** English or Russian speakers

**Psychographics:**
- Interested in self-reflection, spirituality, and personal growth
- Not necessarily deep tarot practitioners; curious explorers
- Values aesthetics and user experience
- Uses mobile devices frequently
- Seeks daily rituals and mindfulness practices

**Behaviors:**
- Checks horoscopes, reads self-help content
- Uses wellness apps (meditation, journaling)
- Engages with social media (Instagram, Pinterest)
- Willing to pay for premium features that add value

**Motivations:**
- Daily guidance and clarity
- Self-reflection and journaling
- Aesthetic pleasure and ritual
- Curiosity about the future

**Pain Points:**
- Existing tarot apps feel cluttered or overly mystical
- Hard to remember meanings of cards
- No history/tracking of past readings
- Inconsistent quality of interpretations

### 3.2 Secondary Persona: **The Devoted Practitioner**

**Demographics:**
- **Gender:** Female (70%)
- **Age:** 30-55
- **Experience:** Intermediate to advanced tarot knowledge

**Psychographics:**
- Studies tarot seriously, owns physical decks
- Values accuracy and depth in interpretations
- Interested in multiple deck options
- Willing to pay for premium features

**Motivations:**
- Deeper insights and pattern tracking
- Access to multiple deck aesthetics
- Advanced spread options (future)
- Integration with journaling practice

**Pain Points:**
- Digital tools lack depth
- Want more control over spreads and timing
- Need better history and pattern analysis

---

## 4. Feature Specifications

### 4.1 Authentication & User Management

#### 4.1.1 Registration

**Requirements:**
- Email-based registration
- Password requirements: minimum 8 characters, must contain letter and number
- Username: 3-20 characters, may include letters, numbers, dots, underscores and hyphens
- Username is optional; if left blank, backend derives it from the email local-part
- Timezone selection (IANA timezone format) - required
- Default role: **Seeker**
- Default theme: **Light**
- Default language: Auto-detect from browser, fallback to English

**Fields:**
| Field | Type | Required | Validation |
|-------|------|----------|------------|
| Email | String | Yes | Valid email format, unique |
| Password | String | Yes | Min 8 chars, 1 letter + 1 number |
| Username | String | No | 3-20 chars, letters/numbers/dot/underscore/hyphen |
| Timezone | String | Yes | Valid IANA timezone |

**Success Flow:**
1. User fills registration form
2. Frontend validates inputs in real-time
3. Backend creates user with role `seeker`
4. User auto-logged in, redirected to dashboard
5. Welcome message shown (localized)

**Error Handling:**
- Email already exists → "This email is already registered"
- Username taken → "Username unavailable"
- Invalid timezone → "Please select a valid timezone"

#### 4.1.2 Login

**Requirements:**
- Email + password authentication
- JWT token issued on success
- Token stored in localStorage
- Token includes: user_id, role, timezone
- Session persistence across browser close

**Success Flow:**
1. User enters credentials
2. Backend validates and issues JWT
3. Frontend stores token
4. User redirected to dashboard
5. User profile loaded

**Error Handling:**
- Invalid credentials → "Email or password incorrect"
- Account disabled → "Your account has been suspended"

#### 4.1.3 Profile Management

**Editable Fields:**
| Field | Editable | Notes |
|-------|----------|-------|
| Username | Yes | Must remain unique |
| Email | No | Cannot be changed in MVP |
| Password | Yes | Requires current password confirmation |
| Timezone | Yes | Affects Card of the Day reset |
| Language | Yes | EN/RU |
| Theme | Yes | Light/Dark |
| Selected Deck | Yes | From available decks |
| Selected Card Back | Yes | From available card backs |

**Timezone Change Impact:**
- If timezone changes, Card of the Day may become available/unavailable based on new 3:00 AM boundary
- Backend recalculates "accounting day" based on new timezone
- Warning shown: "Changing timezone will affect when your Card of the Day refreshes"

---

### 4.2 Card of the Day

#### 4.2.1 Core Behavior

**Deterministic Generation:**
- Seed: `hash(user_id + accounting_date + timezone)`
- Accounting date calculated as:
  ```
  if local_time.hour >= 3:
      accounting_date = local_time.date
  else:
      accounting_date = local_time.date - 1 day
  ```
- One card per accounting day
- Card cannot be regenerated until next accounting day

**Card Selection:**
- Pulled from 78-card deck
- Can be any card (Major or Minor Arcana)
- No reversed cards in MVP

**Presentation:**
```
┌─────────────────────────────┐
│   Card of the Day           │
│                             │
│   ┌───────────────────┐     │
│   │                   │     │
│   │   [Card Image]    │     │
│   │                   │     │
│   └───────────────────┘     │
│                             │
│   The Fool                  │
│   "New beginnings..."       │
│                             │
│   Keywords: freedom, ...    │
│                             │
│   [Full Interpretation]     │
│                             │
│   Next card in: 8h 23m      │
└─────────────────────────────┘
```

**Reveal Behavior:**
- A newly generated card is initially shown face-down using the selected card back
- Title, interpretation, and tags are hidden until the user opens the card
- On click/tap, the card flips to its face and becomes `viewed`
- Viewed state is persisted per reading, so reloading the dashboard keeps already opened cards revealed

**Available Actions:**
- Reveal card contents
- Save to history (automatic)
- View in deck context (optional: show card position in deck)

**Constraints:**
- Cannot regenerate
- Cannot get reversed version
- Tied to user's timezone

#### 4.2.2 API Behavior

**Endpoint:** `POST /api/reading/generate/card-of-day`

**Request:**
```json
{
  "timezone": "Europe/Amsterdam"
}
```

**Response (First Request of Day):**
```json
{
  "success": true,
  "reading": {
    "id": "uuid",
    "type": "card_of_day",
    "generated_at": "2026-03-31T15:30:00Z",
    "accounting_date": "2026-03-31",
    "cards": [
      {
        "key": "the_fool",
        "name": "The Fool",
        "name_ru": "Шут",
        "interpretation": "New beginnings...",
        "interpretation_ru": "Новые начинания...",
        "tags": ["freedom", "adventure"],
        "tags_ru": ["свобода", "приключение"],
        "image_url": "/uploads/decks/default/the_fool.png"
      }
    ],
    "next_available_at": "2026-04-01T03:00:00+02:00"
  }
}
```

**Response (Already Generated Today):**
```json
{
  "success": false,
  "error": "card_of_day_already_generated",
  "message": "You already received your Card of the Day",
  "existing_reading": { ... },
  "next_available_at": "2026-04-01T03:00:00+02:00"
}
```

---

### 4.3 Three-Card Spread

#### 4.3.1 Core Behavior

**Random Generation:**
- Three unique cards selected randomly
- No predetermined positions/meanings in MVP
- Can be regenerated within role limits

**Spread Limits (Daily):**
- **Seeker:** 3 spreads/day
- **Adept:** 5 spreads/day
- **Master:** 10 spreads/day
- **Keeper:** Unlimited

**Admin Configuration:**
- Limits editable via Admin Panel
- Changes apply to new accounting days
- Cannot be changed retroactively

**Presentation:**
```
┌─────────────────────────────────────┐
│   Three-Card Spread                 │
│                                     │
│   ┌──────┐  ┌──────┐  ┌──────┐    │
│   │ Card │  │ Card │  │ Card │    │
│   │  1   │  │  2   │  │  3   │    │
│   └──────┘  └──────┘  └──────┘    │
│                                     │
│   [Interpretations below]           │
│                                     │
│   Spreads remaining today: 2/5      │
│                                     │
│   [Generate New Spread] [Save]      │
└─────────────────────────────────────┘
```

**Reveal Behavior:**
- New spread cards are initially hidden by the selected card back
- Each card opens independently on user click/tap
- Opening one card does not reveal the others
- Viewed state is persisted inside the reading and remains stable after reload

#### 4.3.2 Limit Enforcement

**Backend Logic:**
```
daily_spread_count = count(spreads WHERE 
  user_id = current_user 
  AND accounting_date = today 
  AND type = 'three_card'
)

user_role_limit = role_policies[user.role].three_card_limit

if daily_spread_count >= user_role_limit:
  return error: "daily_limit_reached"
```

**Frontend Display:**
- Show remaining spreads: "Spreads remaining today: 2/5"
- Disable "Generate" button when limit reached
- Show upgrade prompt (post-MVP) when limit hit

---

### 4.4 Reading History

#### 4.4.1 Retention Policy

**History Storage by Role:**
| Role | Retention | Notes |
|------|-----------|-------|
| Seeker | 1 day | Only today's readings |
| Adept | 7 days | Last week |
| Master | 15 days | Last 2 weeks |
| Keeper | 30 days | Last month |

**Cleanup:**
- Daily cron job deletes readings older than retention window
- User notified when readings are about to be deleted (future: email)

#### 4.4.2 History View

**Display:**
```
┌─────────────────────────────────────┐
│   Reading History                   │
│                                     │
│   [Filter: All | Card of Day | 3-Card] │
│                                     │
│   March 31, 2026                    │
│   ├─ Card of the Day: The Fool      │
│   │  Generated at 10:30 AM          │
│   │  [View Details]                 │
│   │                                 │
│   ├─ Three-Card Spread              │
│   │  Generated at 2:15 PM           │
│   │  Cards: The Magician, ...       │
│   │  [View Details]                 │
│   │                                 │
│   └─ Three-Card Spread              │
│      Generated at 5:45 PM           │
│      [View Details]                 │
│                                     │
│   March 30, 2026                    │
│   ├─ Card of the Day: Strength      │
│      ...                            │
│                                     │
│   Your history is kept for 7 days   │
└─────────────────────────────────────┘
```

**Functionality:**
- Click reading → full detail view with all cards and interpretations
- Filter by reading type
- Sort: newest first (default)
- Pagination: 10 readings per page

---

### 4.5 Deck & Card Back Selection

#### 4.5.1 Deck Management (User Perspective)

**Available Decks:**
- User sees all decks with status `done`
- Decks shown with:
  - Preview image (3-4 sample cards)
  - Name and description
  - Author (if applicable)

**Selection:**
- User selects deck from settings
- Deck selection is available only for roles `adept`, `master`, and `keeper`
- Selected deck saved to user profile
- All future readings use selected deck
- Past readings retain original deck (immutable)

**Fallback Behavior:**
- If user has not selected a custom deck, the system uses the global default deck configured by admin

**Default Deck:**
- "Tarot Classic" (loaded from provided JSON)
- Always available to all users

#### 4.5.2 Card Back Selection

**Available Card Backs:**
- Default card back always available
- Admin can upload additional card backs
- Card back is used in the live dashboard reading flow while cards remain unopened

**Selection:**
- Independent from deck selection
- User can mix: Deck A + Card Back B
- Saved to user profile
- Card back selection is available only for roles `adept`, `master`, and `keeper`

**Fallback Behavior:**
- If user has not selected a custom card back, the system uses the global default card back configured by admin

---

### 4.6 Admin Panel

#### 4.6.1 Access Control

**Role Required:** Keeper

**Route Protection:**
- Frontend: Route guard checks `user.role === 'keeper'`
- Backend: All `/api/admin/*` endpoints verify role

**Unauthorized Access:**
- Redirect to dashboard with error: "You do not have permission to access this area"

#### 4.6.2 Dashboard

**Metrics Displayed:**
```
┌─────────────────────────────────────┐
│   Admin Dashboard                   │
│                                     │
│   Total Users: 1,247                │
│   ├─ Seekers: 1,100                 │
│   ├─ Adepts: 120                    │
│   ├─ Masters: 25                    │
│   └─ Keepers: 2                     │
│                                     │
│   Readings Today: 3,842             │
│   ├─ Card of the Day: 1,247         │
│   └─ Three-Card Spreads: 2,595      │
│                                     │
│   Decks: 3 (2 active, 1 processing) │
│   Card Backs: 5                     │
└─────────────────────────────────────┘
```

#### 4.6.3 User Management

**User List:**
| Username | Email | Role | Status | Registered | Actions |
|----------|-------|------|--------|------------|---------|
| alice123 | alice@... | Seeker | Active | Mar 15 | [Edit Role] [Disable] |
| bob_tarot | bob@... | Adept | Active | Mar 20 | [Edit Role] [Disable] |

**Actions:**
- **Edit Role:** Change user role (Seeker/Adept/Master/Keeper)
- **Disable Account:** Set status to `disabled` (user cannot log in)
- **View Details:** See user's reading history, settings

**Search & Filter:**
- Search by username or email
- Filter by role
- Filter by status (active/disabled)
- Sort by registration date, last login

#### 4.6.4 Role Policy Management

**Policy Configuration:**
```
┌─────────────────────────────────────┐
│   Role Policies                     │
│                                     │
│   Seeker                            │
│   ├─ Three-Card Spread Limit: [3▾]  │
│   └─ History Retention: [1▾] day    │
│                                     │
│   Adept                             │
│   ├─ Three-Card Spread Limit: [5▾]  │
│   └─ History Retention: [7▾] days   │
│                                     │
│   Master                            │
│   ├─ Three-Card Spread Limit: [10▾] │
│   └─ History Retention: [15▾] days  │
│                                     │
│   Keeper                            │
│   ├─ Three-Card Spread Limit: ∞     │
│   └─ History Retention: [30▾] days  │
│                                     │
│   [Save Changes]                    │
└─────────────────────────────────────┘
```

**Constraints:**
- Limits must be positive integers
- Retention must be 1-365 days
- Changes apply to new accounting days only
- Confirmation required: "This will affect all users with this role starting tomorrow. Proceed?"

#### 4.6.5 Deck Management

**Deck List:**
```
┌─────────────────────────────────────┐
│   Decks                             │
│                                     │
│   [+ Create New Deck]               │
│                                     │
│   Tarot Classic (Default)           │
│   Status: Done | Cards: 78/78       │
│   [Edit] [View Cards]               │
│                                     │
│   Rider-Waite Modern                │
│   Status: Processing | Cards: 45/78 │
│   [Edit] [Upload Cards]             │
│                                     │
│   Ethereal Dreams                   │
│   Status: Done | Cards: 78/78       │
│   [Edit] [View Cards] [Disable]     │
└─────────────────────────────────────┘
```

**Create Deck Flow:**
1. Enter metadata:
   - Name (required)
   - Description (optional)
   - Author (optional)
   - Tags (optional, comma-separated)
2. Upload card back image separately (PNG/JPG, max 2MB, 2:3 aspect ratio)
3. Upload cards:
   - **Option A:** ZIP archive with structure containing only the 78 card images named by tarot `key`:
     ```
     deck-name.zip
       ├── 0_meta.json               // optional metadata autofill
       ├── the_fool.png
       ├── the_magician.png
       ├── ace_of_wands.jpg
       └── ... (78 cards total)
     ```
4. Backend validates:
   - All 78 cards present
   - File names match known tarot `key` values
   - Valid image formats
   - File sizes under allowed limit
   - `0_meta.json` may be present and is ignored by backend validation
5. Deck status set to `processing`
6. Backend processes images (save to `/uploads/decks/{deck_id}/`)
7. Deck status set to `done`
8. Deck becomes available to users

**Archive and Delete:**
- Active decks and card backs can be moved to archive from admin UI
- Archived assets are not available in user/admin selection lists
- Archived assets can be restored
- Permanent deletion is available only from archive
- Current global defaults cannot be archived or deleted until another default is assigned

**Admin Defaults and Manual Selection:**
- Admin can upload decks and card backs independently
- Admin can choose the global default deck and the global default card back independently
- These defaults apply to all users who have not chosen their own deck or card back
- Users with roles `adept`, `master`, and `keeper` can override the global defaults from profile settings

#### 4.6.6 Global Settings

**Settings:**
```
┌─────────────────────────────────────┐
│   Global Settings                   │
│                                     │
│   Maintenance Mode                  │
│   [○ Disabled  ● Enabled]           │
│   Message: "Site under maintenance" │
│                                     │
│   Default Role for New Users        │
│   [Seeker ▾]                        │
│                                     │
│   Registration Enabled              │
│   [✓] Allow new registrations       │
│                                     │
│   [Save Settings]                   │
└─────────────────────────────────────┘
```

---

### 4.7 Gateway Integration (Stub)

#### 4.7.1 Current State

**Implementation:** `v0.1.0` used a stub/placeholder; `v0.1.1` switches to direct AI integration through `OpenRouter`

**Frontend:**
- AI interpretation is shown in the existing `Meaning` / `Значение расклада` block after the reading is revealed
- UI keeps the existing reading screen instead of creating a separate AI page

**Backend:**
- `POST /api/ask` creates or reuses an AI interpretation request
- `GET /api/ask/:request_id` returns status and completed result
- `POST /api/gateway/callback` endpoint exists (no-op)
- `GET /api/gateway/health` returns `{ status: "stub" }`

**Database:**
- `ask_requests` stores AI request history and provider metadata
- the latest successful AI result is additionally cached on the reading itself for fast reuse

**Future Integration:**
- Will be added via separate instruction
- Frontend will enable "AI Interpretation" button when Gateway is live
- Role limits will apply (e.g., Seeker: 0, Adept: 3/day, Master: 10/day)

**Planning update for `v0.1.1`:**
- the next version uses direct AI integration through `OpenRouter` as the primary path
- the direct AI response is expected to populate the existing `Meaning` / `Значение расклада` block after a card of the day or a fully revealed spread is opened
- the prompt must include the revealed cards and require only a compact interpretation capped at `100` words
- the legacy gateway path remains as an architectural placeholder for a future fallback design

---

## 5. User Roles & Permissions

### 5.1 Role Hierarchy

```
Keeper (Admin)
    ↑
  Master (Premium)
    ↑
  Adept (Paid)
    ↑
  Seeker (Free)
```

### 5.2 Permission Matrix

| Feature | Seeker | Adept | Master | Keeper |
|---------|--------|-------|--------|--------|
| **Authentication** |
| Register/Login | ✅ | ✅ | ✅ | ✅ |
| Edit Profile | ✅ | ✅ | ✅ | ✅ |
| Change Password | ✅ | ✅ | ✅ | ✅ |
| Change Timezone | ✅ | ✅ | ✅ | ✅ |
| **Readings** |
| Card of the Day | ✅ (1/day) | ✅ (1/day) | ✅ (1/day) | ✅ (1/day) |
| Three-Card Spread | ✅ (3/day) | ✅ (5/day) | ✅ (10/day) | ✅ (∞) |
| **History** |
| View History | ✅ (1 day) | ✅ (7 days) | ✅ (15 days) | ✅ (30 days) |
| Export History | ❌ | ❌ | ❌ | ❌ |
| **Customization** |
| Select Deck | ❌ (default only) | ✅ | ✅ | ✅ |
| Select Card Back | ❌ (default only) | ✅ | ✅ | ✅ |
| Theme (Light/Dark) | ✅ | ✅ | ✅ | ✅ |
| Language (EN/RU) | ✅ | ✅ | ✅ | ✅ |
| **AI Features (Future)** |
| AI Interpretation | ❌ | ✅ (3/day) | ✅ (10/day) | ✅ (∞) |
| Chat with AI | ❌ | ✅ (5 msg/reading) | ✅ (20 msg/reading) | ✅ (∞) |
| **Admin Panel** |
| Access Admin Panel | ❌ | ❌ | ❌ | ✅ |
| Manage Users | ❌ | ❌ | ❌ | ✅ |
| Manage Decks | ❌ | ❌ | ❌ | ✅ |
| Edit Role Policies | ❌ | ❌ | ❌ | ✅ |
| Global Settings | ❌ | ❌ | ❌ | ✅ |

### 5.3 Role Assignment

**Default:** New users are assigned `Seeker` role

**Change Role:**
- Only Keeper can change user roles
- Changes take effect immediately
- User limits updated on next reading generation

**Future (Post-MVP):**
- Seeker → Adept via payment
- Adept → Master via payment
- Master → Keeper: manual assignment only (not purchasable)

---

## 6. User Flows & Journeys

### 6.1 New User Registration & First Reading

```
┌─────────────────────────────────────┐
│ 1. Landing Page                     │
│    → [Sign Up] button               │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ 2. Registration Form                │
│    - Email                          │
│    - Username                       │
│    - Password                       │
│    - Timezone (auto-detected)       │
│    → [Create Account]               │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ 3. Dashboard (Auto-login)           │
│    Welcome message:                 │
│    "Welcome, Alice! Get your first  │
│     Card of the Day"                │
│    → [Get Card of the Day] (CTA)    │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ 4. Card of the Day Generated        │
│    - Card displayed                 │
│    - Interpretation shown           │
│    - Tags/keywords visible          │
│    "Your next card in 14h 32m"      │
│    → [Try 3-Card Spread]            │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ 5. First 3-Card Spread              │
│    - 3 cards displayed              │
│    - Interpretations shown          │
│    "Spreads remaining: 2/3"         │
│    → [Generate Another] or [Done]   │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ 6. Return Tomorrow                  │
│    User completes activation        │
│    Success metric: returned for 2nd │
│    Card of the Day                  │
└─────────────────────────────────────┘
```

**Key Success Indicators:**
- Time to first reading < 2 minutes
- Completion rate of registration flow > 85%
- Return rate for Day 2 Card of the Day > 40%

---

### 6.2 Daily Ritual User (Returning User)

```
┌─────────────────────────────────────┐
│ 1. User Opens App                   │
│    Logged in automatically (JWT)    │
│    → Dashboard loads                │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ 2. Dashboard View                   │
│    Status shown:                    │
│    ✓ Card of the Day ready          │
│    OR                               │
│    ⏰ Next card in 5h 12m            │
│    → [Get Today's Card] (if ready)  │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ 3. Card of the Day                  │
│    - Card revealed                  │
│    - User reads interpretation      │
│    - Saved to history automatically │
│    → [Generate 3-Card Spread?]      │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ 4. Optional: 3-Card Spread          │
│    User may generate 0-3 spreads    │
│    based on interest/limits         │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ 5. Optional: View History           │
│    User reflects on past readings   │
│    Sees patterns over time          │
└─────────────────────────────────────┘
```

**Typical Session:**
- Duration: 3-8 minutes
- Pages viewed: 2-4
- Actions: Get card, read interpretation, maybe generate 1 spread

---

### 6.3 Admin User Journey

```
┌─────────────────────────────────────┐
│ 1. Keeper Logs In                   │
│    → Dashboard shows admin link     │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ 2. Admin Panel Home                 │
│    - Overview metrics               │
│    - Quick actions                  │
│    → Choose task                    │
└────────────┬────────────────────────┘
             ↓
    ┌────────┼────────┬────────────┐
    ↓        ↓        ↓            ↓
┌───────┐ ┌─────┐ ┌──────┐ ┌──────────┐
│ Users │ │Decks│ │Roles │ │ Settings │
└───┬───┘ └──┬──┘ └───┬──┘ └────┬─────┘
    ↓        ↓        ↓          ↓
```

**Common Admin Tasks:**

**A. Upload New Deck:**
1. Admin → Decks → Create New Deck
2. Fill metadata (name, description)
3. Upload ZIP file with 78 cards + back
4. System validates structure
5. Deck status: Processing → Done
6. Deck available to users

**B. Change User Role:**
1. Admin → Users → Search user
2. Click [Edit Role]
3. Select new role (Seeker/Adept/Master)
4. Confirm change
5. User's limits updated immediately

**C. Adjust Role Policies:**
1. Admin → Role Policies
2. Edit spread limits (e.g., Adept: 5→7)
3. Save changes
4. Changes apply to all users with that role

---

### 6.4 Edge Cases & Error Flows

#### 6.4.1 User Hits Spread Limit

```
┌─────────────────────────────────────┐
│ User: "Generate 3-Card Spread"      │
│ System: daily_count >= limit        │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ Error Message:                      │
│ "You've reached your daily limit    │
│  of 3 spreads. Your limit resets    │
│  tomorrow at 3:00 AM."              │
│                                     │
│  [View History] [OK]                │
└─────────────────────────────────────┘
```

**Future (Post-MVP):**
- Add upgrade CTA: "Upgrade to Adept for 5 spreads/day"

#### 6.4.2 User Tries to Regenerate Card of the Day

```
┌─────────────────────────────────────┐
│ User: Clicks "Get Card of the Day"  │
│ System: Card already generated      │
└────────────┬────────────────────────┘
             ↓
┌─────────────────────────────────────┐
│ Show Existing Card:                 │
│ "You already received your Card of  │
│  the Day. Your next card will be    │
│  available tomorrow at 3:00 AM."    │
│                                     │
│  [View Card] [Generate 3-Card?]     │
└─────────────────────────────────────┘
```

#### 6.4.3 Timezone Change Affects Card Availability

**Scenario:** User in UTC+2 changes timezone to UTC-5

```
Before: 
  Local time: 10:00 AM (UTC+2)
  Card of Day: Already generated today

After Timezone Change:
  Local time: 3:00 AM (UTC-5) 
  → Still same accounting day
  → Card remains the same
  
BUT if they changed at 2:00 AM (UTC+2):
  → Would become 7:00 PM yesterday (UTC-5)
  → New accounting day!
  → Card of Day becomes available
```

**System Behavior:**
- Recalculate accounting day on timezone change
- If new accounting day → Card of Day available
- Show notification: "Your Card of the Day schedule has updated based on your new timezone"

---

## 7. Internationalization

### 7.1 Supported Languages

**MVP:** English (EN), Russian (RU)

**Language Selection:**
- `v0.2.0` defaults to Russian and falls back to Russian
- User switches language in the global header
- Preference is saved to `localStorage` as `nebula-tarot:lang`
- Profile synchronization remains a post-HTTP-integration requirement

### 7.2 Translated Content

**Frontend (UI):**
- Navigation labels
- Buttons and CTAs
- Form labels and placeholders
- Error messages
- Validation messages
- Help text

**Backend (Content):**
- Card names (e.g., "The Fool" / "Шут")
- Card interpretations
- Card keywords/tags

**Not Translated:**
- Usernames
- User-generated content (future: notes, questions)
- Email addresses
- Technical error codes

### 7.3 Implementation

**Frontend (i18next + react-i18next):** locale resources live in `frontend/src/features/i18n/locales/{ru,en}/`. `common` owns header, footer and registration copy; `home` owns Card of the Day and oracle copy. Reading content is returned by the selected typed API implementation for the requested language rather than hardcoded in components.

**Backend (MongoDB):**
```javascript
// Card document structure
{
  "_id": "the_fool",
  "key": "the_fool",
  "arcana": "major",
  "names": {
    "en": "The Fool",
    "ru": "Шут"
  },
  "interpretations": {
    "en": "New beginnings, innocence, spontaneity...",
    "ru": "Новые начинания, невинность, спонтанность..."
  },
  "tags": {
    "en": ["freedom", "adventure", "new start"],
    "ru": ["свобода", "приключение", "новое начало"]
  }
}
```

**API Response (Localized):**
```javascript
// Request header: Accept-Language: ru
// Response automatically includes Russian content

{
  "card": {
    "key": "the_fool",
    "name": "Шут",
    "interpretation": "Новые начинания...",
    "tags": ["свобода", "приключение"]
  }
}
```

### 7.4 Date/Time Localization

**Display Format:**
- English: "March 31, 2026 at 3:00 PM"
- Russian: "31 марта 2026 в 15:00"

**Implementation:**
```javascript
// Frontend (using day.js or date-fns)
import { format } from 'date-fns';
import { enUS, ru } from 'date-fns/locale';

const locale = currentLanguage === 'ru' ? ru : enUS;
const formatted = format(date, 'PPPp', { locale });
```

---

## 8. Success Metrics

### 8.1 Key Performance Indicators (KPIs)

**User Acquisition:**
- **Registration Completion Rate:** >85% of users who start registration complete it
- **Time to First Reading:** <2 minutes from registration to first card

**Engagement:**
- **Daily Active Users (DAU):** Track daily login rate
- **Card of the Day Completion:** >70% of active users get their daily card
- **Return Rate (D1):** >40% of users return next day
- **Return Rate (D7):** >25% of users return after 7 days

**Retention:**
- **Week 1 Retention:** >40% of users return within first week
- **Week 2 Retention:** >25% of users active in week 2
- **Monthly Active Users (MAU):** Track monthly engagement

**Feature Usage:**
- **3-Card Spread Usage:** % of users who generate at least one 3-card spread
- **Average Spreads per Active User:** Track engagement depth
- **History View Rate:** % of users who view their reading history

**Technical:**
- **Page Load Time:** <2s for dashboard
- **API Response Time:** <500ms for p95
- **Error Rate:** <1% of API requests fail
- **Uptime:** >99.5% monthly

### 8.2 Success Criteria for MVP Launch

**Must Hit:**
- ✅ All core features functional (auth, readings, history, admin)
- ✅ Zero critical bugs
- ✅ Mobile-responsive on iOS Safari and Android Chrome
- ✅ Both languages (EN/RU) complete and accurate
- ✅ 78 cards loaded with interpretations
- ✅ Admin can upload and manage decks

**Nice to Have:**
- 🎯 50 beta users testing before public launch
- 🎯 <5 minor bugs reported in beta
- 🎯 Positive feedback on UI/UX from beta testers

---

## 9. Out of Scope

### 9.1 Deferred to Post-MVP

**Monetization:**
- Stripe/PayPal integration
- Subscription management
- Payment processing
- Upgrade flows

**Communication:**
- Email notifications (verification, password reset, reading reminders)
- Push notifications (mobile)
- In-app notifications

**AI Features:**
- Gateway integration (stub exists, but not functional)
- AI-powered interpretations
- Chat with AI about readings

**Advanced Features:**
- 5-card spreads
- Custom spread creation
- Reversed card interpretations
- Spread position meanings (e.g., Past/Present/Future)
- Export reading history (PDF, JSON)
- Social features (sharing, community)
- Journaling with readings

**Infrastructure:**
- Object storage (S3/R2) for deck images
- CDN for image delivery
- Email service integration
- Advanced analytics (Mixpanel, Amplitude)
- A/B testing framework

**User Experience:**
- Onboarding tutorial
- Interactive card meanings guide
- Card of the day reminder emails
- Streak tracking and rewards
- Achievement system

### 9.2 Explicitly Excluded

**Never Planned:**
- Real-time multiplayer readings
- Video/audio content
- Marketplace for third-party decks
- API for third-party integrations

---

## 10. Appendix

### 10.1 Glossary

**Terms:**
- **Accounting Day:** The 24-hour period from 3:00 AM to 2:59 AM next day in user's local timezone
- **Major Arcana:** 22 trump cards (The Fool through The World)
- **Minor Arcana:** 56 suit cards (Wands, Cups, Swords, Pentacles)
- **Spread:** A reading with multiple cards
- **Card of the Day:** Single-card deterministic daily reading
- **Reversed Card:** Card drawn upside-down (not supported in MVP)
- **Reading:** A complete tarot consultation (Card of the Day or Spread)
- **Deck:** A complete set of 78 tarot card images
- **Card Back:** The design shown on the back of cards (cosmetic)

### 10.2 Card Structure Reference

**78 Cards Total:**
- **22 Major Arcana:**
  - 0: The Fool
  - I: The Magician
  - II: The High Priestess
  - ... (through XXI: The World)

- **56 Minor Arcana:**
  - **Wands (14 cards):** Ace, 2-10, Page, Knight, Queen, King
  - **Cups (14 cards):** Ace, 2-10, Page, Knight, Queen, King
  - **Swords (14 cards):** Ace, 2-10, Page, Knight, Queen, King
  - **Pentacles (14 cards):** Ace, 2-10, Page, Knight, Queen, King

### 10.3 Timezone Calculation Example

**User Timezone:** Europe/Amsterdam (UTC+2 in summer)  
**Current UTC Time:** 2026-03-31 01:30:00 UTC  
**Local Time:** 2026-03-31 03:30:00 (UTC+2)

**Accounting Day Calculation:**
```
local_hour = 03 (>= 3)
→ accounting_date = 2026-03-31

Card of Day seed: hash(user_id + "2026-03-31" + "Europe/Amsterdam")
```

**Edge Case (Before 3 AM):**
**UTC Time:** 2026-03-31 00:30:00 UTC  
**Local Time:** 2026-03-31 02:30:00 (UTC+2)

```
local_hour = 02 (< 3)
→ accounting_date = 2026-03-30 (previous day)

Card of Day seed: hash(user_id + "2026-03-30" + "Europe/Amsterdam")
```

### 10.4 References

**Tarot Card Data:**
- Source: Uploaded JSON file `tarot-cards.json`
- Contains 78 cards with EN/RU translations
- Structure: `id`, `key`, `name`, `name_en`, `interpretation`, `interpretation_en`, `tags`, `tags_en`, `arcana`, `type`

**Design Inspirations:**
- Labyrinthos Tarot (aesthetic, clean UI)
- Co-Star (minimalism, daily ritual)
- Golden Thread Tarot (journaling, history)

---

**End of Product Requirements Document**

**Next Steps:**
1. Review and approve PRD
2. Create Technical Specification
3. Create Implementation Roadmap
4. Begin development

**Document Owner:** Product Team  
**Stakeholders:** Engineering, Design, Business  
**Review Cycle:** Monthly or as needed

---
