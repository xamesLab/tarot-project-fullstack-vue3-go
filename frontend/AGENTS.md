# Frontend Agent Rules

This file is mandatory for any coding agent working in `frontend/`.

## UI Kit First

Before implementing any UI change, use UI Kit as the primary reference:
- route: `/ui-kit` (dev-only)
- source: `src/pages/dev/UIKit.vue`
- detailed policy: `UI_KIT_RULE.md`

Applies to:
- new components,
- markup/layout changes,
- style/theme/token changes,
- refactors affecting visual behavior.

### Required workflow

1. Check existing UI Kit patterns first.
2. Reuse existing tokens/components/variants.
3. Avoid introducing conflicting visual language.
4. If you need a new pattern, add/update it in UI Kit first (or in the same PR), then apply in product screens.
5. Mention UI Kit alignment in your final change summary.

## Background Freeze Rule

Never touch background markup/styles unless the user explicitly requests background changes.

Protected files:
- `src/components/AppBackground.vue`
- `src/assets/styles/base.css` (background/decorative layers)
- `src/assets/styles/variables.css` (background tokens)
- `src/App.vue` (background mounting and z-index layering)

If the task intersects these files, pause and ask for explicit permission before editing.
