#!/usr/bin/env bash

set -Eeuo pipefail
umask 077

fail() {
  echo "Error: $*" >&2
  exit 1
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"
}

[[ "$(uname -s)" == Linux ]] || fail "deployment must run on the Linux server"
[[ "$(uname -m)" == x86_64 ]] || fail "expected x86_64 architecture"
[[ "$(id -un)" == tarot-deploy ]] || fail "run this script as tarot-deploy"

for command_name in git go npm node tar sha256sum awk ssh flock sudo install mktemp readlink df tee find; do
  require_command "$command_name"
done

repo_root=$(git rev-parse --show-toplevel 2>/dev/null) || fail "run this script inside the repository clone"
cd "$repo_root"

expected_origin=git@github.com:xamesLab/tarot-project-fullstack-vue3-go.git
origin=$(git remote get-url origin 2>/dev/null) || fail "origin remote is not configured"
[[ "$origin" == "$expected_origin" ]] || fail "unexpected origin remote: $origin"

lock_dir=$HOME/.cache
mkdir -p "$lock_dir"
chmod 0700 "$lock_dir"
if [[ ${TAROT_DEPLOY_LOCK_HELD:-} == 1 ]]; then
  [[ -e /proc/$$/fd/9 ]] || fail "deployment lock was lost during restart"
else
  exec 9>"$lock_dir/tarot-deploy.lock"
  flock -n 9 || fail "another deployment is already running"
  export TAROT_DEPLOY_LOCK_HELD=1
fi

if [[ -n $(git status --porcelain) ]]; then
  echo "The deployment clone has local changes:" >&2
  git status --short >&2
  fail "restore a clean deployment clone before continuing"
fi

current_branch=$(git branch --show-current)
[[ "$current_branch" == main ]] || fail "deployment clone must be on main (currently: $current_branch)"

export GIT_SSH_COMMAND='ssh -o BatchMode=yes -o StrictHostKeyChecking=yes'
echo "Fetching origin/main..."
git fetch --quiet --no-tags origin +refs/heads/main:refs/remotes/origin/main
release_id=$(git rev-parse refs/remotes/origin/main)
[[ "$release_id" =~ ^[0-9a-f]{40}$ ]] || fail "GitHub returned an invalid commit SHA"

if [[ ${TAROT_DEPLOY_SOURCE_SHA:-} != "$release_id" ]]; then
  echo "Checking out release $release_id..."
  git reset --hard "$release_id"
  export TAROT_DEPLOY_SOURCE_SHA="$release_id"
  exec bash "$repo_root/deploy.sh"
fi

[[ "$(git rev-parse HEAD)" == "$release_id" ]] || fail "checkout does not match the selected SHA"
[[ -z $(git status --porcelain) ]] || fail "checkout changed after fetch"

incoming_dir=/var/lib/tarot-deploy/incoming
[[ -d "$incoming_dir" && -w "$incoming_dir" ]] || fail "cannot write to $incoming_dir"

for env_override in \
  "$repo_root/frontend/.env" \
  "$repo_root/frontend/.env.local" \
  "$repo_root/frontend/.env.test.local" \
  "$repo_root/frontend/.env.production.local"; do
  [[ ! -e "$env_override" ]] || fail "remove local Vite environment override: ${env_override#"$repo_root"/}"
done

build_dir=$(mktemp -d "${TMPDIR:-/tmp}/tarot-release.XXXXXXXX")
cleanup() {
  rm -rf -- "$build_dir"
}
trap cleanup EXIT

release_dir=$build_dir/release
mkdir -p "$release_dir/backend" "$release_dir/frontend"

node -e 'const [major, minor] = process.versions.node.split(".").map(Number); process.exit((major === 20 && minor >= 19) || (major === 22 && minor >= 12) || major > 22 ? 0 : 1)' \
  || fail "Node.js 20.19+ or 22.12+ is required (found: $(node --version))"

print_resources() {
  echo "Node.js: $(node --version)"
  awk '/^(MemTotal|MemAvailable|SwapTotal|SwapFree):/ { print }' /proc/meminfo
  awk '{ print "Load average (1/5/15 min):", $1, $2, $3 }' /proc/loadavg
  df -Pm "$repo_root" "$build_dir" "$incoming_dir"
}

echo "Resources before build:"
print_resources

echo "Testing backend..."
(
  cd "$repo_root/backend"
  CGO_ENABLED=0 go test ./...
)

echo "Building backend..."
(
  cd "$repo_root/backend"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
      -o "$release_dir/backend/tarot-app" .
)

echo "Installing frontend dependencies..."
npm --prefix "$repo_root/frontend" ci --include=dev

echo "Checking frontend..."
npm --prefix "$repo_root/frontend" run lint

clear_vite_environment() {
  while IFS= read -r variable_name; do
    if [[ "$variable_name" == VITE_* ]]; then
      unset "$variable_name"
    fi
  done < <(compgen -e)
}

clear_vite_environment
test_log=$(mktemp "$lock_dir/tarot-deploy-tests-$release_id.XXXXXXXX.log")
echo "Testing frontend sequentially; log: $test_log"
# Keep diagnostics outside the checkout and temporary build directory on failure.
# pipefail ensures tee cannot hide a failed test command.
if ! {
  print_resources &&
  npm --prefix "$repo_root/frontend" run test:deploy
} 2>&1 | tee "$test_log"; then
  fail "frontend checks failed; see $test_log"
fi

echo "Building production frontend..."
npm --prefix "$repo_root/frontend" run build
cp -a "$repo_root/frontend/dist/." "$release_dir/frontend/"

archive_name=tarot-$release_id.tar.gz
checksum_name=tarot-$release_id.sha256
archive=$build_dir/$archive_name
checksum=$build_dir/$checksum_name

echo "Creating release archive and checksum..."
# The private build/log umask must not make published runtime assets unreadable.
# Normalize only the assembled release; archives and incoming files stay private.
find "$release_dir/backend" "$release_dir/frontend" -type d -exec chmod 0755 {} +
find "$release_dir/backend" "$release_dir/frontend" -type f -exec chmod 0644 {} +
chmod 0755 "$release_dir/backend/tarot-app"
tar -czf "$archive" -C "$release_dir" backend frontend
archive_hash=$(sha256sum "$archive" | awk '{print $1}')
printf '%s  %s\n' "$archive_hash" "$archive_name" > "$checksum"

incoming_archive=$incoming_dir/$archive_name
incoming_checksum=$incoming_dir/$checksum_name
install -m 0600 -- "$archive" "$incoming_archive"
install -m 0600 -- "$checksum" "$incoming_checksum"

active_release=$(readlink -f /opt/tarot-app/current 2>/dev/null || true)
if [[ "$active_release" == "/opt/tarot-app/releases/$release_id" ]] && /usr/local/sbin/smoke-test-tarot; then
  rm -f -- "$incoming_archive" "$incoming_checksum"
  echo "Release $release_id is already active and passed smoke tests."
  exit 0
fi

echo "Installing release $release_id..."
sudo -n /usr/local/sbin/install-tarot-release "$release_id" "$incoming_archive" "$incoming_checksum"

echo
echo "Deployment complete: $release_id"
echo "https://nebula-tarot.ru"
