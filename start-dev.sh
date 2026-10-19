#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"
RUNTIME_DIR="$ROOT_DIR/.dev-runtime"
BACKEND_PID_FILE="$RUNTIME_DIR/backend.pid"
FRONTEND_PID_FILE="$RUNTIME_DIR/frontend.pid"
BACKEND_LOG_FILE="$RUNTIME_DIR/backend.log"

backend_pid=""
frontend_pid=""
backend_log_pid=""

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
    if [[ -n "$backend_pid" ]] && ! kill -0 "$backend_pid" >/dev/null 2>&1; then
      echo "Backend process exited before becoming available." >&2
      return 2
    fi

    if (echo >"/dev/tcp/$host/$port") >/dev/null 2>&1; then
      return 0
    fi

    sleep 1
  done

  echo "Timed out waiting for backend on $host:$port." >&2
  return 3
}

stop_backend_log_tail() {
  if [[ -n "$backend_log_pid" ]] && kill -0 "$backend_log_pid" >/dev/null 2>&1; then
    kill "$backend_log_pid" >/dev/null 2>&1 || true
    wait "$backend_log_pid" >/dev/null 2>&1 || true
  fi
  backend_log_pid=""
}

stop_process_tree() {
  local pid="$1"
  local child_pid

  while IFS= read -r child_pid; do
    [[ -n "$child_pid" ]] && stop_process_tree "$child_pid"
  done < <(ps -ef | awk -v parent="$pid" 'NR > 1 && $3 == parent { print $2 }')

  kill "$pid" >/dev/null 2>&1 || true
}

cleanup() {
  if [[ -n "$frontend_pid" ]] && kill -0 "$frontend_pid" >/dev/null 2>&1; then
    echo
    echo "Stopping frontend..."
    stop_process_tree "$frontend_pid"
    wait "$frontend_pid" >/dev/null 2>&1 || true
  fi

  if [[ -n "$backend_pid" ]] && kill -0 "$backend_pid" >/dev/null 2>&1; then
    echo "Stopping backend..."
    stop_process_tree "$backend_pid"
    wait "$backend_pid" >/dev/null 2>&1 || true
  fi

  stop_backend_log_tail
  rm -f "$BACKEND_PID_FILE" "$FRONTEND_PID_FILE" "$BACKEND_LOG_FILE"
}

trap cleanup EXIT INT TERM

require_command npm

requested_data_mode="${DEV_DATA_MODE:-auto}"
case "$requested_data_mode" in
  auto|mongo|mock) ;;
  *)
    echo "Invalid DEV_DATA_MODE '$requested_data_mode'. Use auto, mongo, or mock." >&2
    exit 1
    ;;
esac

if [[ ! -d "$BACKEND_DIR" ]]; then
  echo "Backend directory not found: $BACKEND_DIR" >&2
  exit 1
fi

if [[ ! -d "$FRONTEND_DIR" ]]; then
  echo "Frontend directory not found: $FRONTEND_DIR" >&2
  exit 1
fi

mkdir -p "$RUNTIME_DIR"
: >"$BACKEND_LOG_FILE"

data_mode="$requested_data_mode"
if [[ "$data_mode" == "mongo" && ! -f "$BACKEND_DIR/.env" ]]; then
  echo "Missing backend env file required by DEV_DATA_MODE=mongo: $BACKEND_DIR/.env" >&2
  exit 1
fi

if [[ "$data_mode" == "auto" && ! -f "$BACKEND_DIR/.env" ]]; then
  echo "Backend env file not found; starting in mock mode."
  data_mode="mock"
fi

if [[ "$data_mode" != "mock" ]]; then
  require_command go
fi

backend_host="127.0.0.1"
backend_port="$(read_env_value "$BACKEND_DIR/.env" "APP_PORT" "8888")"
frontend_url="$(read_env_value "$FRONTEND_DIR/.env" "VITE_API_URL" "http://localhost:8888/api")"
if [[ -n "${VITE_API_URL:-}" ]]; then
  frontend_url="$VITE_API_URL"
fi

frontend_api_mode="${VITE_API_MODE:-}"
if [[ -z "$frontend_api_mode" ]]; then
  frontend_api_mode="$(read_env_value "$FRONTEND_DIR/.env" "VITE_API_MODE" "mock")"
fi

frontend_auth_api_mode="${VITE_AUTH_API_MODE:-}"
if [[ -z "$frontend_auth_api_mode" ]]; then
  frontend_auth_api_mode="$(read_env_value "$FRONTEND_DIR/.env" "VITE_AUTH_API_MODE" "http")"
fi

if [[ "$data_mode" != "mock" ]]; then
  echo "Starting backend on http://$backend_host:$backend_port ..."
  tail -n 0 -f "$BACKEND_LOG_FILE" &
  backend_log_pid="$!"
  (
    cd "$BACKEND_DIR"
    exec env GO111MODULE=on go run .
  ) >"$BACKEND_LOG_FILE" 2>&1 &
  backend_pid="$!"
  printf '%s\n' "$backend_pid" >"$BACKEND_PID_FILE"

  if wait_for_backend "$backend_host" "$backend_port"; then
    echo "Backend is up."
  else
    backend_wait_result="$?"
    wait "$backend_pid" >/dev/null 2>&1 || true

    if [[ "$requested_data_mode" == "auto" && "$backend_wait_result" -eq 2 ]] && \
      grep -Eqi 'load config: (connect mongo:|ping mongo:|missing required environment variable: MONGO_URI)' "$BACKEND_LOG_FILE"; then
      stop_backend_log_tail
      echo "MongoDB is not configured or unavailable; continuing with the frontend mock API."
      data_mode="mock"
      backend_pid=""
      rm -f "$BACKEND_PID_FILE"
    else
      stop_backend_log_tail
      cat "$BACKEND_LOG_FILE" >&2
      echo "Backend failed to start. Fix the backend error or use DEV_DATA_MODE=mock to run the UI with mocks." >&2
      exit 1
    fi
  fi
fi

if [[ "$data_mode" == "mock" ]]; then
  frontend_api_mode="mock"
  frontend_auth_api_mode="mock"
  echo "Data mode: mock. The Go backend is not running; frontend changes are temporary."
elif [[ "$requested_data_mode" == "mongo" ]]; then
  frontend_auth_api_mode="http"
  echo "Data mode: MongoDB."
else
  echo "Data mode: backend available; auth API mode is $frontend_auth_api_mode."
fi

echo "Starting frontend (API: $frontend_url) ..."
(
  cd "$FRONTEND_DIR"
  VITE_API_URL="$frontend_url" \
    VITE_API_MODE="$frontend_api_mode" \
    VITE_AUTH_API_MODE="$frontend_auth_api_mode" \
    npm run dev
) &
frontend_pid="$!"
printf '%s\n' "$frontend_pid" >"$FRONTEND_PID_FILE"

echo "Frontend is up."
echo "Run ./stop-dev.sh to stop both services."

if [[ -n "$backend_pid" ]]; then
  wait "$backend_pid" "$frontend_pid"
else
  wait "$frontend_pid"
fi
