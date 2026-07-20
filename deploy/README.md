# Deployment assets

Production releases are built by GitHub Actions and installed under
`/opt/tarot-app/releases/<commit-sha>`. The active release is selected by the
`/opt/tarot-app/current` symlink.

Runtime state is kept outside releases:

- environment: `/etc/tarot-app/app.env`;
- uploads: `/var/lib/tarot-app/uploads`;
- MongoDB data: `/var/lib/mongodb`.

The installer and smoke-test scripts must be installed as root-owned files:

```text
/usr/local/sbin/install-tarot-release
/usr/local/sbin/smoke-test-tarot
```

The deployment account may receive passwordless sudo permission for the exact
installer command only. Do not grant it unrestricted sudo access.

The incoming directory is `/var/lib/tarot-deploy/incoming`. The installer
rejects release inputs outside this directory and accepts only regular files
and directories inside the release archive.

Operational scripts and units are installed as root-owned files:

```text
/usr/local/sbin/backup-tarot
/usr/local/sbin/health-check-tarot
/etc/systemd/system/tarot-backup.service
/etc/systemd/system/tarot-backup.timer
/etc/systemd/system/tarot-health.service
/etc/systemd/system/tarot-health.timer
```

The backup timer creates a daily all-database MongoDB archive and an uploads
archive under `/var/backups/tarot-app`. A restore dry-run and SHA-256 checks
must pass before a backup becomes visible. Seven scheduled generations are
retained; named pre-deployment backups are not included in that rotation.
MongoDB tools read the dedicated backup credentials from the root-only
`/etc/tarot-app/mongodump.yml`; credentials are never passed on the command
line.

The health timer verifies all four services, local and public health endpoints,
and a 512 MiB minimum free-disk threshold every five minutes. It writes results
to journald; external alert delivery remains a separate infrastructure concern.

TLS is intentionally not managed by these files. HTTPS for `nebula-tarot.ru`
is terminated outside this Nginx server configuration.
