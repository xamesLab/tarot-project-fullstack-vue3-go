# Arcana — Background Implementation Prompt

> Add a layered mystical background system to the Arcana tarot web service.
> The background must work in both **dark mode** (default) and **light mode**,
> respond to theme switching, and never interfere with UI readability.

---

## Overview

The background is a **purely decorative, non-interactive layer** rendered behind
all UI components. It is composed of:

1. A solid base color
2. Two or three radial color blooms (CSS or SVG gradients)
3. A centered geometric mandala (inline SVG, very low opacity)
4. Two constellation groups — lines + star dots (inline SVG)
5. Radial vignette darkening the edges
6. Static decorative glyphs at fixed positions

All layers are `position: fixed`, `z-index: 0`, `pointer-events: none`.
All UI content sits at `z-index: 1` or higher.

---

## File Structure

Create or modify these files:

```
src/
├── components/
│   └── AppBackground.vue     ← new component (self-contained)
├── App.vue                   ← import & mount AppBackground
└── assets/
    └── background.css        ← optional: CSS custom properties only
```

---

## AppBackground.vue — Full Specification

### Template Structure

```html
<template>
  <div class="arcana-bg" :class="{ 'light-mode': isLight }">

    <!-- Layer 1: base + blooms via CSS -->
    <div class="bg-blooms" aria-hidden="true"></div>

    <!-- Layer 2: SVG canvas (mandala + constellations + glyphs) -->
    <svg
      class="bg-svg"
      viewBox="0 0 1440 900"
      preserveAspectRatio="xMidYMid slice"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden="true"
      focusable="false"
    >
      <!-- defs: gradients for blooms and star glows -->
      <defs> ... </defs>

      <!-- bloom rects -->
      <rect class="bloom-top"    width="1440" height="900" />
      <rect class="bloom-br"     width="1440" height="900" />
      <rect class="bloom-bl"     width="1440" height="900" />

      <!-- mandala group -->
      <g class="mandala" transform="translate(720, 450)"> ... </g>

      <!-- constellation group A — top, gold -->
      <g class="constellation gold"> ... </g>

      <!-- constellation group B — bottom, violet -->
      <g class="constellation violet"> ... </g>

      <!-- decorative glyphs -->
      <g class="glyphs"> ... </g>

      <!-- vignette rect -->
      <rect class="vignette" width="1440" height="900" />
    </svg>

  </div>
</template>
```

### Props

```typescript
// AppBackground receives the current theme state as a prop
defineProps<{
  isLight: boolean   // false = dark mode (default), true = light mode
}>()
```

### Mounting in App.vue

```html
<!-- App.vue -->
<template>
  <AppBackground :is-light="isLight" />
  <div class="app-shell" :class="{ 'light-mode': isLight }">
    <!-- router-view / page content here -->
  </div>
</template>
```

```typescript
// App.vue script setup
import { ref } from 'vue'
import AppBackground from '@/components/AppBackground.vue'

const isLight = ref(false)
provide('isLight', isLight)  // so child components can call toggleTheme
```

---

## Layer 1 — CSS Blooms

The `.bg-blooms` div sits behind the SVG and provides the warm/cool
gradient atmosphere using pure CSS `radial-gradient`.

### Dark Mode CSS

```css
.arcana-bg {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
}

.bg-blooms {
  position: absolute;
  inset: 0;

  /* Base */
  background-color: #0d0d12;

  /* Bloom 1: deep violet, top-center */
  background-image:
    radial-gradient(
      ellipse 80% 60% at 50% -5%,
      rgba(61, 42, 128, 0.75) 0%,
      transparent 100%
    ),
    /* Bloom 2: amber, bottom-right */
    radial-gradient(
      ellipse 65% 55% at 95% 105%,
      rgba(107, 46, 8, 0.60) 0%,
      transparent 100%
    ),
    /* Bloom 3: deep indigo, left-mid */
    radial-gradient(
      ellipse 50% 45% at -5% 65%,
      rgba(30, 14, 85, 0.55) 0%,
      transparent 100%
    ),
    /* Bloom 4: subtle center pulse */
    radial-gradient(
      ellipse 40% 35% at 50% 50%,
      rgba(42, 18, 96, 0.30) 0%,
      transparent 100%
    );
}
```

