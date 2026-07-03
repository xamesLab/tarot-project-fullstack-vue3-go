# NEBULA-TAROT

Frontend v1 for the NEBULA-TAROT card-of-the-day experience. The application recreates the supplied Home and UI Library prototypes with React, TypeScript, a shared design system, mock APIs, RU/EN localization and dark/light themes.

## Setup

Requirements: Node.js 18+ (Node.js 20 LTS is recommended) and npm.

```bash
npm install
npm run dev
```

Open `/` for Home and `/ui-kit` for the development-only component showcase.

## Scripts

- `npm run dev` — start Vite.
- `npm run build` — strict TypeScript check and production build.
- `npm run preview` — preview the production bundle.
- `npm run typecheck` — TypeScript only.
- `npm run lint` — ESLint.
- `npm test` — Vitest smoke tests.
- `npm run extract:references` — unpack the original prototype templates into ignored `reference/` files.

## Configuration

`VITE_API_MODE=mock|http` selects the API implementation. Both development and production use `mock` until the backend exists; the HTTP adapter intentionally throws a clear not-implemented error.

Theme and language are persisted in localStorage as `nebula-tarot:theme` and `nebula-tarot:lang`. Add design values to both theme blocks in `src/shared/styles/tokens.css`. Add locale strings under `src/features/i18n/locales/<lang>/`; localized reading content belongs to mock/API fixtures.

## Architecture

- `src/shared` — tokens, icons, hooks and reusable controlled UI components.
- `src/features` — preferences, registration and oracle chat behavior.
- `src/widgets` — page-level composition.
- `src/pages` — Home and dev-only UI Kit routes.
- `src/api` — stable typed API surface with mock and future HTTP implementations.

New reusable visuals should first be implemented in `shared/ui` and demonstrated on `/ui-kit`. The original `Arcana Home.html` and `Arcana UI Library.html` files are immutable visual references; their legacy names do not represent the product brand.
