# PRD: Tarot Daily Web App (MVP)

## 1. Product Overview

### 1.1 Working Name
Tarot Daily

### 1.2 Goal
An entertainment web app in the Tarot domain where a user gets daily:
- 2 tarot spreads
- 1 card of the day

Each card includes:
- image
- description
- interpretation/prediction

### 1.3 Business Problem
The target audience wants quick, personalized mystical content 1-2 times a day in a convenient personal account.

## 2. Target Audience

- Girls and women interested in mysticism and tarot.
- Typical behavior: short sessions 1-2 times per day to get a forecast/entertainment content.

## 3. MVP Scope

### 3.1 In Scope (MVP)
- Sign-up/sign-in with username and password.
- User personal account.
- User roles: `user`, `admin`.
- Admin page, accessible only to `admin` role.
- Generation and display of:
  - 2 tarot spread types
  - 1 card of the day
- Random card draw from the deck.
- Card state support: `open` or `closed`.
- Display card image, description, and interpretation.
- User can select a card image set/theme (placeholders in MVP).
- Responsive interface (mobile-first).
- Local run in dev environment.

### 3.2 Out of Scope (MVP)
- Production deployment and production-grade infrastructure.
- Payments, subscriptions, CRM integrations.
- Mandatory AI-generated interpretations (low priority, optional).
- Full CMS: in MVP, admin can only add image set variants (basic minimum).

## 4. User Scenarios

1. As a new user, I want to register with username/password to use my personal account.
2. As an authenticated user, I want to open my account and get the card of the day.
3. As an authenticated user, I want to get 2 spreads and see drawn cards with interpretations.
4. As a user, I want to choose a visual card image set in my account settings.
5. As a user, I want to return 1-2 times per day and get daily content.
6. As an admin, I want to open the admin page and add new image set variants.

## 5. Functional Requirements

### 5.1 Authentication and Roles
- Registration:
  - unique username
  - password (stored as hash)
- Login with username/password.
- Logout.
- Session via JWT (or equivalent token mechanism).
- User role in profile:
  - default `user`
  - `admin` created/assigned manually at dev stage

### 5.2 User Personal Account
- Main screen after login:
  - card of the day (one card)
  - two spreads (type A and type B for MVP)
- User settings:
  - timezone (default `UTC+00:00`)
  - selected card image set (`deckTheme`)
- Last generation date block.
- Re-login on the same day shows the same result.

### 5.3 Admin Panel (MVP)
- Dedicated page `/admin`.
- Access only for `admin` role (backend + frontend guard).
- MVP admin features:
  - add a new card image set variant
  - view list of available sets
- Non-admin user receives 403/redirect.

### 5.4 Card Content
- System contains a deck card catalog (minimum: 78 cards).
- For each card:
  - `name`
  - `description`
  - `interpretationOpen`
  - `interpretationClosed`
- Images depend on selected `deckTheme`; placeholders are allowed in MVP.
- Ability to replace/extend interpretations with AI in future.

### 5.5 Generation Rules
- Card of the day: 1 card.
- Spread A: 3 cards.
- Spread B: 5 cards.
- Random draw from deck without duplicates inside the same spread.
- Duplicates between different blocks on the same day are allowed (MVP simplification).
- Each drawn card gets a random state: `open` or `closed`.
- Interpretation is selected by card state.
- Daily logic:
  - content is fixed for the current date in user's timezone from account settings
  - if timezone is not set, use `UTC+00:00`
  - on the next day (in that timezone), generate a new result

### 5.6 History (Minimum)
- Store results by date for each user.
- In MVP UI, it is acceptable to show only the current day (history prepared at DB/API level for extension).

## 6. Non-Functional Requirements

- Responsiveness: correct behavior on mobile screens (priority).
- Performance: key screen loads within 2-3 seconds in local dev with typical data volume.
- Security (basic):
  - password hashing (bcrypt/argon2)
  - basic input validation
  - private endpoint protection through auth middleware
  - `admin` role check for admin endpoints
- Codebase should remain simple for a small team.

## 7. Tech Stack and Architecture

### 7.1 Frontend
- Vue 3
- Pinia
- Vue Router
- Responsive layout (mobile-first)

