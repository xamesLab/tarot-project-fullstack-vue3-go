# Agent Transfer Playbook

Last updated: March 30, 2026.

This document is written for a coding agent that must recreate the UI Kit in a new project with high fidelity.

## Goal

Rebuild the Arcana UI Kit as a reusable design-system reference, not as a one-off static mockup.

## Minimum Transfer Outcome

The new project should contain:

- a standalone UI Kit page or route
- the same visual token system
- the same typography hierarchy
- the same core components and variants
- the same 11 reference sections
- the same light and dark parity
- the same responsive breakpoints

## Recommended Implementation Order

### 1. Recreate the token layer first

Implement:

- dark and light palettes
- semantic aliases like `--bg`, `--surface`, `--accent`, `--text-primary`
- radius, shadow, and font tokens

Do not start with feature pages before this exists.

### 2. Install correct fonts

Prefer:

- `Cormorant Garamond`
- `DM Sans`

If exact font loading is impossible at first:

- still keep the token names and document the fallback
- do not silently swap the visual identity to generic system fonts

### 3. Recreate the page shell

Build:

- root container
- sticky header
- sticky section-nav row
- section stack card layout

Preserve:

- max width
- rounded surfaces
- thin borders
- soft elevation

### 4. Recreate sections in the original order

Order matters because the page teaches the system progressively:

1. Brand
2. Palette
3. Typography
4. Iconography
5. Buttons
6. Forms
7. Cards
8. Navigation
9. Dividers
10. Motion
11. Do and Don't

### 5. Recreate component variants exactly before inventing abstractions

Build concrete visual parity first for:

- button variants
- segmented language toggle
- circular theme toggle
- status badge
- surface cards
- tarot back and front cards
- top navigation row
- divider family

Only then generalize into reusable components.

### 6. Recreate interaction logic

Implement:

- theme toggle
- smooth scroll to section
- active section tracking via viewport observation
- reveal and press demos

For fidelity, use intersection-based active state logic, not just click state.

### 7. Apply responsive behavior last

After desktop parity is correct, match:

- `1100px`
- `900px`
- `760px`
- `480px`

Do not improvise alternate breakpoints unless the new project has a hard technical reason.

## Fidelity Rules

When there is a conflict between:

- generic app design conventions
- the Arcana UI Kit source

prefer the Arcana UI Kit source.

In particular:

- prefer serif display + sans body pairing
- prefer thin borders over thick outlined fantasy panels
- prefer restrained symbolism over decorative overload
- prefer soft surfaces over flat white cards or pure black slabs

## Common Transfer Mistakes to Avoid

- flattening the layout into ordinary SaaS cards
- replacing the typography with Inter-like defaults
- turning glyph accents into large filled icon buttons
- losing the sticky dual-layer header and nav behavior
- changing the accent hierarchy so violet becomes dominant everywhere
- porting the background system together with the UI Kit even when it is not needed

## What Can Be Adapted

These parts may be adapted carefully if the target project requires it:

- route path
- framework syntax
- state-management implementation
- exact font loading mechanism
- icon implementation method

These parts should remain visually stable:

- token values
- typography hierarchy
- component proportions
- breakpoints
- motion intent
- section order and examples

## Practical Acceptance Checklist for the New Project

- the new UI Kit page feels recognizably Arcana at first glance
- light and dark mode both work
- section nav is sticky and highlights the active section
- buttons, toggles, cards, dividers, and motion demos match the source behavior
- mobile layout still feels elegant below `760px` and `480px`
- no background-only dependencies are required for the UI Kit to function
