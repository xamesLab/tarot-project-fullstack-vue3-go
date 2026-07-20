#!/usr/bin/env bash
set -Eeuo pipefail

minimum_free_kb=${MINIMUM_FREE_KB:-524288}

systemctl is-active --quiet tarot-app.service
systemctl is-active --quiet mongod.service
systemctl is-active --quiet nginx.service
systemctl is-active --quiet gtwapi.service

curl --fail --silent --show-error --max-time 10 \
  http://127.0.0.1:8888/api/health >/dev/null
curl --fail --silent --show-error --max-time 15 \
  https://nebula-tarot.ru/api/health >/dev/null

free_kb=$(df --output=avail / | tail -n 1 | tr -d ' ')
if (( free_kb < minimum_free_kb )); then
  echo "free disk space is below threshold: ${free_kb} KB" >&2
  exit 1
fi

echo "health check passed; free_kb=$free_kb"
