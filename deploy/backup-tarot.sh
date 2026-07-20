#!/usr/bin/env bash
set -Eeuo pipefail

backup_root=${BACKUP_ROOT:-/var/backups/tarot-app}
uploads_root=${UPLOADS_ROOT:-/var/lib/tarot-app/uploads}
keep_count=${BACKUP_KEEP_COUNT:-7}
lock_file=/run/lock/tarot-app-backup.lock
timestamp=$(date -u +%Y%m%d-%H%M%S)
final_dir="$backup_root/$timestamp"
work_dir="$backup_root/.incomplete-$timestamp-$$"

exec 9>"$lock_file"
flock -n 9 || {
  echo "another tarot backup is already running" >&2
  exit 1
}

install -d -o root -g root -m 0700 "$backup_root"
mkdir -m 0700 "$work_dir"
trap 'rm -rf -- "$work_dir"' EXIT

mongodump \
  --host 127.0.0.1 \
  --port 27017 \
  --archive="$work_dir/mongodb-all.archive.gz" \
  --gzip

tar --create --gzip --file "$work_dir/uploads.tar.gz" \
  --directory "$uploads_root" .

mongorestore \
  --host 127.0.0.1 \
  --port 27017 \
  --archive="$work_dir/mongodb-all.archive.gz" \
  --gzip \
  --dryRun >/dev/null

(
  cd "$work_dir"
  sha256sum mongodb-all.archive.gz uploads.tar.gz > SHA256SUMS
  sha256sum --check SHA256SUMS >/dev/null
  {
    printf 'created_utc=%s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    printf 'host=%s\n' "$(hostname -f)"
    printf 'mongodb_bytes=%s\n' "$(stat -c %s mongodb-all.archive.gz)"
    printf 'uploads_bytes=%s\n' "$(stat -c %s uploads.tar.gz)"
    printf 'restore_dry_run=passed\n'
  } > MANIFEST
)

mv "$work_dir" "$final_dir"
trap - EXIT

mapfile -t old_backups < <(
  find "$backup_root" -mindepth 1 -maxdepth 1 -type d \
    -regextype posix-extended \
    -regex '.*/[0-9]{8}-[0-9]{6}' \
    -printf '%f\n' | sort -r | tail -n "+$((keep_count + 1))"
)

for old_backup in "${old_backups[@]}"; do
  rm -rf -- "$backup_root/$old_backup"
done

echo "backup completed: $final_dir"
