# Agent Rules

This file defines mandatory working rules for coding agents in this repository.

## UI Kit First (Frontend)

For any frontend task that includes:
- creating/updating Vue components,
- changing markup/layout,
- changing CSS/styles/theme/tokens,

you must follow the UI Kit as the source of truth:
- route: `/ui-kit` (dev-only)
- implementation reference: `frontend/src/pages/dev/UIKit.vue`
- policy doc: `frontend/UI_KIT_RULE.md`

### Mandatory behavior

1. Reuse existing UI Kit patterns, tokens, and variants first.
2. Do not introduce conflicting ad-hoc styles.
3. If a new visual pattern is needed, add it to UI Kit first (or in the same PR), then use it in feature code.
4. In task summaries, explicitly confirm UI Kit alignment for UI-related changes.

## Background Freeze Rule (Frontend)

Never modify background layout or background styles unless the user explicitly asks for it.

Protected scope includes:
- `frontend/src/components/AppBackground.vue`
- global background layers and decorative background effects in `frontend/src/assets/styles/base.css`
- background tokens in `frontend/src/assets/styles/variables.css`
- root background mounting/stacking in `frontend/src/App.vue`

If a task may require touching these files indirectly, stop and request explicit user confirmation first.

## Scope

- This rule is mandatory for all agents working in this repository.
