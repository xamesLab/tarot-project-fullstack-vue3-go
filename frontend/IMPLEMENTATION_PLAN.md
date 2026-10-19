# NEBULA-TAROT — Frontend Implementation Plan (v1)

> **Historical implementation brief.** The initial React foundation described here has been completed and subsequently extended with backend-connected Card of the Day, profile, spreads, admin and production deployment. Use `frontend/README.md` and `specs/docs/AI_PROJECT_MAP.md` for current behavior; retain this file as the original design/implementation contract.

You are building the first version of the **NEBULA-TAROT** web frontend (a tarot / "card of the day" product) from two HTML prototype files located in the same folder. The source prototype filenames retain their original Arcana names, but all product-facing copy uses NEBULA-TAROT:

- `Arcana Home.html` — the Home page ("Card of the Day")
- `Arcana UI Library.html` — the UI component library page (the design-system showcase)

Your deliverable is a new project created **in the same folder as these HTML files**, implementing both pages **pixel-faithfully**, with a scalable architecture, mock API stubs, theme switching (dark/light), i18n (RU/EN), and a registration modal.

---

## 0. How to read the prototype files (IMPORTANT — do this first)

The two `.html` files are **not plain HTML**. They are self-unpacking bundles: the real page markup, CSS, and behavior script are stored as JSON inside `<script type="__bundler/manifest">` (base64/gzip assets: fonts + runtime) and `<script type="__bundler/template">` (a JSON string containing the full page HTML).

Two ways to study them:

1. **Open in a browser** — each file unpacks itself and renders the live page. Use this to observe the final look, hover states, and animations.
2. **Extract the source** (recommended, gives you exact CSS and logic). Example extraction script:

```python
import json, re, base64, gzip

for name in ['Arcana Home', 'Arcana UI Library']:
    html = open(f'{name}.html', encoding='utf-8').read()
    template = json.loads(re.search(
        r'<script type="__bundler/template">(.*?)</script>', html, re.S).group(1))
    open(f'{name}_extracted.html', 'w', encoding='utf-8').write(template)
```

Inside the extracted template you will find:
- All page CSS in `<style>` blocks (design tokens, component styles, `@keyframes`).
- The page markup with `{{ mustache }}` placeholders bound to a small state class.
- A `<script type="text/x-dc" data-dc-script>` block — a `Component extends DCLogic` class. **This is the prototype's interaction logic** (state, handlers, i18n dictionary, mount effects). Read it carefully: it defines exactly which elements are stateful and how.

Key facts already known from the prototypes (verify them yourself):

- **Design tokens (dark theme, `:root`):** `--gold:#C9A96E; --gold-light:#E8D5A3; --gold-dim:#8A6F47; --ink:#0A0A0F; --ink-mid:#13131C; --ink-soft:#1C1C2A; --mist:#E8E2D9; --mist-dim:#9B9189; --cream:#F5F0E8; --line:rgba(201,169,110,.12); --line-soft:rgba(201,169,110,.1)`
- **Light theme** is a `.light` class on the app root that overrides the same variables (e.g. `--gold:#8A6A2C; --ink:#F0E8D9; --mist:#221C2D; ...`). A few components have extra `.light`-scoped overrides (e.g. `.light .echo`). Port every one of them.
- **Fonts:** `Cinzel` (display/logo/headings), `Cormorant Garamond` (serif body, leads, card names), `Raleway` (UI text, labels, buttons). Load via Google Fonts or `@fontsource`. **Include both `latin` and `cyrillic` subsets** — the site is RU-first. Note: Cinzel has no Cyrillic glyphs; check what the prototype falls back to for Cyrillic headings and replicate the same font stack order so the fallback is identical.
- **CSS keyframe animations (Home):** `twinkle` (stars), `float` (orbs, `6s ease-in-out infinite`), `spin` (zodiac wheel, `120s linear infinite`).
- **JS-driven animations:**
  - Home: star field is generated at mount — 70 `div.star` elements with randomized size (0.5–2.5px), position, `--dur` (3–8s), `--delay` (0–4s), opacity (0.1–0.4).
  - Home: tarot card **3D flip** on click (`.flipped` class, `transform-style: preserve-3d`, `transition: transform .85s cubic-bezier(...)`).
  - UI Library: tarot card **tilt-on-hover** (mousemove → `rotateY(±10deg) rotateX(±10deg) scale(1.02)`, smooth return on mouseleave with a temporary `0.6s cubic-bezier(0.23,1,0.32,1)` transition). Only enabled when `matchMedia('(hover: hover)')` matches.
  - Modal: overlay open/close via `.open` class + CSS transitions; closes on overlay click and on `Escape`; inner panel stops propagation.
  - Home: "Expand interpretation" accordion (`.open` class on button + expandable spheres block).
