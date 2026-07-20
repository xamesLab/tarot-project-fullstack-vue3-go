# Production deployment: от Git push до обновления VPS

## 1. Назначение и границы процесса

Документ описывает фактически настроенный путь доставки Nebula Tarot:

```text
commit → push/PR → CI → merge в main
       → ручной Deploy production или push тега v*
       → повторные тесты → production build
       → release archive + SHA-256
       → SCP на VPS:49401
       → проверка и staging
       → атомарное переключение current
       → restart + smoke tests
       → успешный deploy или автоматический rollback
```

Обычный release обновляет только frontend и Go-backend. Он не заменяет:

- `/etc/tarot-app/app.env`;
- MongoDB;
- `/var/lib/tarot-app/uploads`;
- Nginx и внешний TLS;
- `/opt/gtwapi` и `gtwapi.service`.

Основные файлы автоматизации:

- `.github/workflows/ci.yml`;
- `.github/workflows/deploy-production.yml`;
- `deploy/install-release.sh`;
- `deploy/smoke-test.sh`;
- `deploy/tarot-app.service`;
- `deploy/nginx-nebula-tarot.conf`.

## 2. От push до разрешения на production deploy

### 2.1. CI после push и pull request

Workflow `CI` запускается:

- для каждого pull request;
- для push в `main`;
- для push в `codex/**`.

Backend job:

1. checkout точного commit;
2. установка Go по `backend/go.mod`;
3. `go test ./...`;
4. `go vet ./...`;
5. `go build ./...`.

Frontend job:

1. checkout того же commit;
2. установка Node.js по `.nvmrc`;
3. `npm ci`;
4. `npm audit --omit=dev --audit-level=high`;
5. `npm run typecheck`;
6. `npm run lint`;
7. `npm test`;
8. `npm run build`.

Оба jobs рекомендуется сделать обязательными в branch protection для `main`.
Merge не должен выполняться при failed CI.

### 2.2. Merge в main не обновляет production

Push в `main` запускает CI, но не production deployment. Production
обновляется только одним из способов:

- вручную: GitHub Actions → `Deploy production` → `Run workflow`;
- push тега с префиксом `v`, например `v1.4.0`.

Для ручного запуска нужно выбирать `main` либо конкретный проверенный tag.

Пример release через tag:

```bash
git switch main
git pull --ff-only origin main
git tag -a v1.4.0 -m "Release v1.4.0"
git push origin v1.4.0
```

Тег должен указывать на commit, уже прошедший CI.

## 3. GitHub Environment production

Deployment job использует GitHub Environment `production`.

| GitHub secret | Назначение |
| --- | --- |
| `VPS_HOST` | SSH hostname VPS |
| `VPS_PORT` | внешний SSH-порт `49401` |
| `VPS_USER` | пользователь `tarot-deploy` |
| `VPS_SSH_KEY` | private deploy key |
| `VPS_HOST_KEY` | закреплённый SSH host key |

Рекомендуется включить required reviewer для Environment `production`.
Проверка host key не отключается.

Эти secrets существуют только для доставки release. Они не копируются в
application env и не попадают в artifact.

Перед первым автоматическим deploy нужно убедиться, что workflow-файлы уже
merged в default branch, а все пять secrets добавлены в Environment
`production`. До выполнения этих prerequisites используется только ручная
установка release на VPS.

## 4. Повторные проверки и production build

Workflow `Deploy production` повторно выполняет тесты. Это гарантирует, что
release собирается из того же `GITHUB_SHA`, который проверяется внутри
production job.

### 4.1. Backend artifact

Backend собирается со следующими параметрами:

```text
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
```

Результат:

```text
artifact/backend/tarot-app
```

Используются `-trimpath` и `-ldflags='-s -w'`.

### 4.2. Frontend artifact и frontend env

После `npm ci`, audit, typecheck, lint и tests выполняется `npm run build`.
Содержимое `frontend/dist` копируется в:

```text
artifact/frontend/
```

Vite во время сборки читает tracked-файл:

```text
frontend/.env.production
```

