#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RUNTIME_DIR="$ROOT_DIR/.dev-runtime"
BACKEND_PID_FILE="$RUNTIME_DIR/backend.pid"
FRONTEND_PID_FILE="$RUNTIME_DIR/frontend.pid"

stop_process_tree() {
  local pid="$1"
  local child_pid

  while IFS= read -r child_pid; do
    [[ -n "$child_pid" ]] && stop_process_tree "$child_pid"
  done < <(ps -ef | awk -v parent="$pid" 'NR > 1 && $3 == parent { print $2 }')

  kill "$pid" >/dev/null 2>&1 || true
}

stop_process() {
  local name="$1"
  local pid_file="$2"

  if [[ ! -f "$pid_file" ]]; then
    echo "$name is not running (missing PID file)."
    return
  fi

  local pid
  pid="$(cat "$pid_file")"

  if [[ -z "$pid" ]]; then
    echo "$name PID file is empty."
    rm -f "$pid_file"
    return
  fi

  if ! kill -0 "$pid" >/dev/null 2>&1; then
    echo "$name is not running (stale PID $pid)."
    rm -f "$pid_file"
    return
  fi

  echo "Stopping $name (PID $pid)..."
  stop_process_tree "$pid"

  for _ in {1..20}; do
    if ! kill -0 "$pid" >/dev/null 2>&1; then
      rm -f "$pid_file"
      echo "$name stopped."
      return
    fi
    sleep 0.5
  done

  echo "$name did not stop gracefully, sending SIGKILL..."
  kill -9 "$pid" >/dev/null 2>&1 || true
  rm -f "$pid_file"
}

stop_process "frontend" "$FRONTEND_PID_FILE"
stop_process "backend" "$BACKEND_PID_FILE"

if [[ -d "$RUNTIME_DIR" ]] && [[ -z "$(ls -A "$RUNTIME_DIR")" ]]; then
  rmdir "$RUNTIME_DIR"
fi
