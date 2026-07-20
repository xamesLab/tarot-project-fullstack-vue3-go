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

TLS is intentionally not managed by these files. HTTPS for `nebula-tarot.ru`
is terminated outside this Nginx server configuration.
