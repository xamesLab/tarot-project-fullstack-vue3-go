#!/usr/bin/env bash
set -Eeuo pipefail

attempts=${SMOKE_ATTEMPTS:-15}
delay=${SMOKE_DELAY_SECONDS:-2}
local_health=http://127.0.0.1:8888/api/health
public_health=https://nebula-tarot.ru/api/health
public_home=https://nebula-tarot.ru/

if ! systemctl is-active --quiet tarot-app.service; then
  echo "tarot-app.service is not active" >&2
  exit 1
fi

for ((attempt = 1; attempt <= attempts; attempt++)); do
  if curl --fail --silent --show-error --max-time 5 "$local_health" >/dev/null; then
    break
  fi
  if (( attempt == attempts )); then
    echo "local health check failed after $attempts attempts" >&2
    exit 1
  fi
  sleep "$delay"
done

curl --fail --silent --show-error --max-time 10 "$public_health" >/dev/null
curl --fail --silent --show-error --max-time 10 "$public_home" | grep -qi '<html'

echo "smoke tests passed"