### Light Mode CSS (add `.light-mode .bg-blooms` override)

```css
.arcana-bg.light-mode .bg-blooms {
  background-color: #f2ede4;

  background-image:
    /* Bloom 1: soft lavender, top-center */
    radial-gradient(
      ellipse 80% 60% at 50% -5%,
      rgba(196, 176, 240, 0.55) 0%,
      transparent 100%
    ),
    /* Bloom 2: warm gold, bottom-right */
    radial-gradient(
      ellipse 65% 55% at 95% 105%,
      rgba(212, 168, 74, 0.45) 0%,
      transparent 100%
    ),
    /* Bloom 3: muted violet, left-mid */
    radial-gradient(
      ellipse 50% 45% at -5% 65%,
      rgba(168, 152, 216, 0.38) 0%,
      transparent 100%
    );
}
```

---

## Layer 2A — SVG Gradient Defs

Place these inside `<defs>` in the SVG. They power the bloom rects
and the star glow halos.

```xml
<defs>

  <!-- Dark mode bloom gradients -->
  <radialGradient id="d-top" cx="50%" cy="0%" r="75%" gradientUnits="objectBoundingBox">
    <stop offset="0%"   stop-color="#3d2a80" stop-opacity="0.75"/>
    <stop offset="100%" stop-color="#0d0d12" stop-opacity="0"/>
  </radialGradient>

  <radialGradient id="d-br" cx="92%" cy="98%" r="65%" gradientUnits="objectBoundingBox">
    <stop offset="0%"   stop-color="#6b2e08" stop-opacity="0.60"/>
    <stop offset="100%" stop-color="#0d0d12" stop-opacity="0"/>
  </radialGradient>

  <radialGradient id="d-bl" cx="0%" cy="65%" r="50%" gradientUnits="objectBoundingBox">
    <stop offset="0%"   stop-color="#1e0e55" stop-opacity="0.55"/>
    <stop offset="100%" stop-color="#0d0d12" stop-opacity="0"/>
  </radialGradient>

  <!-- Dark mode vignette -->
  <radialGradient id="d-vig" cx="50%" cy="50%" r="72%" gradientUnits="objectBoundingBox">
    <stop offset="50%"  stop-color="#0d0d12" stop-opacity="0"/>
    <stop offset="100%" stop-color="#07070d" stop-opacity="0.72"/>
  </radialGradient>

  <!-- Light mode bloom gradients -->
  <radialGradient id="l-top" cx="50%" cy="0%" r="75%" gradientUnits="objectBoundingBox">
    <stop offset="0%"   stop-color="#c4b0f0" stop-opacity="0.55"/>
    <stop offset="100%" stop-color="#f2ede4" stop-opacity="0"/>
  </radialGradient>

  <radialGradient id="l-br" cx="92%" cy="98%" r="65%" gradientUnits="objectBoundingBox">
    <stop offset="0%"   stop-color="#d4a84a" stop-opacity="0.45"/>
    <stop offset="100%" stop-color="#f2ede4" stop-opacity="0"/>
  </radialGradient>

  <radialGradient id="l-bl" cx="0%" cy="65%" r="50%" gradientUnits="objectBoundingBox">
    <stop offset="0%"   stop-color="#a898d8" stop-opacity="0.38"/>
    <stop offset="100%" stop-color="#f2ede4" stop-opacity="0"/>
  </radialGradient>

  <!-- Light mode vignette -->
  <radialGradient id="l-vig" cx="50%" cy="50%" r="72%" gradientUnits="objectBoundingBox">
    <stop offset="55%"  stop-color="#f2ede4" stop-opacity="0"/>
    <stop offset="100%" stop-color="#ddd5c4" stop-opacity="0.60"/>
  </radialGradient>

  <!-- Star glow halos — dark -->
  <radialGradient id="glow-gold" cx="50%" cy="50%" r="50%">
    <stop offset="0%"   stop-color="#d4af72" stop-opacity="1"/>
    <stop offset="100%" stop-color="#d4af72" stop-opacity="0"/>
  </radialGradient>

  <radialGradient id="glow-violet" cx="50%" cy="50%" r="50%">
    <stop offset="0%"   stop-color="#b08fff" stop-opacity="1"/>
    <stop offset="100%" stop-color="#b08fff" stop-opacity="0"/>
  </radialGradient>

  <!-- Star glow halos — light -->
  <radialGradient id="glow-gold-l" cx="50%" cy="50%" r="50%">
    <stop offset="0%"   stop-color="#7a5210" stop-opacity="0.8"/>
    <stop offset="100%" stop-color="#7a5210" stop-opacity="0"/>
  </radialGradient>

  <radialGradient id="glow-violet-l" cx="50%" cy="50%" r="50%">
    <stop offset="0%"   stop-color="#5530b8" stop-opacity="0.7"/>
    <stop offset="100%" stop-color="#5530b8" stop-opacity="0"/>
  </radialGradient>

</defs>
```