Его `VITE_*` значения встраиваются в JavaScript. На VPS runtime env для
frontend отсутствует. Изменение `VITE_*` требует commit, новой сборки и
нового deploy.

Текущий production frontend использует относительный API URL `/api`. Браузер
обращается к `nebula-tarot.ru`, а Nginx проксирует `/api/` в Go-backend.

Текущие публичные build-time значения:

```dotenv
VITE_API_MODE=mock
VITE_AUTH_API_MODE=http
VITE_API_URL=/api
```

### 4.3. Release archive

Создаются:

```text
tarot-app-<40-char-GITHUB_SHA>.tar.gz
tarot-app-<40-char-GITHUB_SHA>.tar.gz.sha256
```

Внутри archive разрешены только:

```text
backend/
└── tarot-app

frontend/
├── index.html
└── assets/
```

Env, uploads и MongoDB в archive не входят. GitHub Actions artifact хранится
30 дней.

## 5. Доставка на VPS

GitHub Actions подключается к:

```text
tarot-deploy@d51673dc67a4.vps.myjino.ru:49401
```

Archive и checksum загружаются в:

```text
/var/lib/tarot-deploy/incoming/
```

Upload ещё не меняет production. После него workflow выполняет:

```bash
sudo /usr/local/sbin/install-tarot-release \
  <commit-sha> \
  /var/lib/tarot-deploy/incoming/tarot-app-<sha>.tar.gz \
  /var/lib/tarot-deploy/incoming/tarot-app-<sha>.tar.gz.sha256
```

`tarot-deploy` не имеет произвольного root shell. Sudo разрешает только
`/usr/local/sbin/install-tarot-release *`.

## 6. Проверки и staging на сервере

Installer:

1. принимает только 40-символьный lowercase commit SHA;
2. принимает archive/checksum только из `incoming`;
3. отклоняет symlink вместо входных файлов;
4. блокирует параллельный deploy через `flock`;
5. проверяет SHA-256;
6. отклоняет абсолютные пути и `..`;
7. разрешает только `backend/**` и `frontend/**`;
8. отклоняет links и неподдерживаемые типы файлов.

Распаковка сначала выполняется в:

```text
/opt/tarot-app/releases/.staging-<sha>-<pid>
```

Installer проверяет:

```text
backend/tarot-app
frontend/index.html
```

Uploads из artifact не используются. Вместо них создаётся:

```text
backend/uploads -> /var/lib/tarot-app/uploads
```

Release получает владельца `root:root`, backend binary — mode `0755`.
После проверок staging атомарно переименовывается в:

```text
/opt/tarot-app/releases/<sha>
```

## 7. Переключение active release

Текущая версия выбирается ссылкой:

```text
/opt/tarot-app/current
```

Installer запоминает её предыдущую цель, создаёт временную ссылку на новый
release и атомарно заменяет `current`.

Systemd использует:

```text
WorkingDirectory=/opt/tarot-app/current/backend
ExecStart=/opt/tarot-app/current/backend/tarot-app
```

Nginx использует:

```text
root /opt/tarot-app/current/frontend
```

После переключения выполняется:

```bash
systemctl restart tarot-app.service
```

Nginx reload для обычного application release не требуется.

## 8. Smoke tests и rollback

`/usr/local/sbin/smoke-test-tarot` проверяет:

1. `tarot-app.service` active;
2. `http://127.0.0.1:8888/api/health`;
3. `https://nebula-tarot.ru/api/health`;
4. главную страницу.

Local health повторяется до 15 раз с интервалом 2 секунды.

При ошибке installer:

1. атомарно возвращает `current` на предыдущий release;
2. повторно запускает `tarot-app.service`;
3. завершает GitHub Actions job с ошибкой;
4. не изменяет env, uploads и MongoDB.

После успеха входящий archive и checksum удаляются. Распакованный release
остаётся доступным для аудита и ручного rollback.

## 9. Структура проекта и данных на VPS

### 9.1. Полная схема

