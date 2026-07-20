#!/usr/bin/env bash
set -Eeuo pipefail

if [[ $EUID -ne 0 ]]; then
  echo "install-release must run as root" >&2
  exit 1
fi

if [[ $# -ne 3 ]]; then
  echo "usage: install-release <40-char-commit-sha> <archive> <checksum-file>" >&2
  exit 1
fi

release_id=$1
archive=$2
checksum_file=$3
app_root=/opt/tarot-app
releases_dir=$app_root/releases
current_link=$app_root/current
uploads_dir=/var/lib/tarot-app/uploads
release_dir=$releases_dir/$release_id
staging_dir=$releases_dir/.staging-$release_id-$$
lock_file=/run/lock/tarot-app-deploy.lock

if [[ ! $release_id =~ ^[0-9a-f]{40}$ ]]; then
  echo "invalid release id: $release_id" >&2
  exit 1
fi

for required_file in "$archive" "$checksum_file"; do
  if [[ ! -f $required_file ]]; then
    echo "missing file: $required_file" >&2
    exit 1
  fi
done

exec 9>"$lock_file"
if ! flock -n 9; then
  echo "another deployment is already running" >&2
  exit 1
fi

cleanup() {
  rm -rf -- "$staging_dir"
}
trap cleanup EXIT

cd "$(dirname "$archive")"
sha256sum --check "$(basename "$checksum_file")"

if tar -tzf "$archive" | grep -Eq '(^/|(^|/)\.\.(/|$))'; then
  echo "archive contains an unsafe path" >&2
  exit 1
fi

install -d -m 0755 "$releases_dir" "$uploads_dir"
chown tarot-user:tarot-user "$uploads_dir"

if [[ ! -d $release_dir ]]; then
  install -d -m 0755 "$staging_dir"
  tar -xzf "$archive" -C "$staging_dir"

  if [[ ! -x $staging_dir/backend/tarot-app ]]; then
    echo "release does not contain an executable backend" >&2
    exit 1
  fi
  if [[ ! -f $staging_dir/frontend/index.html ]]; then
    echo "release does not contain frontend/index.html" >&2
    exit 1
  fi

  rm -rf -- "$staging_dir/backend/uploads"
  ln -s "$uploads_dir" "$staging_dir/backend/uploads"
  chown -R root:root "$staging_dir"
  chmod 0755 "$staging_dir/backend/tarot-app"
  mv "$staging_dir" "$release_dir"
fi

previous_release=
if [[ -L $current_link ]]; then
  previous_release=$(readlink -f "$current_link" || true)
fi

new_link=$app_root/.current-$release_id-$$
ln -s "$release_dir" "$new_link"
mv -Tf "$new_link" "$current_link"

rollback() {
  if [[ -n $previous_release && -d $previous_release ]]; then
    rollback_link=$app_root/.rollback-$release_id-$$
    ln -s "$previous_release" "$rollback_link"
    mv -Tf "$rollback_link" "$current_link"
    systemctl restart tarot-app.service || true
  fi
}

if ! systemctl restart tarot-app.service; then
  rollback
  exit 1
fi

if ! /usr/local/sbin/smoke-test-tarot; then
  rollback
  /usr/local/sbin/smoke-test-tarot || true
  exit 1
fi

rm -f -- "$archive" "$checksum_file"
echo "release $release_id installed successfully"