---

## Layer 2B — SVG Bloom Rects

These rects use the gradient defs and are toggled by the theme class.

```xml
<!-- Dark mode blooms (hidden in light mode via CSS) -->
<rect class="bloom bloom-d" fill="url(#d-top)" width="1440" height="900"/>
<rect class="bloom bloom-d" fill="url(#d-br)"  width="1440" height="900"/>
<rect class="bloom bloom-d" fill="url(#d-bl)"  width="1440" height="900"/>

<!-- Light mode blooms (hidden in dark mode via CSS) -->
<rect class="bloom bloom-l" fill="url(#l-top)" width="1440" height="900"/>
<rect class="bloom bloom-l" fill="url(#l-br)"  width="1440" height="900"/>
<rect class="bloom bloom-l" fill="url(#l-bl)"  width="1440" height="900"/>
```

```css
/* show/hide bloom sets based on theme */
.arcana-bg .bloom-d { display: block; }
.arcana-bg .bloom-l { display: none; }

.arcana-bg.light-mode .bloom-d { display: none; }
.arcana-bg.light-mode .bloom-l { display: block; }
```

---

## Layer 2C — Mandala (SVG group)

Centered at `transform="translate(720, 450)"` on a 1440×900 viewBox.
All strokes use `fill="none"`. Toggle colors via CSS classes.

```xml
<g class="mandala" transform="translate(720, 450)" fill="none">

  <!-- Concentric circles: alternating gold / violet -->
  <circle r="185" class="m-violet" stroke-width="0.6"/>
  <circle r="148" class="m-gold"   stroke-width="0.8"/>
  <circle r="110" class="m-violet" stroke-width="0.6"/>
  <circle r="72"  class="m-gold"   stroke-width="0.9"/>
  <circle r="36"  class="m-violet" stroke-width="0.8"/>
  <circle r="12"  class="m-gold"   stroke-width="1.0"/>

  <!-- 8-fold axes: 0°, 90°, 45°, 135° -->
  <line x1="0"    y1="-185" x2="0"    y2="185"  class="m-gold-line"   stroke-width="0.4"/>
  <line x1="-185" y1="0"    x2="185"  y2="0"    class="m-gold-line"   stroke-width="0.4"/>
  <line x1="-131" y1="-131" x2="131"  y2="131"  class="m-violet-line" stroke-width="0.4"/>
  <line x1="131"  y1="-131" x2="-131" y2="131"  class="m-violet-line" stroke-width="0.4"/>

  <!-- Upward triangle (fire / gold) -->
  <polygon
    points="0,-90 77.9,45 -77.9,45"
    class="m-gold"
    stroke-width="1.0"
  />

  <!-- Downward triangle (water / violet) -->
  <polygon
    points="0,90 77.9,-45 -77.9,-45"
    class="m-violet"
    stroke-width="1.0"
  />

  <!-- Tick marks on outer ring at 8 cardinal/intercardinal points -->
  <g class="m-gold-line" stroke-width="0.6">
    <line x1="0"      y1="-148" x2="0"      y2="-138"/>
    <line x1="0"      y1="138"  x2="0"      y2="148"/>
    <line x1="-148"   y1="0"    x2="-138"   y2="0"/>
    <line x1="138"    y1="0"    x2="148"    y2="0"/>
    <line x1="-104.6" y1="-104.6" x2="-97.6" y2="-97.6"/>
    <line x1="97.6"   y1="97.6"   x2="104.6" y2="104.6"/>
    <line x1="104.6"  y1="-104.6" x2="97.6"  y2="-97.6"/>
    <line x1="-97.6"  y1="97.6"   x2="-104.6" y2="104.6"/>
  </g>

  <!-- Center dot -->
  <circle r="4" class="dot-gold"/>
  <circle r="2" class="dot-cream"/>

</g>
```

