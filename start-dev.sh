#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"
RUNTIME_DIR="$ROOT_DIR/.dev-runtime"
BACKEND_PID_FILE="$RUNTIME_DIR/backend.pid"
FRONTEND_PID_FILE="$RUNTIME_DIR/frontend.pid"

backend_pid=""
frontend_pid=""

prepend_to_path() {
  local bin_dir="$1"

  case ":$PATH:" in
    *":$bin_dir:"*) ;;
    *) PATH="$bin_dir:$PATH" ;;
  esac
}

find_nvm_node_bin_dir() {
  local nvm_dir="${NVM_DIR:-$HOME/.nvm}"
  local npm_path=""

  if [[ -s "$nvm_dir/nvm.sh" ]]; then
    # shellcheck disable=SC1090
    . "$nvm_dir/nvm.sh" >/dev/null 2>&1 || true

    if command -v nvm >/dev/null 2>&1; then
      if nvm use --silent default >/dev/null 2>&1 || nvm use --silent node >/dev/null 2>&1; then
        npm_path="$(command -v npm || true)"
        if [[ -n "$npm_path" ]]; then
          dirname "$npm_path"
          return 0
        fi
      fi
    fi
  fi

  if compgen -G "$nvm_dir/versions/node/*/bin/npm" >/dev/null 2>&1; then
    npm_path="$(printf '%s\n' "$nvm_dir"/versions/node/*/bin/npm | sort -V | tail -n 1)"
    dirname "$npm_path"
    return 0
  fi

  return 1
}

ensure_command_available() {
  local command_name="$1"
  local node_bin_dir=""

  if command -v "$command_name" >/dev/null 2>&1; then
    return 0
  fi

  case "$command_name" in
    node|npm|npx)
      node_bin_dir="$(find_nvm_node_bin_dir || true)"
      if [[ -n "$node_bin_dir" ]]; then
        prepend_to_path "$node_bin_dir"
      fi
      ;;
  esac

  command -v "$command_name" >/dev/null 2>&1
}

require_command() {
  local command_name="$1"

  if ! ensure_command_available "$command_name"; then
    echo "Missing required command: $command_name" >&2
    exit 1
  fi
}

read_env_value() {
  local file_path="$1"
  local key="$2"
  local default_value="$3"
  local value=""

  if [[ -f "$file_path" ]]; then
    value="$(sed -n "s/^${key}=//p" "$file_path" | tail -n 1)"
  fi

  if [[ -n "$value" ]]; then
    printf '%s\n' "$value"
    return
  fi

  printf '%s\n' "$default_value"
}

wait_for_backend() {
  local host="$1"
  local port="$2"
  local attempts=30

  for ((i = 1; i <= attempts; i++)); do
    if (echo >"/dev/tcp/$host/$port") >/dev/null 2>&1; then
      return 0
    fi

    if [[ -n "$backend_pid" ]] && ! kill -0 "$backend_pid" >/dev/null 2>&1; then
      echo "Backend process exited before becoming available." >&2
      return 1
    fi

    sleep 1
  done

  echo "Timed out waiting for backend on $host:$port." >&2
  return 1
}

cleanup() {
  if [[ -n "$frontend_pid" ]] && kill -0 "$frontend_pid" >/dev/null 2>&1; then
    echo
    echo "Stopping frontend..."
    kill "$frontend_pid" >/dev/null 2>&1 || true
    wait "$frontend_pid" >/dev/null 2>&1 || true
  fi

  if [[ -n "$backend_pid" ]] && kill -0 "$backend_pid" >/dev/null 2>&1; then
    echo "Stopping backend..."
    kill "$backend_pid" >/dev/null 2>&1 || true
    wait "$backend_pid" >/dev/null 2>&1 || true
  fi

  rm -f "$BACKEND_PID_FILE" "$FRONTEND_PID_FILE"
}

trap cleanup EXIT INT TERM

require_command go
require_command npm

if [[ ! -d "$BACKEND_DIR" ]]; then
  echo "Backend directory not found: $BACKEND_DIR" >&2
  exit 1
fi

if [[ ! -d "$FRONTEND_DIR" ]]; then
  echo "Frontend directory not found: $FRONTEND_DIR" >&2
  exit 1
fi

if [[ ! -f "$BACKEND_DIR/.env" ]]; then
  echo "Missing backend env file: $BACKEND_DIR/.env" >&2
  exit 1
fi

if [[ ! -f "$FRONTEND_DIR/.env" ]]; then
  echo "Missing frontend env file: $FRONTEND_DIR/.env" >&2
  exit 1
fi

mkdir -p "$RUNTIME_DIR"

backend_host="127.0.0.1"
backend_port="$(read_env_value "$BACKEND_DIR/.env" "APP_PORT" "8080")"
frontend_url="$(read_env_value "$FRONTEND_DIR/.env" "VITE_API_URL" "http://localhost:8080/api")"

echo "Starting backend on http://$backend_host:$backend_port ..."
(
  cd "$BACKEND_DIR"
  go run .
) &
backend_pid="$!"
printf '%s\n' "$backend_pid" >"$BACKEND_PID_FILE"

wait_for_backend "$backend_host" "$backend_port"
echo "Backend is up."

echo "Starting frontend (API: $frontend_url) ..."
(
  cd "$FRONTEND_DIR"
  npm run dev
) &
frontend_pid="$!"
printf '%s\n' "$frontend_pid" >"$FRONTEND_PID_FILE"

echo "Frontend is up."
echo "Run ./stop-dev.sh to stop both services."

wait "$backend_pid" "$frontend_pid"