- **The prototype's i18n dictionary** (`DICT = { ru: {...}, en: {...} }` in the Home script) contains all Home page strings in both languages. Reuse these strings verbatim as the seed for your translation files.

---

## 1. Goals and scope of v1

1. **Home page (`/`)** — a verbatim recreation of `Arcana Home.html`: layout, decor SVGs, animations, dark/light theme, RU/EN switch, card flip, expandable interpretation, oracle chat (mocked), all of it.
2. **UI Kit page (`/ui-kit`)** — a verbatim recreation of `Arcana UI Library.html`, but built **from the real shared components** of the project (the page is the living showcase of the component library, not a copy-pasted mockup). Available **only in dev mode**.
3. Home opens **without authentication**. Clicking any of the 4 header nav items ("Card of the Day", "Spreads", "Horoscope", "Learn") opens a **registration modal** (see §8).
4. Theme switching (dark default / light) and language switching (RU default / EN) work globally and persist.
5. All data goes through a **mock API layer** that can later be swapped for a real backend without touching UI code.
6. Architecture must support future growth: more pages, real auth, real API.

Out of scope for v1: real backend, real auth/sessions, routing to Spreads/Horoscope/Learn pages (nav items open the registration modal instead), SSR, tests beyond a basic smoke setup.

---

## 2. Tech stack

| Concern | Choice | Notes |
|---|---|---|
| Build | **Vite** + **TypeScript** | strict mode on |
| UI | **React 18+** | function components + hooks only |
| Styling | **Tailwind CSS** + CSS variables + a small amount of plain CSS | see §4 for the strategy |
| State | **Zustand** | small, scalable, slices per domain |
| Routing | **react-router** (v6/v7) | data-router style |
| i18n | **i18next + react-i18next** | RU/EN namespaces, prepared for lazy loading |
| Forms | **react-hook-form + zod** | registration form validation |
| Server state (future) | **@tanstack/react-query** | install and wire the provider now; mock services return promises so migration is trivial |
| Lint/format | ESLint + Prettier | plus `eslint-plugin-react-hooks` |

Do not add a component framework (MUI, shadcn, etc.) — the design system is fully bespoke.

---

## 3. Project structure

Create the project directly in the same folder as the HTML files; keep the prototypes untouched.

Feature-oriented structure:

