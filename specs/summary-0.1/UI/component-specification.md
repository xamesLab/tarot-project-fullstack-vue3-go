# Component Specification

Last updated: March 30, 2026.

## Brand Block

Composition:

- lunar symbol
- uppercase wordmark `ARCANA`

Visual rules:

- display font
- symbol in secondary accent
- wordmark in primary accent
- tracking around `0.12em`

## Header Theme Toggle

Shape:

- circle
- `34x34`

Style:

- transparent background
- thin border
- muted text by default
- accent color on hover

Content:

- sun when current mode is dark
- moon when current mode is light

## Section Tabs

Style:

- pill shape
- thin border
- uppercase micro-label style
- transparent by default

Active state:

- accent text
- accent border
- elevated surface fill

## Buttons

Variants in the current UI Kit:

- Primary
- Secondary
- Ghost
- Violet Accent

### Primary

- filled with primary accent
- dark ink text
- uppercase label
- expanded tracking
- hover slightly reduces opacity and increases tracking

### Secondary

- transparent fill
- accent border
- accent text
- subtle focus-ring-like hover fill

### Ghost

- transparent fill
- muted border
- secondary text color

### Violet Accent

- violet fill
- white text

Shared rules:

- pill radius
- minimum height around `44px`
- label size around `11px`
- uppercase metadata feel

## Inputs

Text input rules:

- width capped around `320px`
- height `44px`
- medium rounded corners
- `surface2` fill
- thin strong border
- accent focus ring

Placeholder:

- muted text color

## Language Toggle

Structure:

- inline segmented control
- two segments
- `EN` and `RU` in the UI Kit demo

Style:

- outer container border
- overflow hidden
- segment padding around `6px 14px`
- active segment uses elevated surface and accent text

Important real-app alignment:

- `ModeToggle.vue` implements the same visual idea with actual locale values `ru` and `en`

## Theme Toggle

The UI Kit control and the real app `ModeToggle.vue` are aligned on:

- `34x34` circular size
- transparent fill
- thin strong border
- accent hover
- moon or sun symbol

## Status Badge

Composition:

- small pulsing dot
- uppercase micro-label text

Style:

- accent color only
- very small scale
- never large or noisy

## Cards

Three reference card surfaces exist:

1. neutral surface card
2. reading card with pullquote and tag pills
3. tarot card showcase with back and front states

### Surface Card

- neutral content container
- `surface2` background
- medium rounding
- thin border

### Reading Quote Card

- pullquote with left accent border
- serif italic quote
- tag pills for positions like Past, Present, Future

### Tarot Card Showcase

Back card:

- `140x240`
- decorative SVG geometry
- border frame and concentric circular motif

Front card:

- roman numeral
- large central glyph
- serif card title

Responsive sizes:

- `126x216` below `760px`
- `112x192` below `480px`

## Top Navigation Demo

Structure:

- horizontal row
- fixed height around `48px`
- scrolling overflow allowed

State model:

- active tab uses accent color
- active tab gets bottom border
- glyphs remain tiny

## Dividers

Supported divider types:

- simple line
- gold gradient line
- ornament divider with symbol cluster
- standalone symbol divider

These should be reused as low-cost decorative separators across the design system.

## Motion Patterns

The UI Kit demonstrates these standard motion primitives:

- card flip: `0.6s ease`
- hover lift: `0.3s ease`
- text reveal: `0.5s ease`
- status pulse: `2s infinite`
- button press: `0.15s`

Transfer rule:

- prefer a few meaningful motions
- do not add generic micro-animations everywhere

## Do and Don't Rules from Source

Do:

- use Cormorant Garamond for all display and heading text
- use gold accent sparingly as the primary highlight color
- keep symbol glyphs small and precise, around `9-16px`, as accents only
- maintain thin `1px` borders on interactive containers

Don't:

- use system fonts like Arial, Inter, or Roboto as the primary UI Kit expression
- use more than two accent colors simultaneously in one component
- use filled backgrounds on symbol or glyph elements
- use radius larger than the card radius except on pill buttons
