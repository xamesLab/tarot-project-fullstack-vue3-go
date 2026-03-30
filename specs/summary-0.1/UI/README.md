# UI Transfer Pack

Last updated: March 30, 2026.

This folder is a focused handoff pack for a coding agent that needs to recreate the Tarot Daily UI Kit in another project as accurately as possible.

## Source of Truth

Primary implementation source:

- `frontend/src/pages/dev/UIKit.vue`

Supporting policy and alignment files:

- `frontend/UI_KIT_RULE.md`
- `frontend/src/components/ModeSwitcher.vue`
- `frontend/src/components/ModeToggle.vue`
- `frontend/src/router/dev-tools.routes.ts`
- `frontend/src/router/guards.ts`

## Recommended Reading Order

1. `ui-kit-source-of-truth.md`
2. `visual-tokens-and-typography.md`
3. `component-specification.md`
4. `layout-responsive-and-interaction.md`
5. `agent-transfer-playbook.md`

## Important Scope Boundary

This pack is about the UI Kit itself, not the entire app shell.

Do not treat these files as part of the transferable UI Kit baseline:

- `frontend/src/components/AppBackground.vue`
- `frontend/src/assets/styles/base.css` background layers
- `frontend/src/assets/styles/variables.css` global app theme tokens
- `frontend/src/App.vue` background mounting

Reason:

- the Arcana background system is a separate concern
- the UI Kit should be transferable even without the background implementation
- the repository has an explicit `Background Freeze Rule`

## Accuracy Notes

- `/ui-kit` exists today and is the real visual source of truth
- `/dev` does not exist yet and should not be assumed
- the UI Kit page uses its own scoped token set and Google Fonts import
- the app-wide token file is related, but not identical, to the UI Kit page