```
frontend/
├─ index.html
├─ vite.config.ts
├─ tailwind.config.ts
├─ tsconfig.json                # paths: "@/*" -> "src/*"
├─ .env.development             # VITE_API_MODE=mock
├─ .env.production
└─ src/
   ├─ main.tsx
   ├─ app/
   │  ├─ App.tsx                # providers composition
   │  ├─ router.tsx             # routes, dev-only ui-kit registration
   │  └─ providers/             # ThemeProvider, I18nProvider, QueryProvider
   ├─ shared/                   # the design system — everything reusable
   │  ├─ styles/
   │  │  ├─ tokens.css          # :root + .light variable sets (verbatim from prototype)
   │  │  ├─ fonts.css
   │  │  └─ animations.css      # @keyframes twinkle/float/spin + shared transitions
   │  ├─ ui/                    # one folder per component (see §5)
   │  │  ├─ Button/ (Button.tsx, index.ts)
   │  │  ├─ Chip/  Badge/  Field/  Toggle/  Checkbox/  RadioCard/
   │  │  ├─ SegmentedControl/  Avatar/  UserPill/  NavBar/
   │  │  ├─ Ornament/  Divider/  SectionLabel/
   │  │  ├─ TarotCard/          # face, back, flip, tilt variants
   │  │  ├─ AspectBlock/  CornerPanel/  Tooltip/
   │  │  ├─ Modal/              # overlay, head, body, foot, close
   │  │  └─ Logo/
   │  ├─ icons/                 # every decorative SVG as a component
   │  ├─ lib/                   # hooks: useTilt, useEscape, useLockBodyScroll, usePrefersReducedMotion
   │  └─ config/                # env access, constants
   ├─ entities/                 # domain models & mock data
   │  └─ card/                  # TarotCard type, cardOfDay mock, keywords, spheres
   ├─ features/
   │  ├─ theme/                 # store slice + ThemeToggle component
   │  ├─ i18n/                  # store/init + LangSwitch component + locales/
   │  │  └─ locales/ (ru/common.json, ru/home.json, en/…)
   │  ├─ auth/                  # RegistrationModal + form + auth store stub
   │  └─ oracle-chat/           # chat UI + mock reply service
   ├─ widgets/                  # page-level composites
   │  ├─ Header/                # logo + nav + user pill + lang + theme toggle
   │  ├─ StarField/  Orbs/  ConstellationDecor/  ZodiacWheel/
   │  ├─ DateStrip/  MoonPhases/
   │  ├─ CardOfDayStage/        # side runes, echoes, glow, TarotCard, hint
   │  ├─ Interpretation/        # card label, keywords, brief, expand, spheres grid
   │  └─ MessageOfDay/
   ├─ pages/
   │  ├─ home/HomePage.tsx
   │  └─ ui-kit/UiKitPage.tsx   # + section components per showcase block
   └─ api/
      ├─ client.ts              # fetch wrapper (base URL, errors) — unused by mocks but ready
      ├─ types.ts
      └─ mock/                  # mock implementations (see §9)
```

Rules:
- **No component over ~150 lines.** If a component grows, decompose it.
- `shared/ui` components know nothing about the store, routing, or i18n — they take props (including text) and are fully controlled or minimally self-contained.
- `widgets` may consume stores and translations.
- Barrel exports (`index.ts`) per component folder.

---

## 4. Styling strategy (Tailwind + tokens)

The prototype's visual identity lives in **CSS custom properties** and a lot of bespoke CSS (gradients, borders, ornaments, 3D transforms). Do not try to translate everything into raw Tailwind utilities — you will lose fidelity. Use this hybrid:

1. Copy the `:root` and `.light` variable blocks **verbatim** into `shared/styles/tokens.css`. These variables are the single source of truth.
2. Extend Tailwind's theme to reference them:
   ```ts
   colors: {
     gold: 'var(--gold)', 'gold-light': 'var(--gold-light)', 'gold-dim': 'var(--gold-dim)',
     ink: 'var(--ink)', 'ink-mid': 'var(--ink-mid)', 'ink-soft': 'var(--ink-soft)',
     mist: 'var(--mist)', 'mist-dim': 'var(--mist-dim)', cream: 'var(--cream)',
     line: 'var(--line)', 'line-soft': 'var(--line-soft)',
   },
   fontFamily: { display: ['Cinzel', 'serif'], serif: ['"Cormorant Garamond"', 'serif'], ui: ['Raleway', 'sans-serif'] }
   ```
3. Use Tailwind utilities for layout/spacing/typography where they map cleanly; keep **component-scoped CSS (CSS modules or a co-located `.css` file)** for complex pieces the prototype styles with multi-property rules: the tarot card and its frames, echoes and glow, corner-ornament panels, moon-phase strip, toggles/knobs, modal transitions, tooltip bubbles. Copy those rules from the extracted prototype CSS **as-is**, only renaming class names to match your components.
4. Port all `@keyframes` into `animations.css` unchanged.
5. Theme switching = toggling the `.light` class on the app root (`<html>` or top-level wrapper) — exactly like the prototype. Nothing else should change per theme.

