# Visual Tokens and Typography

Last updated: March 30, 2026.

## Fonts

The UI Kit page imports and uses:

- `Cormorant Garamond` for display, headings, pullquotes, logo
- `DM Sans` for body text, labels, controls, metadata

Important transfer note:

- these are the real UI Kit fonts
- the broader app still has older global font fallbacks in `variables.css`
- if the goal is accurate transfer of the UI Kit, prefer the UI Kit fonts, not the older app-wide font fallbacks

## Brand Voice Keywords

- Mystical
- Precise
- Trustworthy
- Intimate

## Dark Theme Tokens

- `--dark-bg`: `#0d0d12`
- `--dark-surface`: `#16161f`
- `--dark-surface2`: `#1e1e2e`
- `--dark-surface3`: `#2a2a3e`
- `--dark-accent`: `#d4af72`
- `--dark-accent2`: `#b08fff`
- `--dark-accent3`: `#f0ddb8`
- `--dark-text-primary`: `#f5f2ed`
- `--dark-text-secondary`: `#c8c2d4`
- `--dark-text-muted`: `#7a748a`

Supporting dark borders:

- `--dark-border`: `rgba(200, 180, 255, 0.22)`
- `--dark-border2`: `rgba(200, 180, 255, 0.38)`
- `--dark-border-gold`: `rgba(201, 169, 110, 0.35)`

## Light Theme Tokens

- `--light-bg`: `#f2ede4`
- `--light-surface`: `#faf6ee`
- `--light-surface2`: `#ede7da`
- `--light-surface3`: `#e0d8c8`
- `--light-accent`: `#7a5210`
- `--light-accent2`: `#5530b8`
- `--light-accent3`: `#4a3010`
- `--light-text-primary`: `#14101c`
- `--light-text-secondary`: `#3a3248`
- `--light-text-muted`: `#7a6e8a`

Supporting light borders:

- `--light-border`: `rgba(80, 50, 140, 0.2)`
- `--light-border2`: `rgba(80, 50, 140, 0.35)`
- `--light-border-gold`: `rgba(120, 80, 20, 0.3)`

## Shared Semantic Tokens

These are the practical transfer-level tokens the page actually consumes:

- `--bg`
- `--surface`
- `--surface2`
- `--surface3`
- `--border`
- `--border2`
- `--border-gold`
- `--accent`
- `--accent2`
- `--accent3`
- `--text-primary`
- `--text-secondary`
- `--text-muted`
- `--card-bg`
- `--font-display`
- `--font-body`
- `--radius-sm`: `8px`
- `--radius-md`: `12px`
- `--radius-card`: `16px`
- `--radius-pill`: `40px`

Additional constants:

- `--ink-strong`: `#0d0b08`
- `--white`: `#ffffff`
- `--do-color`: `#4a9e6f`
- `--dont-color`: `#c0392b`
- `--shadow-card`: `0 20px 34px rgba(0, 0, 0, 0.24)` in dark mode
- `--shadow-lift`: `0 16px 26px rgba(0, 0, 0, 0.2)` in dark mode
- `--focus-ring`: `rgba(212, 175, 114, 0.15)`

## Typography Hierarchy

Display:

- `Cormorant Garamond`
- very large
- letter spacing around `0.06em`
- light or regular weight

Heading 1:

- `Cormorant Garamond`
- around `28-36px`
- soft tracking

Heading 2:

- `Cormorant Garamond`
- around `23-28px`

Subheading:

- `Cormorant Garamond`
- italic
- secondary text color

Body:

- `DM Sans`
- around `14-15px`
- lightweight, readable, not dense

Caption and metadata:

- `DM Sans`
- `10-11px`
- uppercase
- extended letter spacing

Pullquote:

- `Cormorant Garamond`
- italic
- accent3 color

## Accent Policy

Use accents in this order of importance:

1. gold accent as primary highlight
2. violet accent as secondary highlight
3. pale accent for refined text emphasis

Avoid:

- using more than two accents prominently inside the same component
- filling icon blocks with strong accent backgrounds
- using neon contrast or glossy gradients that break the restrained mystical tone
