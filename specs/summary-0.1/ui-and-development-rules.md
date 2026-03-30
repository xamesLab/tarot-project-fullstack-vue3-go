# UI and Development Rules Summary

Last updated: March 30, 2026.

## UI Governance

Frontend development currently follows explicit design-governance rules.

### UI Kit First

For any frontend work that touches:

- Vue components
- markup or layout
- CSS, tokens, or styles

the source of truth is:

- route: `/ui-kit`
- implementation reference: `frontend/src/pages/dev/UIKit.vue`
- policy file: `frontend/UI_KIT_RULE.md`

### Background Freeze Rule

Background visuals must not be changed unless explicitly requested by the user.

Protected area includes:

- `frontend/src/components/AppBackground.vue`
- global background layers in base styles
- background tokens
- root background mounting behavior in `App.vue`

## Current Dev-Tools Reality

- `/ui-kit` exists today
- `/dev` is planned, but not implemented yet
- do not document `/dev` as available until the route and page actually exist in code

## Current Automation Policy

At the current moment:

- automated tests are postponed until explicit user instruction
- visual regression work is postponed until explicit user instruction
- automation is not the current active delivery track

## Current Storage Policy

At the current moment:

- external or object-based asset storage work is postponed until explicit user instruction
- current asset handling remains local and backend-served

## Working Recommendation

For the current project phase, prioritize:

- manual validation
- documentation quality
- integration clarity
- reproducible local runtime behavior

before:

- test automation expansion
- infrastructure migration
