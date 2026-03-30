# UI Kit Source of Truth

Last updated: March 30, 2026.

## Purpose

The current UI Kit is a standalone development reference page for the Arcana visual language. It is not just a moodboard. It is the canonical visual contract for future UI work.

## Route and Runtime Model

Route:

- `/ui-kit`

Registration model:

- registered only in dev mode
- guarded again at router level through `meta.devOnly`
- redirected to `not-found` outside dev mode

This means the page is intentionally an internal design reference, not a production feature page.

## Page Structure

The UI Kit page is a single scrolling reference document with sticky controls and 11 sections.

Top structure:

- root shell
- sticky header
- sticky section navigation
- stacked content sections

Sections:

1. Brand Identity
2. Color Palette
3. Typography
4. Iconography and Symbols
5. Components: Buttons
6. Components: Form Controls
7. Components: Cards
8. Components: Navigation
9. Dividers and Ornaments
10. Motion and Animation Principles
11. Do and Don't rules

## Header

The header contains:

- Arcana wordmark
- lunar symbol
- circular theme toggle

Behavior:

- sticky near the top of viewport
- rounded pill container
- bordered and elevated
- compact spacing, never oversized

## Section Navigation

The navigation row contains one tab per section.

Behavior:

- sticky under the header
- horizontally scrollable
- active tab updates from section visibility
- click scrolls smoothly to target section

The active section is derived with `IntersectionObserver`, not manual scroll math.

## Visual Philosophy

The UI Kit aims for:

- mystical, intimate, precise presentation
- decorative restraint instead of overloaded fantasy styling
- serif display typography with clean sans-serif body text
- light and dark parity
- small symbolic accents instead of large icon-heavy chrome

## Non-Negotiable Rules

- reuse existing patterns before inventing new ones
- reuse tokens before introducing new ones
- keep glyphs as accents, not as dominant filled graphics
- keep interactive containers on thin borders with soft elevation
- preserve uppercase micro-label rhythm for UI metadata