**Fidelity check:** every visual detail (adapted letter-spacing on `NEBULA-TAROT`, roman numeral `II`, opacity of the zodiac wheel `.13`, blur radius of orbs `120px`, star opacity ranges, border radii, box shadows) must match the prototype. When in doubt, copy the prototype's exact CSS value rather than approximating with a Tailwind step.

---

## 5. Component inventory (from the UI Library page)

The UI Library page contains 16 sections. Each maps to shared components you must build. Build them exactly as shown, with the props below (extend as the markup requires):

1. **Foundations — Palette** — showcase-only (swatch grid `sw/sw-chip/sw-hex/sw-name`); implement as a UI-kit-page component that reads tokens.
2. **Typography** — showcase of the three families (`t-cinzel`, `t-corm`, `t-rale` rows with meta).
3. **Buttons** — `Button` with `variant: 'primary' | 'outline' | 'ghost'`, `size: 'md' | 'sm'`, `iconOnly?: boolean`, disabled state.
4. **Chips & Badges** — `Chip` (`solid?: boolean`), `Badge` (`variant: 'gold' | 'soft' | 'dim'`, optional dot).
5. **Input fields** — `Field` (label + input/textarea, focus ring in gold, placeholder styling).
6. **Switchers** — `Toggle` (switch with knob), `Checkbox` (`cbx`), `RadioCard` (`rad`/`rad-opt`/`rad-dot` — the spread selector cards), `SegmentedControl` (`seg`/`seg-opt` — used for Direct/Reversed and reused in the header for RU/EN).
7. **Avatars & User Pill** — `Avatar` (`size: sm|md|lg`, `ring?: boolean`, initials), `UserPill` (avatar + name + zodiac line).
8. **Navigation** — `NavBar` (active item state; on Home its items trigger the registration modal).
9. **Ornaments & Dividers** — `Ornament` (star flourish), `Divider` (`divlbl` with/without label, dashed variant).
10. **Tarot card** — `TarotCard`: props for `roman`, `name`, `illustration`, `flippable?: boolean` (Home: click-to-flip with back side `cb-frame`), `tilt?: boolean` (UI kit: hover tilt via a `useTilt` hook). Face and back are sub-components.
11. **Aspect blocks** — `AspectBlock` (icon + title + text; used as the 4 "spheres" on Home).
12. **Corner-ornament panel** — `CornerPanel` (the `panel` with 4 `corner tl/tr/bl/br` ornaments; used for Message of the Day on Home).
13. **Tooltips & hints** — `Tooltip` (`tip`/`tip-bubble`, placements incl. `below`, and the `info-dot` trigger + rich popover variant `pop-title/pop-text`).
14. **Modal** — `Modal` (overlay fade, panel, `ModalHead/Body/Foot`, close button, Escape + overlay-click close, body scroll lock, focus trap). This is the base for the registration modal.
15. **Icons & Decor** — every inline SVG from both pages becomes a named component in `shared/icons` (logo sigil, star ornament, moon phases, zodiac wheel, constellations, aspect icons, close icon, send icon, etc.). No SVG duplication in page code.
16. **Logo** — `Logo` (`size: sm|md|lg`, lockup with wordmark + tagline).

The **UiKitPage** rebuilds the showcase layout of the prototype (rail/sections with `sec-num`, `sec-title`, `sec-desc`, demo rows) and demonstrates each real component, including the interactive states that the prototype wires up (toggles, checkboxes, segmented control, radio cards, nav active state, modal open, theme toggle on that page).

---

## 6. Home page composition

Recreate `Arcana Home.html` from widgets (top to bottom):