### Mandala CSS Classes

```css
/* --- Dark mode --- */
.arcana-bg .m-gold        { stroke: #d4af72; opacity: 0.26; }
.arcana-bg .m-violet      { stroke: #b08fff; opacity: 0.22; }
.arcana-bg .m-gold-line   { stroke: #d4af72; opacity: 0.16; }
.arcana-bg .m-violet-line { stroke: #b08fff; opacity: 0.16; }
.arcana-bg .dot-gold      { fill: #d4af72;   opacity: 0.70; }
.arcana-bg .dot-cream     { fill: #f0ddb8;   opacity: 0.90; }

/* --- Light mode overrides --- */
.arcana-bg.light-mode .m-gold        { stroke: #7a5210; opacity: 0.22; }
.arcana-bg.light-mode .m-violet      { stroke: #5530b8; opacity: 0.18; }
.arcana-bg.light-mode .m-gold-line   { stroke: #7a5210; opacity: 0.14; }
.arcana-bg.light-mode .m-violet-line { stroke: #5530b8; opacity: 0.14; }
.arcana-bg.light-mode .dot-gold      { fill: #7a5210;   opacity: 0.55; }
.arcana-bg.light-mode .dot-cream     { fill: #5530b8;   opacity: 0.70; }
```

---

## Layer 2D — Constellations

Two groups. Group A (top half) uses gold. Group B (bottom half) uses violet.
Each group has: connection lines, star circles, and glow halos on 2–3 key stars.

### Group A — Top, Gold (dark: #d4af72, light: #7a5210)

**Lines** (stroke-width 0.7):

| x1  | y1  | x2  | y2  |
|-----|-----|-----|-----|
| 55  | 110 | 175 | 72  |
| 175 | 72  | 295 | 145 |
| 295 | 145 | 415 | 92  |
| 415 | 92  | 478 | 195 |
| 175 | 72  | 238 | 192 |
| 238 | 192 | 295 | 145 |

**Star dots** (all gold fill):

| cx  | cy  | r   | opacity (dark) | opacity (light) |
|-----|-----|-----|----------------|-----------------|
| 55  | 110 | 2.2 | 0.75           | 0.60            |
| 175 | 72  | 3.0 | 0.85           | 0.70  ← key star |
| 295 | 145 | 2.4 | 0.75           | 0.60            |
| 415 | 92  | 2.8 | 0.80           | 0.65  ← key star |
| 478 | 195 | 2.0 | 0.65           | 0.50            |
| 238 | 192 | 1.8 | 0.65           | 0.52            |

**Glow halos** on key stars (ellipse, fill="url(#glow-gold)" or "#glow-gold-l"):

| cx  | cy | rx | ry | opacity (dark) | opacity (light) |
|-----|----|----|----|----------------|-----------------|
| 175 | 72 | 10 | 10 | 0.35           | 0.28            |
| 415 | 92 |  8 |  8 | 0.28           | 0.22            |

### Group B — Bottom, Violet (dark: #b08fff, light: #5530b8)

**Lines** (stroke-width 0.7):

| x1  | y1  | x2  | y2  |
|-----|-----|-----|-----|
| 72  | 710 | 195 | 748 |
| 195 | 748 | 348 | 688 |
| 348 | 688 | 458 | 755 |
| 195 | 748 | 258 | 828 |
| 348 | 688 | 420 | 620 |

**Star dots** (all violet fill):

