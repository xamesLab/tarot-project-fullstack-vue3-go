# Frontend

Client application for the Tarot Daily project.

## Stack

- Vue 3 + TypeScript
- Pinia
- Vue Router
- Vite

## Setup

1. Install dependencies:
   - `npm install`
2. Run dev server:
   - `npm run dev`
3. Build production bundle:
   - `npm run build`
4. Preview production bundle locally:
   - `npm run preview`

## Scripts

- `npm run dev` - start Vite dev server
- `npm run build` - run type-check and build
- `npm run preview` - preview production build
- `npm run typecheck` - run `vue-tsc --noEmit`

## Environment

Use template:
- `frontend/.env.example`

Local file:
- `frontend/.env`

Variables:
- `VITE_API_BASE_URL` - backend base URL (`/` in dev via Vite proxy)
- `VITE_APP_NAME` - app display name

## Dev Proxy

`vite.config.ts` proxies:
- `/api` -> `http://127.0.0.1:8080`
- `/assets/cards` -> `http://127.0.0.1:8080`
- `/uploads` -> `http://127.0.0.1:8080`

## Routing and Access

Main routes:
- public-only: `/login`, `/register`
- authenticated: `/dashboard`, `/profile`, `/settings`
- admin-only: `/admin/*`
- fallback: `/:pathMatch(.*)*`

Guards are implemented in:
- `src/router/guards.ts`

### Dev-only routes

Dev tools routes are registered via:
- `src/router/dev-tools.routes.ts`

They are connected conditionally in:
- `src/router/index.ts` using `import.meta.env.DEV`

Current dev-only route:
- `/dev` (`src/pages/dev/DevIndexPage.vue`)
- `/ui-kit` (`src/pages/dev/UIKit.vue`)

Additional runtime protection:
- routes with `meta.devOnly` are redirected to `not-found` outside dev mode.

## Background System

Global background is implemented as a separate component:
- `src/components/AppBackground.vue`

Integration point:
- `src/App.vue`

Notes:
- background is purely decorative and non-interactive (`pointer-events: none`)
- content layer stays above background (`app-shell` at higher z-index)
- light/dark appearance is synchronized with app theme

## UI Kit Page

`/ui-kit` is a standalone Arcana UI reference page for internal development:
- sticky header + theme toggle
- sticky section navigation + smooth scroll
- active section highlight via `IntersectionObserver`
- 11 sections: brand, palette, typography, iconography, buttons, forms, cards, navigation, dividers, motion, do/don't
- improved mobile harmony (`<=760px` and `<=480px`)

## Dev Index

`/dev` is the dev-only entry point for internal tools:
- current catalog entry: `/ui-kit`
- future dev-only pages can be linked from here
- protected by the same `meta.devOnly` routing rules as `/ui-kit`

## Autotests Status

Все задачи, связанные с автотестами и visual/e2e checks, отложены до отдельного явного запроса пользователя.

Это означает:
- не считать Playwright и любые e2e/visual regression задачи текущим активным этапом;
- не расширять и не чинить автотестовый контур без отдельной команды;
- приоритет остается за ручной валидацией и функциональными изменениями.

## Mode Switcher

`ModeSwitcher` and `ModeToggle` are aligned with UI Kit section 6:
- language segmented control (EN/RU)
- theme circular toggle (34x34)
- symbols follow UI Kit visual language
- `Icon` component remains in markup as SVG fallback hook

## UI Rules

Project-level UI policy is documented in:
- `frontend/UI_KIT_RULE.md`

Agent-level rules are documented in:
- `AGENTS.md`
- `frontend/AGENTS.md`

Mandatory highlights:
- `UI Kit First`: for components/layout/styles, use `/ui-kit` as source of truth
- `Background Freeze Rule`: do not modify background markup/styles without explicit user request

## Related Specs

- `specs/version-0.1.3/front/requirements.md`
- `specs/version-0.1.3/front/step-by-step.md`
- `specs/version-0.1.3/front/arcana-uikit-prompt.md`
- `specs/version-0.1.3/front/arcana-background-prompt.md`