1. Background layer: `StarField` (70 randomized stars, generated once on mount), three `Orbs`, two fixed `ConstellationDecor` SVGs.
2. `Header`: logo group, `NavBar` (4 items; **each click opens RegistrationModal**, active state stays on "Card of the Day"), `UserPill` (static mock user "Мария А. / Maria A."), RU/EN `SegmentedControl`, theme toggle button (☾/☀).
3. `DateStrip`: `MoonPhases` SVG strip + localized date line.
4. Ornament + section label ("Card of the Day").
5. `CardOfDayStage`: zodiac wheel (slow `spin`), side runes left/right, card glow, echo cards, flippable `TarotCard` (High Priestess, II), "tap to flip" hint.
6. `Interpretation`: card label (roman + title + subtitle), 4 keyword `Chip`s, brief text, ornament, expand button → animated reveal of 4 `AspectBlock` spheres (Love, Career, Health, Spirit).
7. `MessageOfDay`: `CornerPanel` with the two-paragraph message (contains an `<em>` — keep rich text via i18n `<Trans>` or markup-in-translation).
8. `OracleChat`: header, scrollable log (oracle greeting seeded per language), input + send (Enter submits), mock oracle reply; auto-scroll to bottom on new message. When the language changes, the seeded greeting message is replaced with the other language's greeting (replicate the prototype's `switchLang` behavior).
9. Footer line ("The cards renew each day at dawn").

All content (card of the day, spheres, message, oracle replies) must come from the **mock API / entities layer**, not be hardcoded in components — this is a deliberate difference from the prototype, where it was inlined.

---

## 7. Theme & language

