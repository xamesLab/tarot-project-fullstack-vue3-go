#!/usr/bin/env bash

set -euo pipefail

export NVM_DIR="${NVM_DIR:-$HOME/.nvm}"

if [ ! -s "$NVM_DIR/nvm.sh" ]; then
  echo "nvm was not found at $NVM_DIR" >&2
  exit 1
fi

# Load nvm explicitly because Codex commands run in a non-interactive shell.
. "$NVM_DIR/nvm.sh"

if [ $# -eq 0 ]; then
  echo "Usage: scripts/with-node.sh <command> [args...]" >&2
  exit 1
fi

exec "$@"