| cx  | cy  | r   | opacity (dark) | opacity (light) |
|-----|-----|-----|----------------|-----------------|
| 72  | 710 | 2.2 | 0.70           | 0.55            |
| 195 | 748 | 3.0 | 0.82           | 0.65  ← key star |
| 348 | 688 | 2.0 | 0.68           | 0.55            |
| 458 | 755 | 2.4 | 0.72           | 0.58            |
| 258 | 828 | 1.8 | 0.62           | 0.50            |
| 420 | 620 | 2.0 | 0.65           | 0.52            |

**Glow halo** on key star:

| cx  | cy  | rx | ry | opacity (dark) | opacity (light) |
|-----|-----|----|----|----------------|-----------------|
| 195 | 748 | 10 | 10 | 0.32           | 0.26            |

### Scattered background stars

Spread across the full canvas. Alternate gold and violet. Use r=1.2–1.6, opacity 0.42–0.60 (dark), 0.35–0.50 (light).

Suggested positions (cx, cy) — feel free to adjust for natural distribution:

| cx  | cy  | color  |
|-----|-----|--------|
| 380 | 315 | gold   |
| 128 | 395 | gold   |
| 28  | 275 | gold   |
| 488 | 575 | gold   |
| 62  | 832 | gold   |
| 408 | 862 | gold   |
| 168 | 558 | violet |
| 318 | 498 | violet |
| 455 | 338 | violet |
| 92  | 520 | violet |
| 1380| 200 | gold   |
| 1200| 120 | violet |
| 1300| 680 | gold   |
| 1100| 780 | violet |
| 1050| 400 | gold   |
| 1380| 550 | violet |

### Constellation CSS

```css
/* Dark mode */
.arcana-bg .constellation.gold   { stroke: #d4af72; fill: none; opacity: 0.22; }
.arcana-bg .constellation.violet { stroke: #b08fff; fill: none; opacity: 0.20; }
.arcana-bg .star.gold            { fill: #d4af72; }
.arcana-bg .star.violet          { fill: #b08fff; }

/* Light mode */
.arcana-bg.light-mode .constellation.gold   { stroke: #7a5210; opacity: 0.20; }
.arcana-bg.light-mode .constellation.violet { stroke: #5530b8; opacity: 0.18; }
.arcana-bg.light-mode .star.gold            { fill: #7a5210; }
.arcana-bg.light-mode .star.violet          { fill: #5530b8; }
```

---

## Layer 2E — Decorative Glyphs

Fixed Unicode characters placed at precise positions in the SVG.
Font: `Georgia, serif`. All `text-anchor="middle"`.

| Glyph | x    | y   | font-size | class       | purpose         |
|-------|------|-----|-----------|-------------|-----------------|
| `☽`   | 720  | 52  | 40        | `glyph-crescent` | top center crescent |
| `✦`   | 28   | 42  | 16        | `glyph-star` | top-left corner |
| `✦`   | 1412 | 42  | 16        | `glyph-star` | top-right corner |
| `✦`   | 28   | 876 | 16        | `glyph-star` | bottom-left corner |
| `✦`   | 1412 | 876 | 16        | `glyph-star` | bottom-right corner |
| `⊕`   | 720  | 455 | 14        | `glyph-center` | center divider symbol |

Flanking lines for center divider (the `⊕`):

```xml
<line x1="80"  y1="450" x2="680" y2="450" class="divider-line"/>
<line x1="760" y1="450" x2="1360" y2="450" class="divider-line"/>
```

### Glyphs CSS

```css
/* Dark mode */
.arcana-bg .glyph-crescent { fill: #b08fff; opacity: 0.35; }
.arcana-bg .glyph-star     { fill: #d4af72; opacity: 0.38; }
.arcana-bg .glyph-center   { fill: #d4af72; opacity: 0.40; }
.arcana-bg .divider-line   { stroke: #d4af72; stroke-width: 0.6; opacity: 0.22; }

/* Light mode */
.arcana-bg.light-mode .glyph-crescent { fill: #5530b8; opacity: 0.28; }
.arcana-bg.light-mode .glyph-star     { fill: #7a5210; opacity: 0.32; }
.arcana-bg.light-mode .glyph-center   { fill: #7a5210; opacity: 0.35; }
.arcana-bg.light-mode .divider-line   { stroke: #7a5210; opacity: 0.20; }
```

---

## Layer 2F — Inner Frame Border