- **Theme store** (Zustand slice): `theme: 'dark' | 'light'`, default `dark`, persisted in `localStorage` (`nebula-tarot:theme`). Applies/removes `.light` on the root element via an effect.
- **Language**: i18next with `ru` (default) and `en`, persisted in `localStorage` (`nebula-tarot:lang`). Namespaces at minimum: `common` (header, footer, modal, form) and `home`. Seed all strings from the prototype's `DICT` verbatim. Structure translation keys semantically (`home.spheres.love.title`), not `sph1t`.
- Everything user-visible is translated — including the registration modal, validation errors, and the UI-kit page chrome may stay as-is (it's dev-only; translating it is optional but the components it shows must render translated content when given).
- Date line: for v1 you may localize the mock date via i18n values; design the `DateStrip` to accept a `Date` so it can later be generated with `Intl.DateTimeFormat`.

---

## 8. Registration modal (new — not in the prototypes)

Triggered by clicking any header nav item on Home (and any future gated action). Built from the shared `Modal` styled exactly like the UI-library modal (head/body/foot, gold-bordered panel, overlay fade).

Content (design it in the same visual language — Cinzel title, Cormorant lead, gold `Field`s and `Button`s, ornament divider):

- Title: "Join NEBULA-TAROT" / «Присоединиться к NEBULA-TAROT» (put in i18n).
- Lead sentence inviting registration to unlock spreads/horoscope/learning.
- Form (`react-hook-form` + `zod`): **name**, **email**, **password** (min 8), optional **birth date** (used later for the natal/zodiac features — a nice product-fitting touch). Localized inline validation errors.
- Footer: primary `Button` "Create account" + ghost/outline "Later" (closes). Optional secondary line "Already have an account? Sign in" — render it as a stub link (no login flow in v1).
- Submit calls `api.auth.register(payload)` (mock: ~600ms delay, returns success), shows a success state inside the modal (ornament + "Welcome" message), then closes. Store the "registered" flag in the auth store stub (no real session).
- Accessibility: focus trap, initial focus on first field, `Escape`/overlay close, `aria-modal`, restore focus to the trigger on close.

---

## 9. Mock API layer

- `src/api/types.ts`: domain DTOs — `TarotCardDto`, `DailyReadingDto` (card, keywords, brief, spheres, message), `OracleReplyDto`, `RegisterPayload`, `UserDto`.
- `../mocks/frontend/`: implementations returning typed data via `Promise` + artificial `delay(300–800ms)`. Content = the prototype's texts, **in both languages** (mock accepts a `lang` param or returns both locales; pick one approach and document it).
- Public surface: a single `api` object (`api.readings.getDaily()`, `api.oracle.ask(question)`, `api.auth.register(data)`), selected by `VITE_API_MODE` (`mock` | `http`). The `http` variant can throw `NotImplemented` for now — the point is that swapping is a one-file change.
- Wire `@tanstack/react-query` provider and fetch the daily reading through `useQuery` on Home (with a simple loading state that doesn't break the layout — e.g. render the stage immediately and fill texts when ready, or a subtle skeleton in the brand style). Oracle chat uses a mutation.

---

## 10. Routing & dev-only UI kit

- Routes: `/` → HomePage; `/ui-kit` → UiKitPage.
- The UI-kit route must exist **only in dev**:
  ```ts
  const routes = [
    { path: '/', element: <HomePage /> },
    ...(import.meta.env.DEV
      ? [{ path: '/ui-kit', lazy: () => import('@/pages/ui-kit') }]
      : []),
  ];
  ```
  Lazy import + the `import.meta.env.DEV` guard ensures the page is both unreachable **and tree-shaken out of the production bundle**. Add a catch-all 404 that redirects to `/`.
- In dev, add a small unobtrusive link or console hint to `/ui-kit` if you want, but do not alter the Home design.

---

## 11. Animations checklist (must all be ported)

- [ ] Star twinkle (`@keyframes twinkle`, per-star `--dur`/`--delay` custom properties).
- [ ] Orb float (`float 6s ease-in-out infinite`, staggered per orb if the prototype staggers them).
- [ ] Zodiac wheel rotation (`spin 120s linear infinite`).
- [ ] Tarot card 3D flip on Home (`.flipped`, `preserve-3d`, `.85s` cubic-bezier, backface visibility).
- [ ] Tarot card hover tilt on UI kit (`useTilt` hook; pointer-fine only; smooth spring-back).
- [ ] Card glow / echo cards static effects.
- [ ] Expand/collapse of the interpretation spheres (height/opacity transition + button chevron/label state).
- [ ] Modal overlay + panel enter/leave transitions.
- [ ] All hover transitions on buttons, chips, nav items, toggles (copy transition durations/easings from prototype CSS).
- [ ] Respect `prefers-reduced-motion`: disable twinkle/float/spin/tilt via a media query or the `usePrefersReducedMotion` hook (flip and modal can remain, shortened).

Implement animations with CSS (as the prototype does) — do **not** introduce framer-motion for v1.

---

## 12. Quality bar & acceptance criteria

1. `npm run dev` serves the app; `/` matches `Arcana Home.html` side-by-side in both themes and both languages (open the prototype in a browser next to your app and compare at 1280px and ~1440px widths; also verify a mobile width if the prototype includes responsive rules — port them if so).
2. `/ui-kit` matches `Arcana UI Library.html` and every interactive demo works (toggles, checkboxes, radio cards, segmented, nav, modal, tilt, theme toggle).
3. `npm run build` succeeds; the production bundle contains **no** ui-kit page code; visiting `/ui-kit` in a production preview shows 404/redirect.
4. Theme and language persist across reloads; switching either updates the whole page instantly with no visual glitches (including the seeded oracle greeting swap).
5. Nav clicks on Home open the registration modal; the form validates, submits to the mock, shows success, closes; focus management works.
6. No hardcoded user-facing strings outside locale files (except the dev-only ui-kit chrome).
7. TypeScript strict passes; ESLint passes; no `any` in the api/entities layer.
8. No component file exceeds ~150 lines; SVGs are components; no duplicated markup between Home and UI kit — both consume the same `shared/ui` components.
9. `README.md` in the project root: setup, scripts, env vars, folder structure overview, how to add a locale/theme token/component, how the mock↔http API switch works.

## 13. Suggested implementation order

1. Extract prototype sources (§0); keep the extracted files in a `reference/` folder (git-ignored) for side-by-side comparison.
2. Scaffold Vite + TS + Tailwind + ESLint/Prettier + aliases.
3. Tokens, fonts, keyframes, theme mechanism (`.light` class + store + persistence).
4. i18n setup + locales seeded from `DICT`.
5. Shared UI components bottom-up (icons → primitives → composite: modal, tarot card).
6. UiKitPage (this validates the whole component library early), with the dev-only route guard.
7. Mock API + entities.
8. Home widgets and page assembly; animations; oracle chat.
9. Registration modal + auth stub + nav wiring.
10. Fidelity pass (side-by-side), reduced-motion, a11y pass, README, build check.
