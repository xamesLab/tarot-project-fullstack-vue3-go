#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

if [[ ! -f ".env" ]]; then
  echo "warning: '.env' not found in backend folder; app will rely on process env vars"
fi

GO_BIN="go"
if ! command -v go >/dev/null 2>&1; then
  if [[ -x "/c/Program Files/Go/bin/go.exe" ]]; then
    GO_BIN="/c/Program Files/Go/bin/go.exe"
  else
    echo "go binary not found. Install Go or add it to PATH." >&2
    exit 1
  fi
fi

export GO111MODULE=on
"$GO_BIN" run ./cmd/api
