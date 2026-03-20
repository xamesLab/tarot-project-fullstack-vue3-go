# Arcana UI Kit — Frontend Agent Prompt

> Build a standalone **UI Kit Reference page** as a Vue 3 single-file component (`UIKit.vue`) for a desktop tarot reading web service called **Arcana**.

---

## Tech Stack

- Vue 3 (Composition API, `<script setup>`)
- No external UI libraries (no Vuetify, no Element Plus)
- Google Fonts via `@import` or `<link>`: `Cormorant Garamond` (300, 400, 500, 600, italic variants) + `DM Sans` (300, 400, 500)
- Pure CSS (scoped), no Tailwind
- No icon libraries — use Unicode glyphs as symbols: `☽ ✦ ⊕ △ ◇ ♃`

---

## Design Tokens

Define all tokens as CSS custom properties on `:root`. The page must support **dark mode** (default) and **light mode**, toggled by adding class `.light-mode` to the root wrapper.

### Dark Mode (default)

```css
--bg: #0d0d12;
--surface: #16161f;
--surface2: #1e1e2e;
--surface3: #2a2a3e;
--border: rgba(200, 180, 255, 0.22);
--border2: rgba(200, 180, 255, 0.38);
--border-gold: rgba(201, 169, 110, 0.35);
--accent: #d4af72;        /* primary gold */
--accent2: #b08fff;       /* violet */
--accent3: #f0ddb8;       /* cream */
--text-primary: #f5f2ed;
--text-secondary: #c8c2d4;
--text-muted: #7a748a;
--card-bg: #13131e;
```

### Light Mode (`.light-mode` class on root wrapper)

```css
--bg: #f2ede4;
--surface: #faf6ee;
--surface2: #ede7da;
--surface3: #e0d8c8;
--border: rgba(80, 50, 140, 0.20);
--border2: rgba(80, 50, 140, 0.35);
--border-gold: rgba(120, 80, 20, 0.30);
--accent: #7a5210;
--accent2: #5530b8;
--accent3: #4a3010;
--text-primary: #14101c;
--text-secondary: #3a3248;
--text-muted: #7a6e8a;
--card-bg: #e8e2d4;
```

### Typography

```css
--font-display: 'Cormorant Garamond', Georgia, serif;
--font-body: 'DM Sans', system-ui, sans-serif;
```

### Spacing & Radius

```css
--radius-sm: 8px;
--radius-md: 12px;
--radius-card: 16px;
--radius-pill: 40px;
```

---

## Page Layout

The page is a full UI Kit reference divided into **named sections**, each with a section heading. Sections are stacked vertically with `gap: 48px`. Max content width: `1100px`, centered, padding `48px` horizontal.

At the top of the page: a **sticky header** with the Arcana logo (`☽ ARCANA`) on the left, and a **theme toggle button** (circular, `34×34px`) on the right that switches dark/light mode.

Below the header: a **sticky section navigation** — a horizontal pill-style tab bar listing all sections by name. Clicking a tab smoothly scrolls to that section.

---

## Section 1 — Brand Identity

Show the logo in two variants side by side on separate surface cards:

- **Dark variant**: `--bg: #0d0d12` card background, logo text in `--accent` gold
- **Light variant**: `--bg: #f2ede4` card background, logo text in `#7a5210`

Logo typographic spec:
- Font: `Cormorant Garamond`, 28px, weight 500, `letter-spacing: 0.12em`, `text-transform: uppercase`
- Symbol `☽` preceding the wordmark, rendered in `--accent2` violet, `font-size: 22px`

Below logos: a **Brand Voice** block — four keyword pills:
`Mystical` · `Precise` · `Trustworthy` · `Intimate`

Pills: `border: 1px solid --border2`, `border-radius: --radius-pill`, `padding: 6px 20px`, font `DM Sans` 11px, `letter-spacing: 0.12em`, `text-transform: uppercase`, color `--text-secondary`.

---

## Section 2 — Color Palette

Display two subsections: **Dark Scheme** and **Light Scheme**, each as a horizontal row of color swatches.

Each swatch:
- Rectangle `80×64px`, `border-radius: --radius-sm`
- Below: token name in 10px `DM Sans`, `--text-secondary`
- Below: hex value in 10px monospace, `--text-muted`

Dark scheme swatches (in order):

