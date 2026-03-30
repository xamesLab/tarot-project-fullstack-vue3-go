# Domain Model Summary

Last updated: March 30, 2026.

This document summarizes the main business entities currently represented in backend models and MongoDB collections.

## Core Entities

### User

Represents an authenticated account.

Important fields:

- `login`
- `passwordHash`
- `role`
- `status`
- `displayName`
- `timezone`
- `preferredDeckId`
- `deckThemeId` for legacy compatibility

Notes:

- default product role is `user`
- admin-only functionality is guarded by `role == admin`
- `status` supports `active`, `banned`, and `freeze`

### Card

Represents a tarot catalog entry.

Important fields:

- `code`
- `name`
- `arcanaType`
- `description`
- `interpretationOpen`
- `interpretationClosed`

Notes:

- cards are the canonical interpretation source
- deck assets only change visuals, not meaning

### DailyReading

Represents one user-day reading result.

Important fields:

- `userId`
- `date`
- `timezoneUsed`
- `deckThemeId`
- `cardOfDay`
- `spreadA`
- `spreadB`
- `spreadGenerationsCount`

Notes:

- uniqueness is enforced per `userId + date`
- card draws store `cardId` plus `state`
- `state` is `open` or `closed`

### DeckTheme

Legacy theme model used for mapping card asset URLs.

Important fields:

- `name`
- `isPlaceholder`
- `assets`
- `createdBy`

Notes:

- deck themes remain in the codebase for MVP compatibility
- the product is also moving toward deck-based selection flows

### Deck

Represents an uploadable and admin-managed deck entity.

Important fields:

- `name`
- `status`
- `createdByUserID`
- `backImageUrl`
- `cards`
- `cardsUploadedCount`
- `requiredCardsCount`

Status values:

- `process`
- `done`

Notes:

- user-facing deck selection should only expose `done` decks
- current upload storage is local filesystem storage

### RolePolicy

Defines generation limits by role.

Important fields:

- `role`
- `dailySpreadLimit`

Notes:

- `-1` means unlimited
- policies support admin-defined throttling of spread generation

### AskRequest

Represents an async Gateway-backed ask operation.

Important fields:

- `requestId`
- `userId`
- `text`
- `status`
- `resultText`
- `errorCode`
- `errorMessage`
- `completedAt`

Status values:

- `accepted`
- `queued`
- `processing`
- `success`
- `error`

### AdminSettings

Represents singleton-style global admin settings.

Important fields:

- `projectTitle`
- `supportEmail`
- `maintenanceMode`
- `dayResetHour`
- `maxUploadMb`

## Main Relationships

- one `User` has many `DailyReading`
- one `DailyReading` references many `Card` documents through draw entries
- one `User` may reference one preferred `Deck`
- one `User` may still reference one legacy `DeckTheme`
- one `User` may create many `Deck` or `DeckTheme` records in admin flows
- one `User` may create many `AskRequest` records

## Practical Modeling Notes

- the product currently mixes older `deck_theme` compatibility with newer `deck` flows
- uploaded deck assets and static fallback card assets are both active runtime concepts
- the domain model is stable enough for development, but `v0.1.2` validation and planned `v0.1.4` hygiene work still remain