### 7.2 Backend
- Go
- Gin
- REST API
- JWT auth + role-based access middleware

### 7.3 Database
- MongoDB
- Collections: `users`, `cards`, `daily_readings`, `deck_themes`

### 7.4 Local Development
- Local run only (no production deployment)
- Prefer `docker-compose` for MongoDB + backend + frontend (optional for MVP)

## 8. Data Model (Draft)

### 8.1 users
- `_id`
- `login` (unique)
- `passwordHash`
- `role` (`user` | `admin`)
- `timezone` (IANA or fixed UTC offset, default `UTC+00:00`)
- `deckThemeId` (selected image set variant)
- `createdAt`
- `updatedAt`

### 8.2 cards
- `_id`
- `name`
- `arcanaType` (major/minor)
- `description`
- `interpretationOpen`
- `interpretationClosed`

### 8.3 deck_themes
- `_id`
- `name`
- `isPlaceholder` (bool)
- `assets` (mapping: `cardCode -> imageUrl`)
- `createdBy` (admin userId)
- `createdAt`

### 8.4 daily_readings
- `_id`
- `userId`
- `date` (YYYY-MM-DD in user's timezone)
- `timezoneUsed`
- `cardOfDay` (`cardId`, `state`)
- `spreadA` (array of `{ cardId, state }`)
- `spreadB` (array of `{ cardId, state }`)
- `generatedAt`
- unique index: (`userId`, `date`)

## 9. API (Draft Contract)

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/me`
- `PATCH /api/me/settings` (timezone, deckTheme)
- `GET /api/reading/today`
  - returns card of the day + 2 spreads for current date
  - creates and returns if no record exists for that date
- `GET /api/deck-themes` (for users to choose a theme)
- `POST /api/admin/deck-themes` (admin only)
- `GET /api/admin/deck-themes` (admin only)

Example `GET /api/reading/today` response:

```json
{
  "date": "2026-03-07",
  "timezoneUsed": "UTC+03:00",
  "cardOfDay": {
    "id": "card_10",
    "name": "Wheel of Fortune",
    "state": "open",
    "imageUrl": "/themes/default/wheel.jpg",
    "description": "...",
    "interpretation": "..."
  },
  "spreadA": [],
  "spreadB": []
}
```

## 10. UX/UI Requirements (MVP)

- Personal account screen:
  - "Card of the Day" block
  - "Spread A (3 cards)" block
  - "Spread B (5 cards)" block
- Each card shows:
  - image
  - name
  - state (`open`/`closed`)
  - description
  - interpretation
- Settings section:
  - timezone selection
  - image set selection
- Admin page:
  - image set list
  - form to add a new set
- Mobile-first:
  - comfortable vertical scroll layout
  - minimum tap target size 44px

## 11. MVP Success Metrics

- Functional criterion: product runs locally and passes core user flow.
- Technical criterion: authentication + daily content delivery work stably.
- Behavioral criterion (manual testing stage): user can return 1-2 times daily and see fixed result for the date.

## 12. Acceptance Criteria

1. User can register and login with username/password.
2. After login, user sees card of the day and 2 spreads (3 and 5 cards).
3. Cards are selected randomly and without duplicates inside each spread.
4. Every card has `open` or `closed` state, and matching interpretation is shown.
5. Each card displays image, description, and interpretation.
6. Re-login on same day returns same data; next day (by user timezone) returns new data.
7. User can change timezone in personal account; default is `UTC+00:00`.
8. User can choose an available card image set.
9. Admin page is available only to `admin`; regular user has no access.
10. Admin can add a new image set variant.
11. Interface is correctly displayed on mobile resolution.
12. Project runs locally in dev.

## 13. Post-MVP Roadmap

- Integrate AI-generated interpretations (feature flag).
- Extend admin panel to full content management.
- Add reading history in UI.
- Add themes/deck personalization.

## 14. Finalized Decisions

1. Card states in MVP: `open` and `closed`.
2. Spread formats: A = 3 cards, B = 5 cards.
3. Daily generation uses timezone from personal account; default `UTC+00:00`.
4. Card images:
   - placeholders are used in MVP
   - user chooses a variant
   - additional variants are added by admin via admin panel