| Token | Hex | Label |
|---|---|---|
| `--bg` | `#0d0d12` | Background |
| `--surface` | `#16161f` | Surface |
| `--surface2` | `#1e1e2e` | Surface 2 |
| `--surface3` | `#2a2a3e` | Surface 3 |
| `--accent` | `#d4af72` | Gold |
| `--accent2` | `#b08fff` | Violet |
| `--accent3` | `#f0ddb8` | Cream |
| `--text-primary` | `#f5f2ed` | Text Primary |
| `--text-secondary` | `#c8c2d4` | Text Secondary |
| `--text-muted` | `#7a748a` | Text Muted |

Light scheme swatches use the light mode values for the same tokens.
For dark swatches `--bg`, `--surface`, `--surface2`, `--surface3`, add `border: 1px solid rgba(255,255,255,0.1)` so they're visible on the page background.

---

## Section 3 — Typography

A vertical stack of type specimens, each row showing: **label** (left, 10px uppercase `DM Sans`, `--text-muted`, fixed width `90px`) + **specimen text** (right).

| Role | Font | Size | Weight | Style | Letter-spacing | Color | Specimen text |
|---|---|---|---|---|---|---|---|
| Display | Cormorant Garamond | 52px | 300 | normal | 0.06em | `--text-primary` | `The World` |
| Heading 1 | Cormorant Garamond | 36px | 400 | normal | 0.05em | `--text-primary` | `Three Card Reading` |
| Heading 2 | Cormorant Garamond | 28px | 400 | normal | 0.04em | `--text-primary` | `The Narrative of Your Cards` |
| Subheading | Cormorant Garamond | 20px | 300 | italic | 0.02em | `--text-secondary` | `Past · Present · Future` |
| Body | DM Sans | 15px | 300 | normal | 0 | `--text-secondary` | `A cycle closes as wisdom crystallizes — the turbulence of uncertainty gives way to clarity.` |
| Caption | DM Sans | 11px | 500 | normal | 0.15em | `--text-muted` | `✦ READING ACTIVE` (uppercase) |
| Pullquote | Cormorant Garamond | 19px | 300 | italic | 0.01em | `--accent3` | `"What the Moon conceals, the Emperor reveals."` |

Separate each row with a `1px solid --border` horizontal rule.

Below the table: show the **font pairing** as a callout block with `border-left: 2px solid --accent`, padding `16px 24px`, listing:
- Display / Headings / Pullquotes → `Cormorant Garamond`
- Body / UI / Captions → `DM Sans`

---

## Section 4 — Iconography & Symbols

A grid of `4 columns` (responsive: 2 on narrow), each cell:
- Symbol glyph: `font-size: 32px`, color `--accent`, `font-family: var(--font-display)`
- Name below: `DM Sans` 10px uppercase `--text-muted`, `letter-spacing: 0.12em`
- Usage note below: `DM Sans` 11px `--text-secondary`

Symbols:

| Glyph | Name | Usage |
|---|---|---|
| `☽` | Moon | Logo, nav markers, lunar themes |
| `✦` | Star | CTA buttons, section accents |
| `⊕` | Circle Cross | Section dividers, card labels |
| `△` | Triangle | Fire element, upward motion |
| `◇` | Diamond | Air element, clarity |
| `♃` | Jupiter | Planet markers, expansion |
| `✧` | Sparkle | Hover states, decorative accents |
| `☿` | Mercury | Communication, swift themes |

Each cell: `background: --surface`, `border: 1px solid --border`, `border-radius: --radius-card`, `padding: 24px 16px`, `text-align: center`.

---

## Section 5 — Components: Buttons

Four button variants displayed in a row, each labeled below with its variant name.

**Primary button** (`btn-primary`):
- Background: `--accent` (`#d4af72`)
- Text color: `#0d0b08` (near-black)
- Border: none
- Font: `DM Sans` 11px, weight 500, `letter-spacing: 0.2em`, `text-transform: uppercase`
- Padding: `13px 36px`
- Border-radius: `--radius-pill`
- Hover: `opacity: 0.88`, `letter-spacing: 0.26em`
- Example label: `✦ Reveal Cards`

**Secondary button** (`btn-secondary`):
- Background: transparent
- Border: `1.5px solid --accent`
- Text: `--accent`
- Same font, padding, radius as primary
- Hover: background `rgba(212,175,114,0.08)`
- Example label: `☽ New Reading`

