# AI Agent Development Plan
## Phase 6: Deployment

## Mandatory Documentation Sync

Before performing deployment work, the agent must review:

1. `specs/docs/PRD.md`
2. `specs/docs/TECHNICAL_SPEC.md`
3. `specs/docs/IMPLEMENTATION_ROADMAP.md`
4. `specs/v0.1.0/phase-6-deployment/README.md`
5. `specs/v0.1.0/phase-6-deployment/MANUAL_TESTS.md`

The agent must confirm the intended production topology, public URLs, environment variables, and smoke-test scope before deployment actions begin.

## Mission

Deploy the MVP safely to production, verify core functionality, and leave the project with enough operational visibility for the first live users.

## Step-by-Step Execution

### Step 1. Confirm release readiness

- verify all Phase 5 exit criteria are actually met
- ensure known critical bugs are resolved
- ensure required documentation exists and reflects reality

### Step 2. Confirm production configuration

- verify backend env vars
- verify frontend env vars
- verify database connection details
- verify JWT secret and other secrets are production-safe
- verify seed and migration assumptions

### Step 3. Build production artifacts

- produce backend build output
- produce frontend production bundle
- verify build success without relying on development-only behavior

### Step 4. Prepare production infrastructure

- deploy backend binary or service package
- deploy frontend static files
- configure reverse proxy
- configure `/uploads` static serving
- enable HTTPS

### Step 5. Validate service startup

- confirm backend starts successfully in production mode
- confirm frontend is publicly reachable
- confirm database connectivity from the production environment

### Step 6. Execute smoke tests

- health or stub endpoint
- user registration
- user login
- card-of-day generation
- spread generation
- image delivery
- theme switch
- language switch
- admin access flow

### Step 7. Set up monitoring and observability

- verify logs are accessible
- verify basic uptime checks are in place
- document where to look during incidents

### Step 8. Perform launch watch

- monitor the first 24 to 48 hours
- watch for backend errors, auth failures, upload failures, and broken assets
- triage issues quickly and document hotfixes

### Step 9. Close the phase with an operational summary

- what was deployed
- what was verified
- what remains risky
- what follow-up items exist for post-MVP work

## AI Agent Guardrails

- do not deploy if smoke tests already fail in staging or local equivalent
- do not expose the app publicly without HTTPS
- do not treat missing monitoring as acceptable completion
- do not skip smoke tests after deployment

## Completion Standard

This phase is complete only when the production environment is live, smoke-tested, observable, and stable enough for initial real-world usage.

