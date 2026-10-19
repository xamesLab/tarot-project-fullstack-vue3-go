# Development mocks

This directory contains the frontend mock API and its fixture data. `start-dev.sh`
uses it automatically when MongoDB is not configured or cannot be reached, and
it can be selected directly with `DEV_DATA_MODE=mock`.

Mock mode runs the complete frontend against local mock adapters. It does not
start the Go API. Mock changes live in browser memory; the guest daily-card
bootstrap remains cached in local storage for the current accounting date.

Use `DEV_DATA_MODE=mongo` to require the Go API and MongoDB, or leave the mode
unset to try MongoDB first and fall back to mocks when its connection fails.