A subtle decorative rect just inside the SVG edge.

```xml
<rect
  x="18" y="18"
  width="1404" height="864"
  rx="4"
  fill="none"
  class="frame-border"
/>
```

```css
.arcana-bg .frame-border { stroke: #b08fff; stroke-width: 0.5; opacity: 0.10; }
.arcana-bg.light-mode .frame-border { stroke: #5530b8; opacity: 0.08; }
```

---

## Layer 2G — Vignette

The final SVG rect that darkens/lightens the edges, creating depth.

```xml
<!-- Dark mode vignette (hidden in light mode) -->
<rect class="vignette vignette-d" fill="url(#d-vig)" width="1440" height="900"/>

<!-- Light mode vignette (hidden in dark mode) -->
<rect class="vignette vignette-l" fill="url(#l-vig)" width="1440" height="900"/>
```

```css
.arcana-bg .vignette-d { display: block; }
.arcana-bg .vignette-l { display: none; }

.arcana-bg.light-mode .vignette-d { display: none; }
.arcana-bg.light-mode .vignette-l { display: block; }
```

---

## Transition on Theme Switch

All color/opacity changes must animate smoothly when toggling dark/light.

```css
.arcana-bg,
.arcana-bg .bg-blooms,
.arcana-bg .bg-svg {
  transition:
    background-color 0.5s ease,
    opacity          0.4s ease;
}

/* SVG elements that change color on theme switch */
.arcana-bg .mandala,
.arcana-bg .constellation,
.arcana-bg .glyph-crescent,
.arcana-bg .glyph-star,
.arcana-bg .glyph-center,
.arcana-bg .divider-line,
.arcana-bg .frame-border {
  transition: opacity 0.4s ease, stroke 0.4s ease, fill 0.4s ease;
}
```

---

## Z-Index Stack

```
z-index 0  → .arcana-bg (the entire background component)
z-index 1  → .app-shell (all page content: header, nav, main, cards)
z-index 50 → sticky header
z-index 100→ modals / dropdowns
```

---

## Performance Notes

- The SVG uses `preserveAspectRatio="xMidYMid slice"` — it scales to fill the
  viewport at any resolution without repainting.
- `position: fixed` ensures the background stays put during scroll.
- No JavaScript animations — all effects are static SVG + CSS.
- No external images, no canvas, no WebGL. Pure SVG + CSS only.
- Total weight: ~4 KB unminified. No performance impact.

---

## Readability Contract

The background must NEVER compromise UI legibility. Enforce these limits:

| Element | Max opacity allowed |
|---|---|
| Mandala strokes | 0.32 |
| Constellation lines | 0.25 |
| Star dots | 0.85 (at r ≥ 2px only) |
| Glow halos | 0.38 |
| Bloom radial gradients | 0.75 |
| Vignette | 0.72 |
| Decorative glyphs ☽ ✦ ⊕ | 0.40 |

If at any point body text (`--text-secondary: #c8c2d4`) on dark or
(`--text-secondary: #3a3248`) on light falls below **WCAG AA 4.5:1** contrast
against the background, reduce the nearest bloom opacity by 0.10 increments
until it passes.

---

## Checklist for the Agent

- [ ] `AppBackground.vue` created in `src/components/`
- [ ] Imported and mounted as first child in `App.vue` template
- [ ] `isLight` prop wired correctly from parent theme state
- [ ] Both dark and light bloom rects toggle via CSS class (no JS color changes)
- [ ] Mandala renders at center of viewport (translate 720,450)
- [ ] Both constellation groups present with correct coordinates
- [ ] Glow halos on key stars (cx:175/cy:72, cx:415/cy:92, cx:195/cy:748)
- [ ] All 6 decorative glyphs placed at correct positions
- [ ] Inner frame border rect present
- [ ] Vignette rects toggle correctly with theme
- [ ] `pointer-events: none` on the entire component
- [ ] `aria-hidden="true"` on the SVG
- [ ] Smooth theme transition (0.4–0.5s ease) on all changing properties
- [ ] z-index does not block any UI interaction
- [ ] Tested at viewport widths: 1280px, 1440px, 1920px
- [ ] No console errors, no layout shifts
