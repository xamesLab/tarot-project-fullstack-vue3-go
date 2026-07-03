# NEBULA-TAROT — Recommendations for Further Development

Guidance for the agent (and future contributors) on how to evolve the codebase after v1. These are not v1 tasks; they are constraints and directions to keep in mind so v1 decisions don't block them.

---

## 1. Working conventions

- **The UI kit page is the contract.** Any new visual element must first appear as a `shared/ui` component demonstrated on `/ui-kit`, then be used in pages. Never style one-off elements inside pages; if a page needs something new, promote it to the library.
- **Design tokens are the only source of color/spacing truth.** Never hardcode hex values in components. New colors go into `tokens.css` (both themes at once) and the Tailwind theme mapping. Adding a theme later (e.g. "midnight", "solar") must be a matter of adding one variable block and one toggle option.
- **Decomposition discipline:** if a component needs a scroll to read, split it. Widgets compose, `shared/ui` renders, features own state. Keep `shared/ui` free of store/i18n/router imports permanently.
- **Every SVG is a component** with `currentColor` where the prototype uses it, so theming keeps working.
- Keep the `reference/` extracted prototypes until the design stabilizes; delete before public release.
- Conventional Commits + a `CHANGELOG.md` once more than one person/agent works on the repo.

## 2. Architecture evolution

- **API layer:** when the backend appears, implement the `http` variant of the `api` object with the existing DTOs; add an anti-corruption mapping layer (`api DTO → entity`) so backend schema churn doesn't ripple into UI. Keep react-query keys centralized (`src/api/queryKeys.ts`).
- **Auth:** replace the stub store with a real session flow: token storage (httpOnly cookie preferred over localStorage), `api.auth.me()` bootstrap, a `RequireAuth` route wrapper, and a global 401 interceptor that opens the existing registration/login modal. The current "nav opens registration modal" behavior generalizes into a `useGatedAction(action)` hook: run the action if authenticated, otherwise open the modal and continue after success.
- **Routing:** future pages (Spreads `/spreads`, Horoscope `/horoscope`, Learn `/learn`, card encyclopedia `/cards/:slug`, profile `/me`) should be lazy routes with route-level code splitting. Nav items switch from "open modal" to real links once pages exist; keep the gating hook for premium actions.
- **i18n scaling:** move to lazy-loaded namespaces per page as locales grow; add an ICU/plural strategy early (Russian plural rules will bite otherwise); consider extracting keys with `i18next-parser` in CI to catch missing translations. If more languages come, add a language detector but keep RU/EN pinned in the header switcher design.
- **State:** keep Zustand slices small (theme, lang meta, auth, ui/modals). Server data stays in react-query — resist mirroring it into Zustand.
- **Modals:** if modal count grows, add a tiny modal manager (single `ui` slice: `openModal(id, props)`) so stacking, scroll lock, and focus restore live in one place.

## 3. Content & domain

- Model the full tarot domain early: `Card { id, arcana: 'major'|'minor', suit?, number, roman, names: {ru,en}, keywords, meanings: { upright, reversed, spheres } }`. The daily reading, spreads, and the learning section all consume this one model — design it before building Spreads.
- The oracle chat will likely become an LLM-backed feature: keep the chat feature isolated behind `api.oracle.ask()` with a streaming-friendly interface (`onToken?` callback or async iterable) so switching from canned replies to streaming responses doesn't rewrite the UI.
- Card illustrations: the prototype uses inline SVG art. Decide on a pipeline (per-card SVG components vs. sprite vs. server-hosted images) before the deck grows to 78 cards; lazy-load card art.

## 4. Quality, testing, tooling

- **Tests:** Vitest + React Testing Library for shared/ui behavior (modal focus trap, toggle semantics, form validation) and store logic; Playwright for 3 smoke E2E flows (load home, switch theme+lang, open/submit registration modal). Add Playwright visual snapshots of `/ui-kit` — it doubles as visual regression coverage of the whole design system.
- **Storybook is optional:** the dev-only `/ui-kit` page already serves the purpose; only add Storybook if you need isolated interaction testing or docs for external consumers. Don't maintain both showcase systems.
- **CI:** typecheck + lint + unit tests + build on every PR; Playwright on main. Enforce bundle-size budget (e.g. `rollup-plugin-visualizer` + a size-limit check) — fonts and card art are the main risks.
- **A11y:** run axe on Home and the modal in CI. Known watch-points: contrast of `--gold` on `--ink` in light theme, decorative SVGs need `aria-hidden`, the flip card needs a keyboard path (Enter/Space) and an announced state.

## 5. Performance & delivery

- Self-host fonts (`@fontsource` or downloaded woff2) with `font-display: swap` and only the weights actually used; the prototype ships ~17 font files — trim to the used subset (latin + cyrillic, the 3 families, few weights).
- The star field and orbs are cheap, but keep them `pointer-events: none`, `position: fixed`, and behind `prefers-reduced-motion`. If profiling ever shows paint cost, cap star count on mobile.
- Preload the daily reading query on app start; consider `staleTime: Infinity` for a "renews at dawn" resource with a scheduled invalidation at local midnight/dawn.
- When SEO matters (marketing pages, card encyclopedia), evaluate migrating pages to SSR/SSG (e.g. moving to Next/Remix or prerendering); the strict shared/ui isolation makes this migration mostly a routing/layout move. Until then, at least set meta tags, OG images, and localized `<html lang>`.

## 6. Product-fit details worth adding soon

- Persist the flipped/expanded state of the daily card per day (localStorage keyed by date) so revisits feel continuous.
- Real date/moon-phase computation for the DateStrip (`Intl.DateTimeFormat` + a small moon-phase lib) instead of the mocked string.
- Birth-date field from registration should feed a zodiac profile used in the `UserPill` and future horoscope page — store it in the user entity from day one.
- Error/empty states in the brand language (an ornamented "the cards are silent" panel) — design them as shared components on `/ui-kit` first, per the contract in §1.
- Analytics events behind a thin `track()` facade (no vendor lock-in): card flip, expand, oracle question, registration funnel steps.

## 7. Things explicitly deferred from v1 (do not sneak them in early)

- Real authentication/session handling, password reset, OAuth.
- Spreads/Horoscope/Learn pages.
- LLM oracle backend and streaming.
- Payments/premium gating beyond the registration modal.
- SSR, PWA/offline, push notifications ("card is ready at dawn" is a natural future push, note it in the backlog).
