# Development mocks

This directory contains the frontend mock API and its fixture data. `start-dev.sh`
uses it automatically when MongoDB is not configured or cannot be reached, and
it can be selected directly with `DEV_DATA_MODE=mock`.

The mock daily card uses the provided face and card-back PNG fixtures in
`frontend/assets`. The card container follows the face image's proportions. The
back preserves its proportions, fills the card outline, and is cropped from the
center when its aspect ratio differs from the card.

Mock mode runs the complete frontend against local mock adapters. It does not
start the Go API. Mock changes live in browser memory; the guest daily-card
bootstrap remains cached in local storage for the current accounting date.

Use `DEV_DATA_MODE=mongo` to require the Go API and MongoDB, or leave the mode
unset to try MongoDB first and fall back to mocks when its connection fails.
