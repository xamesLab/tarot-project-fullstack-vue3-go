# Layout, Responsive, and Interaction Rules

Last updated: March 30, 2026.

## Page Shell

Main shell:

- max width around `1100px`
- centered
- generous but controlled horizontal padding
- grid layout with vertical gaps

Default shell padding uses `clamp()` and scales between compact mobile and wide desktop layouts.

## Sticky Structure

Sticky elements:

- header
- section navigation

Desktop offsets:

- header top around `8px`
- section nav top around `74px`

Mobile offsets:

- section nav top around `68px` below `760px`
- section nav top around `60px` below `480px`

## Section Containers

Each section is:

- a bordered card-like surface
- rounded with `--radius-card`
- padded with `16-26px` range
- elevated with soft shadow

This means the UI Kit is designed as a stack of consistent surface cards, not as flat document sections.

## Grid Behavior

Desktop:

- brand and rules use two columns
- icon grid uses four columns
- button grid uses four columns
- cards grid uses three columns

Mid-width:

- cards grid becomes two columns below `1100px`
- button and icon grids become two columns below `900px`

Mobile:

- brand, rules, cards, and button grids collapse to one column below `760px`
- icon grid stays two columns below `760px`
- icon grid becomes one column below `480px`

## Responsive Breakpoints

Current breakpoints in source:

- `1100px`
- `900px`
- `760px`
- `480px`

These are meaningful design breakpoints and should be preserved during transfer.

## Typography Scaling

The UI Kit uses `clamp()` for major display sizes.

Transfer expectation:

- keep large text fluid
- keep micro-labels compact
- do not upscale labels or metadata on mobile

## Interaction Model

### Theme switching

- instant state swap
- colors transition smoothly
- no dramatic page animation

### Section navigation

- click causes smooth scroll
- active state tracks viewport visibility

### Hover model

- subtle lift
- slight letter-spacing shift on strong CTA
- no aggressive scale explosions

### Press model

- tiny compression to about `scale(0.97)`

### Reveal model

- opacity from `0` to `1`
- translateY from `12px` to `0`

## Motion Discipline

The design language is calm and ritual-like.

That means:

- use motion to suggest intention and atmosphere
- keep durations short to medium
- avoid bouncy, playful, gamified easing
- avoid stacking multiple motion effects on the same small control

## Transfer Boundaries

Do transfer:

- sticky header and section nav pattern
- section-card layout model
- grids, spacing, and breakpoints
- motion primitives

Do not blindly transfer:

- dev-only routing behavior if the new project is not dev-only
- repo-specific background mounting
- app-specific fixed mode-switcher placement unless the new shell needs it