```text
/opt/tarot-app/
├── current -> /opt/tarot-app/releases/<active-sha>
├── releases/
│   ├── <previous-sha>/
│   │   ├── backend/
│   │   │   ├── tarot-app
│   │   │   └── uploads -> /var/lib/tarot-app/uploads
│   │   └── frontend/
│   │       ├── index.html
│   │       └── assets/
│   ├── <active-sha>/
│   │   ├── backend/
│   │   └── frontend/
│   └── legacy-<timestamp>/
└── backend/                     # старый путь, не active runtime

/etc/tarot-app/
├── app.env
└── mongodump.yml

/var/lib/tarot-app/
└── uploads/

/var/lib/tarot-deploy/
└── incoming/

/var/lib/mongodb/                # MongoDB data files

/var/backups/tarot-app/
├── pre-deploy-<timestamp>/
└── <YYYYMMDD-HHMMSS>/
    ├── mongodb-all.archive.gz
    ├── uploads.tar.gz
    ├── SHA256SUMS
    └── MANIFEST

/usr/local/sbin/
├── install-tarot-release
├── smoke-test-tarot
├── backup-tarot
└── health-check-tarot

/etc/systemd/system/
├── tarot-app.service
├── tarot-backup.service
├── tarot-backup.timer
├── tarot-health.service
└── tarot-health.timer

/etc/nginx/sites-available/
└── nebula-tarot.ru

/opt/gtwapi/                     # отдельный проект
├── .env
└── ...
```

### 9.2. Что является active source

Проверка active release:

```bash
readlink -f /opt/tarot-app/current
```

Не следует редактировать файлы внутри `releases/<sha>` вручную. Следующий
deploy создаёт отдельный release, а ручное изменение нарушает
воспроизводимость.

### 9.3. Постоянные ресурсы

| Ресурс | Путь |
| --- | --- |
| Backend env | `/etc/tarot-app/app.env` |
| Mongo backup credentials | `/etc/tarot-app/mongodump.yml` |
| Uploads | `/var/lib/tarot-app/uploads` |
| MongoDB data | `/var/lib/mongodb` |
| Backup | `/var/backups/tarot-app` |

Эти ресурсы переживают deploy и rollback и не должны находиться внутри
`releases/<sha>`.

## 10. Env-файлы

### 10.1. Главный production env backend

Путь:

```text
/etc/tarot-app/app.env
```

Владелец и права:

```text
root:tarot-user 0640
```

Systemd загружает его напрямую:

```ini
EnvironmentFile=/etc/tarot-app/app.env
```

Backend читает файл при каждом restart. Он не приходит из GitHub и не входит
в release.

В нём находятся:

- `APP_ENV`, `APP_HOST`, `APP_PORT`;
- `JWT_SECRET`, `JWT_TTL`;
- `MONGO_URI`, имя БД, auth source и Mongo timeouts;
- HTTP server timeouts;
- frontend/backend public URLs;
- OpenRouter key, model и referer; optional base URL и timeout могут быть
  добавлены при необходимости;
- gateway URL, callback secret и allowlist.

Значения нельзя выводить в CI logs, хранить в Git, добавлять в документацию
или передавать через аргументы процесса.

Посмотреть только имена переменных:

```bash
sudo awk -F= '/^[A-Za-z_][A-Za-z0-9_]*=/ {print $1}' \
  /etc/tarot-app/app.env
```

### 10.2. MongoDB env для backup

Путь:

```text
/etc/tarot-app/mongodump.yml
```

Права:

```text
root:root 0600
```

Файл содержит URI и password отдельного MongoDB-пользователя с ролью
`backup`. Его использует `/usr/local/sbin/backup-tarot`; backend его не
читает.

### 10.3. Frontend env

Frontend runtime env на VPS отсутствует. Vite использует при build:

```text
frontend/.env.production
```

`VITE_*` являются публичными build-time значениями и попадают в JavaScript.
Изменение требует нового commit/build/deploy. Файл `.env`, созданный рядом с
готовым `index.html` на VPS, не повлияет на frontend.

