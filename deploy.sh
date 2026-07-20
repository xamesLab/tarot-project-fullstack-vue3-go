#!/usr/bin/env bash

set -Eeuo pipefail

repo_root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
cd "$repo_root"

deploy_host=${DEPLOY_HOST:-d51673dc67a4.vps.myjino.ru}
deploy_port=${DEPLOY_PORT:-49401}
deploy_user=${DEPLOY_USER:-root}
deploy_key=${DEPLOY_KEY:-$HOME/.ssh/id_rsa}
deploy_branch=${DEPLOY_BRANCH:-main}
incoming_dir=/var/lib/tarot-deploy/incoming

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Не найдена обязательная команда: $1" >&2
    exit 1
  fi
}

for command_name in git go npm tar sha256sum ssh scp mktemp; do
  require_command "$command_name"
done

if [[ ! -f $deploy_key ]]; then
  echo "SSH-ключ не найден: $deploy_key" >&2
  echo "Укажите путь через DEPLOY_KEY." >&2
  exit 1
fi

current_branch=$(git branch --show-current)
if [[ $current_branch != "$deploy_branch" ]]; then
  echo "Деплой разрешён только из ветки $deploy_branch (сейчас: $current_branch)." >&2
  echo "При необходимости переопределите DEPLOY_BRANCH." >&2
  exit 1
fi

if [[ -n $(git status --porcelain) ]]; then
  echo "В рабочем дереве есть незакоммиченные изменения:" >&2
  git status --short >&2
  exit 1
fi

if ! git rev-parse --verify '@{upstream}' >/dev/null 2>&1; then
  echo "Для ветки $current_branch не настроена upstream-ветка." >&2
  exit 1
fi

echo "Обновляю сведения о remote..."
git fetch --quiet

release_id=$(git rev-parse HEAD)
upstream_id=$(git rev-parse '@{upstream}')
if [[ $release_id != "$upstream_id" ]]; then
  echo "Текущий commit $release_id ещё не совпадает с @{upstream} ($upstream_id)." >&2
  echo "Сначала выполните git push." >&2
  exit 1
fi

build_dir=$(mktemp -d "${TMPDIR:-/tmp}/tarot-release.XXXXXXXX")
cleanup() {
  rm -rf -- "$build_dir"
}
trap cleanup EXIT

release_dir=$build_dir/release
archive_name=tarot-$release_id.tar.gz
checksum_name=tarot-$release_id.sha256
archive=$build_dir/$archive_name
checksum=$build_dir/$checksum_name

mkdir -p "$release_dir/backend" "$release_dir/frontend"

echo "Проверяю backend..."
(
  cd "$repo_root/backend"
  go test ./...
)

echo "Собираю backend для linux/amd64..."
(
  cd "$repo_root/backend"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
      -o "$release_dir/backend/tarot-app" .
)

echo "Устанавливаю frontend-зависимости..."
npm --prefix "$repo_root/frontend" ci

echo "Проверяю frontend..."
npm --prefix "$repo_root/frontend" run lint
npm --prefix "$repo_root/frontend" test

echo "Собираю frontend..."
npm --prefix "$repo_root/frontend" run build
cp -a "$repo_root/frontend/dist/." "$release_dir/frontend/"

echo "Формирую release-архив..."
tar -czf "$archive" -C "$release_dir" backend frontend
archive_hash=$(sha256sum "$archive" | awk '{print $1}')
printf '%s  %s\n' "$archive_hash" "$archive_name" > "$checksum"

ssh_target=$deploy_user@$deploy_host
ssh_options=(-p "$deploy_port" -i "$deploy_key" -o BatchMode=yes)
scp_options=(-P "$deploy_port" -i "$deploy_key" -o BatchMode=yes)

echo "Загружаю релиз $release_id на $deploy_host:$deploy_port..."
scp "${scp_options[@]}" \
  "$archive" "$checksum" \
  "$ssh_target:$incoming_dir/"

remote_archive=$incoming_dir/$archive_name
remote_checksum=$incoming_dir/$checksum_name
installer="/usr/local/sbin/install-tarot-release"
if [[ $deploy_user != root ]]; then
  installer="sudo $installer"
fi

echo "Активирую релиз..."
ssh "${ssh_options[@]}" "$ssh_target" \
  "$installer '$release_id' '$remote_archive' '$remote_checksum'"

echo
echo "Деплой завершён: $release_id"
echo "https://nebula-tarot.ru"
