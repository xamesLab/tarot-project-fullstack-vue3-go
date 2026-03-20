# UI Kit First Rule

This rule is mandatory for all frontend work related to components, layout, and styles.

## Rule

When creating components or changing markup/styles, always align with the UI Kit page:
- Route: `/ui-kit` (dev-only)
- Source: `frontend/src/pages/dev/UIKit.vue`

## Mandatory Requirements

1. Reuse existing visual patterns from UI Kit first.
2. Reuse existing tokens/variables (color, spacing, radius, typography) before introducing new ones.
3. Reuse existing component variants (buttons, form controls, cards, navigation patterns) where possible.
4. Do not introduce ad-hoc styles that conflict with UI Kit language.
5. If a new visual pattern is required, first add it to UI Kit, then use it in feature code.

## Pull Request Checklist (UI)

- [ ] Compared the change with `/ui-kit`.
- [ ] Reused existing tokens/components/variants where applicable.
- [ ] Any new pattern was added to UI Kit first (or in the same PR).
- [ ] No conflicting typography/color/radius conventions introduced.

## Exceptions

Exception is allowed only when:
1. There is a product requirement not covered by UI Kit.
2. The exception is documented in PR description.
3. A follow-up task is created to add/normalize the pattern in UI Kit.