**Ghost button** (`btn-ghost`):
- Background: transparent
- Border: `1px solid --border2`
- Text: `--text-secondary`
- Same font, padding, radius
- Hover: background `--surface2`
- Example label: `Save to Journal`

**Accent violet button** (`btn-violet`):
- Background: `--accent2` (`#b08fff`)
- Text: `#fff`
- Border: none
- Same font, padding, radius
- Hover: `opacity: 0.88`
- Example label: `⊕ Celtic Cross`

Below all buttons: a **states reference row** showing the same button (Primary) in states: Default · Hover · Disabled (opacity 0.4, cursor not-allowed).

---

## Section 6 — Components: Form Controls

Show these controls stacked vertically with `gap: 16px`, each with a label above (11px DM Sans, `--text-secondary`, `letter-spacing: 0.08em`):

**Text Input:**
- Height: `44px`, width: `320px`
- Background: `--surface2`
- Border: `1px solid --border2`
- Border-radius: `--radius-md`
- Padding: `0 16px`
- Font: `DM Sans` 14px `--text-primary`
- Placeholder color: `--text-muted`
- Focus: `border-color: --accent`, `outline: none`, `box-shadow: 0 0 0 3px rgba(212,175,114,0.15)`
- Placeholder text: `Ask the cards…`

**Language Toggle** (segmented control):
- Two segments: `EN` · `RU`
- Outer container: `border: 1px solid --border2`, `border-radius: 20px`, `overflow: hidden`, `display: inline-flex`
- Each segment: `padding: 6px 14px`, `font-size: 11px`, `font-weight: 500`, `letter-spacing: 0.06em`
- Inactive: background transparent, color `--text-muted`
- Active: background `--surface3`, color `--accent`
- Transition: `background 0.2s`

**Theme Toggle** (icon button):
- Size: `34×34px`, `border-radius: 50%`
- Border: `1px solid --border2`
- Icon: `☀` (dark mode) or `☽` (light mode), `font-size: 15px`, color `--text-secondary`
- Hover: `border-color: --accent`, color `--accent`

**Status Badge:**
- Inline flex with a dot + text
- Dot: `5×5px`, `border-radius: 50%`, background `--accent`, animation `pulse` (opacity 1→0.4, 2s infinite)
- Text: `DM Sans` 10px, `letter-spacing: 0.08em`, `text-transform: uppercase`, color `--accent`, weight 500
- Example: `● PREMIUM`

---

## Section 7 — Components: Cards

Show three card types side by side:

**Surface Card** (general content container):
- Background: `--surface`
- Border: `1px solid --border2`
- Border-radius: `--radius-card`
- Padding: `28px`
- Contents: a short heading (`Cormorant Garamond` 20px `--text-primary`) + two lines of body text (`DM Sans` 14px `--text-secondary`)

**Reading Card** (interpretation block):
- Same as Surface Card
- Has a `border-left: 3px solid --accent` on the inner pullquote child
- Contains: title + pullquote text + tag pills row at bottom

**Tarot Card (face-down)**:
- Size: `140×240px`
- Background: `--card-bg` (`#13131e`)
- Border: `1px solid --border2`
- Border-radius: `12px`
- Back pattern (inline SVG centered inside):
  - Outer border rect: `8px inset`, `stroke: rgba(176,143,255,0.35)`, `stroke-width: 1`
  - Inner border rect: `16px inset`, `stroke: rgba(212,175,114,0.28)`, `stroke-width: 0.5`
  - Center mandala: concentric circles at r=36, r=24, r=12 with gold/violet strokes; two overlapping triangles (Star of David style) with `stroke: rgba(212,175,114,0.45)` and `rgba(176,143,255,0.35)`; axis lines; center dot `fill: rgba(212,175,114,0.7)`
  - Corner glyphs at four corners: `font-size: 10px`, `fill: rgba(212,175,114,0.55)`
- Show two states: **face-down** (back pattern) and **face-up** (white roman numeral top `DM Sans 11px --text-secondary`; large center glyph `42px Cormorant --accent`; card name bottom `13px Cormorant 500 --accent3`)

---

## Section 8 — Components: Navigation