### 10.4. Env отдельного gtwapi

Путь существующего ресурса:

```text
/opt/gtwapi/.env
```

`gtwapi` запускается пользователем `gtwuser` из `/opt/gtwapi`. Этот env,
каталог и сервис не входят в Tarot deployment.

Фактический файл принадлежит `gtwuser:gtwuser` и сейчас имеет mode `0644`.
Поскольку `gtwapi` выведен из scope текущего deployment, права не менялись.
Ужесточение до `0600` следует выполнять отдельной согласованной операцией с
проверкой запуска `gtwapi`.

Shared callback secret можно менять только одновременно в Tarot backend и
`gtwapi`.

### 10.5. GitHub secrets не являются server env

`VPS_HOST`, `VPS_PORT`, `VPS_USER`, `VPS_SSH_KEY` и `VPS_HOST_KEY` хранятся
в GitHub Environment и используются только для доставки artifact.

## 11. Изменение production env

Env меняется отдельной административной процедурой:

1. создать root-only backup текущего файла;
2. подготовить новый файл без BOM и CRLF;
3. сохранить `root:tarot-user 0640`;
4. проверить обязательные keys без вывода значений;
5. выполнить `systemctl restart tarot-app`;
6. проверить local/public health и бизнес-функцию;
7. при ошибке восстановить env backup.

Metadata и health:

```bash
sudo stat -c '%n owner=%U:%G mode=%a' /etc/tarot-app/app.env
sudo systemctl restart tarot-app
curl --fail http://127.0.0.1:8888/api/health
curl --fail https://nebula-tarot.ru/api/health
```

Смена `JWT_SECRET` завершает существующие web-сессии. Смена `MONGO_URI`
должна быть согласована с MongoDB user/password. Callback secret меняется
одновременно с `gtwapi`.

## 12. Post-deploy проверка

```bash
readlink -f /opt/tarot-app/current
systemctl is-active tarot-app mongod nginx gtwapi
curl --fail http://127.0.0.1:8888/api/health
curl --fail https://nebula-tarot.ru/api/health
journalctl -u tarot-app --since "10 minutes ago"
```

Дополнительно вручную проверяются:

- главная страница и hashed asset;
- существующий upload;
- login/profile;
- основные reading endpoints;
- AI-интерпретация.

Automatic smoke test не вызывает платный AI provider. Доступность OpenRouter
или другого AI-провайдера контролируется отдельно.

## 13. Ручной rollback

```bash
readlink -f /opt/tarot-app/current
find /opt/tarot-app/releases -mindepth 1 -maxdepth 1 -type d -printf '%f\n'

previous=/opt/tarot-app/releases/<previous-sha>
temporary=/opt/tarot-app/.manual-rollback-$$
sudo ln -s "$previous" "$temporary"
sudo mv -Tf "$temporary" /opt/tarot-app/current
sudo systemctl restart tarot-app
curl --fail http://127.0.0.1:8888/api/health
curl --fail https://nebula-tarot.ru/api/health
```

Rollback кода не изменяет env, MongoDB, uploads, Nginx/TLS или `gtwapi`.
При изменениях схемы данных совместимость оценивается отдельно.

## 14. Диагностика failed deploy

```bash
systemctl status tarot-app
journalctl -u tarot-app --since "30 minutes ago"
journalctl -u tarot-health.service --since "30 minutes ago"
nginx -t
readlink -f /opt/tarot-app/current
find /var/lib/tarot-deploy/incoming -maxdepth 1 -type f -ls
```

| Симптом | Что проверять |
| --- | --- |
| checksum mismatch | соответствие archive и checksum |
| unsafe archive | лишние paths или links внутри artifact |
| backend не стартует | env, Mongo auth, architecture, journal |
| local health работает, public нет | Nginx, edge/TLS, DNS |
| frontend не обновился | `current`, browser/CDN cache, asset hash |
| API работает, AI возвращает 502 | доступ VPS к AI provider |

После failed deploy нужно подтвердить, что `current` вернулся на предыдущий
release и public health снова успешен.