**Top Navigation Bar** mockup:
- Background: `--surface`
- Border-bottom: `1px solid --border2`
- Height: `48px`
- Items: `Today` · `Three Card` (active) · `Celtic Cross` · `Year Ahead` · `Relationship` · `Journal`
- Each item: `DM Sans` 11px, weight 500, `letter-spacing: 0.12em`, `text-transform: uppercase`, padding `0 20px`
- Inactive color: `--text-muted`
- Hover color: `--text-secondary`
- Active: color `--accent`, `border-bottom: 2px solid --accent`
- Each item preceded by a small symbol glyph at `font-size: 9px`, `vertical-align: super`

---

## Section 9 — Dividers & Ornaments

Show four decorative divider styles stacked with labels:

1. **Simple line**: `1px solid --border2`, full width
2. **Gold gradient line**: `height: 1px`, `background: linear-gradient(to right, transparent, --border-gold, transparent)`
3. **Ornament divider**: flex row — `line (flex:1, 1px --border-gold)` + center text `✦ ☽ ✦` (10px, `letter-spacing: 0.3em`, `--accent`) + `line (flex:1)`
4. **Symbol divider**: centered row `⊕` glyph only, `font-size: 16px`, `--accent2`

---

## Section 10 — Motion & Animation Principles

A reference table with 3 columns: **Name** · **CSS spec** · **Live demo**.

| Name | Spec | Demo |
|---|---|---|
| Card flip | `rotateY 0deg→90deg→0deg, 0.6s ease` | A mini tarot card div that flips on hover |
| Hover lift | `translateY 0→-6px, scale 1→1.02, 0.3s ease` | A small card that lifts on hover |
| Text reveal | `opacity 0→1, translateY 12px→0, 0.5s ease` | A text block that fades in on button click |
| Status pulse | `opacity 1→0.4→1, 2s infinite` | The live status dot |
| Button press | `scale 1→0.97, 0.15s` | Active state on button click |

Live demos must actually work using Vue reactive state and CSS transitions/animations.

---

## Section 11 — Do's and Don'ts

Two columns: **Do** (left, `border-top: 3px solid #4a9e6f`) and **Don't** (right, `border-top: 3px solid #c0392b`).

Each column contains 4 rule cards (`background: --surface`, `border: 1px solid --border`, `border-radius: --radius-md`, `padding: 20px`):

**Do:**
1. Use `Cormorant Garamond` for all display and heading text
2. Use `--accent` gold sparingly as the primary highlight color
3. Keep symbol glyphs small and precise (`9–16px`) as accents only
4. Maintain `1px solid --border2` on all interactive containers

**Don't:**
1. Use system fonts (Arial, Inter, Roboto) anywhere in the UI
2. Use more than 2 accent colors simultaneously in one component
3. Use filled backgrounds on symbol/glyph elements — always transparent
4. Use `border-radius` larger than `--radius-card` (16px) except on pill buttons

---

## Vue Component Structure

```
UIKit.vue
├── <template>
│   ├── sticky header (logo + theme toggle)
│   ├── sticky section nav (tab pills)
│   └── main content
│       ├── Section 1: Brand Identity
│       ├── Section 2: Color Palette
│       ├── Section 3: Typography
│       ├── Section 4: Iconography
│       ├── Section 5: Buttons
│       ├── Section 6: Form Controls
│       ├── Section 7: Cards
│       ├── Section 8: Navigation
│       ├── Section 9: Dividers
│       ├── Section 10: Motion
│       └── Section 11: Do's & Don'ts
├── <script setup>
│   ├── const isDark = ref(true)
│   ├── const activeSection = ref('brand')
│   ├── function toggleTheme()
│   ├── function scrollTo(id)
│   └── motion demo reactive state
└── <style scoped>
    ├── :root tokens
    ├── .light-mode overrides
    └── all component styles
```

---

## Additional Requirements

- All sections must have an `id` attribute matching the nav tab for smooth-scroll targeting
- The sticky section nav must highlight the currently visible section using `IntersectionObserver`
- The theme toggle must apply `.light-mode` to the root wrapper `div` reactively via `:class`
- All transitions: `transition: background 0.4s, color 0.4s, border-color 0.3s` on root and interactive elements
- No hardcoded colors outside of CSS custom properties (except the SVG card pattern which uses rgba directly)
- The component must be self-contained and work when dropped into any Vue 3 project with no additional dependencies
