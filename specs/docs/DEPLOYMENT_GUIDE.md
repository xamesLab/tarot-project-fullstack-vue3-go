# Рекомендации по production-деплою

## 1. Назначение документа

Этот документ описывает организацию production-деплоя проекта Nebula Tarot из GitHub на VPS и служит эксплуатационной инструкцией.

Основные принципы:

- `gtwapi` не входит в процесс деплоя и продолжает работать без изменений;
- frontend и backend собираются в CI, а на сервер доставляются готовые артефакты;
- production-секреты и пользовательские файлы не хранятся в Git;
- каждый релиз размещается в отдельном каталоге и может быть быстро откачен;
- переключение версии выполняется через символьную ссылку `current`;
- MongoDB, env-файл и uploads не заменяются при обычном деплое.

Базовая схема применена 21 июля 2026 года. Активный application-релиз:
`2ffc919215b9dcceedd57651b150e9514fc7db94`. Deployment automation и
операционные файлы находятся в ветке `codex/deployment-automation`.

## 2. Текущее состояние production

На VPS используются следующие ресурсы:

| Ресурс | Расположение или адрес |
| --- | --- |
| Домен | `nebula-tarot.ru` |
| Nginx-конфигурация | `/etc/nginx/sites-available/nebula-tarot.ru` |
| Активный релиз | `/opt/tarot-app/current` |
| Текущий frontend | `/opt/tarot-app/current/frontend` |
| Текущий backend | `/opt/tarot-app/current/backend` |
| Backend binary | `/opt/tarot-app/current/backend/tarot-app` |
| Backend service | `/etc/systemd/system/tarot-app.service` |
| Backend user | `tarot-user` |
| Backend env | `/etc/tarot-app/app.env` |
| Backend endpoint | `127.0.0.1:8888` |
| Постоянные uploads | `/var/lib/tarot-app/uploads` |
| MongoDB service | `mongod.service` |
| MongoDB configuration | `/etc/mongod.conf` |
| MongoDB endpoint | `127.0.0.1:27017`, authorization enabled |
| Backup credentials | `/etc/tarot-app/mongodump.yml`, `root:root 0600` |
| Deploy user | `tarot-deploy` |
| Release installer | `/usr/local/sbin/install-tarot-release` |
| Gateway | `/opt/gtwapi` |
| Gateway service | `gtwapi.service` |
| Gateway port | `8080` |

Nginx отдаёт frontend как статические файлы и проксирует:

- `/api/` на `http://127.0.0.1:8888/api/`;
- `/uploads/` на `http://127.0.0.1:8888/uploads/`.

`gtwapi.service`, каталог `/opt/gtwapi` и его env-файл не должны изменяться при внедрении нового процесса деплоя Tarot-приложения.
При внедрении они не изменялись; внешний TLS также не изменялся.

## 3. Рекомендуемая структура сервера

```text
/opt/tarot-app/
├── releases/
│   ├── <commit-sha>/
│   │   ├── backend/
│   │   │   ├── tarot-app
│   │   │   └── uploads -> /var/lib/tarot-app/uploads
│   │   └── frontend/
│   │       ├── index.html
│   │       └── assets/
│   └── ...
└── current -> releases/<commit-sha>

/etc/tarot-app/
├── app.env
└── mongodump.yml

/var/lib/tarot-app/
└── uploads/

/var/backups/tarot-app/
├── pre-deploy-<timestamp>/
└── <YYYYMMDD-HHMMSS>/
    ├── mongodb-all.archive.gz
    ├── uploads.tar.gz
    ├── SHA256SUMS
    └── MANIFEST
```

Назначение каталогов:

- `/opt/tarot-app/releases` — неизменяемые версии backend и frontend;
- `/opt/tarot-app/current` — ссылка на активный релиз;
- `/etc/tarot-app/app.env` — production-конфигурация и секреты;
- `/var/lib/tarot-app/uploads` — постоянные пользовательские файлы;
- `/var/backups/tarot-app` — резервные копии MongoDB и uploads.

Frontend и backend рекомендуется включать в один каталог релиза. Это позволяет одной операцией переключать согласованные версии API и клиентского приложения.

## 4. Подготовка репозитория

### 4.1. Исключение runtime-данных

Собранный Go-бинарник и пользовательские загрузки не должны храниться в Git.

Рекомендуемые правила `.gitignore`:

```gitignore
backend/tarot-app
backend/uploads/*
!backend/uploads/.gitkeep
```

Текущие файлы в `backend/uploads` следует разделить на две категории:

- пользовательские и изменяемые файлы — хранить только в `/var/lib/tarot-app/uploads`;
- обязательные начальные ресурсы — перенести, например, в `backend/deploy/seed/assets` и устанавливать отдельной seed-командой.

### 4.2. Проверка секретов

Следует проверить содержимое и историю файлов:

```text
specs/temp/env.txt
specs/temp/env (1).txt
```

Эти файлы удалены из текущего дерева Git. Если в них когда-либо находились реальные пароли, токены или ключи, необходимо:

1. заменить все раскрытые секреты;
2. удалить секретные файлы из актуальной версии репозитория;
3. при необходимости очистить историю Git;
4. включить автоматическое сканирование репозитория на секреты.

Созданы новые отдельные credentials для плановых MongoDB-backup. Application
MongoDB credentials и JWT secret ротированы 21 июля 2026 года; старый MongoDB
application user удалён, а существующие web-сессии после смены JWT требуют
повторного входа. Общий callback secret нельзя менять без согласованного
обновления `gtwapi`; в рамках внедрения он намеренно не изменялся.

Production-файл `/etc/tarot-app/app.env` не должен копироваться в репозиторий или включаться в release-архив.

### 4.3. Рекомендуемые deployment-файлы

```text
.github/workflows/ci.yml
.github/workflows/deploy-production.yml
deploy/install-release.sh
deploy/smoke-test.sh
deploy/tarot-app.service
deploy/nginx.conf
```

Серверный `install-release.sh` рекомендуется хранить в репозитории, но устанавливать на VPS как root-owned файл, недоступный для изменения deploy-пользователем.

## 5. CI для pull request

Слияние в `main` следует разрешать только после успешного прохождения автоматических проверок.

### 5.1. Backend

```bash
cd backend
go test ./...
go vet ./...
go build ./...
```

Версию Go необходимо брать из `backend/go.mod`. На момент подготовки документа проект использует Go `1.23.1`.

### 5.2. Frontend

```bash
cd frontend
npm ci
npm run typecheck
npm run lint
npm test
npm run build
```

Версию Node.js рекомендуется зафиксировать в `.nvmrc` и использовать одинаковую версию локально и в GitHub Actions.

Каталог `frontend-legacy` не следует включать в production pipeline. Актуальный frontend расположен в `frontend` и использует React, TypeScript и Vite.

## 6. Production-конфигурация frontend

Текущий `frontend/.env.production` содержит:

```env
VITE_API_MODE=mock
VITE_AUTH_API_MODE=http
VITE_API_URL=/api
```

Это означает:

- авторизация обращается к production backend;
- часть контентного API использует mock-реализацию;
- относительный `/api` направляется Nginx на backend этого же домена.

Нельзя переключать `VITE_API_MODE=http`, пока HTTP-адаптер в `frontend/src/api/http.ts` не будет полностью реализован и проверен.

Все переменные `VITE_*` встраиваются в JavaScript bundle и доступны пользователю браузера. В них запрещено помещать пароли, API-ключи и другие секреты.

## 7. Сборка release-артефакта

### 7.1. Backend

Backend рекомендуется собирать в GitHub Actions для Linux `amd64`:

```bash
mkdir -p artifact/backend
cd backend

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -trimpath -ldflags="-s -w" \
-o ../artifact/backend/tarot-app .
```

### 7.2. Frontend

```bash
cd frontend
npm ci
npm run build

mkdir -p ../artifact/frontend
cp -R dist/. ../artifact/frontend/
```

### 7.3. Архив релиза

Результат сборки:

```text
artifact/
├── backend/
│   └── tarot-app
└── frontend/
    ├── index.html
    └── assets/
```

Архив следует именовать по commit SHA или release tag:

```text
tarot-app-<commit-sha>.tar.gz
tarot-app-<commit-sha>.tar.gz.sha256
```

Checksum должен проверяться на VPS до распаковки.

## 8. Доступ GitHub Actions к VPS

Для деплоя рекомендуется создать отдельного пользователя `tarot-deploy`. Backend при этом продолжает работать как `tarot-user`.

Не рекомендуется:

- подключать GitHub Actions как `root`;
- использовать личный SSH-ключ администратора;
- отключать проверку SSH host key;
- давать deploy-пользователю неограниченный `sudo`.

В GitHub Environment `production` следует хранить:

```text
VPS_HOST=d51673dc67a4.vps.myjino.ru
VPS_PORT=49401
VPS_USER=tarot-deploy
VPS_SSH_KEY=<private deploy key>
VPS_HOST_KEY=<pinned server host key>
```

GitHub Environment желательно защитить ручным подтверждением production-деплоя.

Deploy-пользователь должен иметь право:

- загружать новый release-архив во временный каталог;
- запускать один root-owned deployment script через ограниченное правило `sudo`;
- читать результат health check.

Он не должен получать произвольный root shell.

SSH принимает только public-key authentication. Для root действует
`PermitRootLogin prohibit-password`, а `tarot-user` исключён из группы
`sudo`. Ограниченный `tarot-deploy` может через sudo запускать только
`/usr/local/sbin/install-tarot-release *`. Применяемый drop-in хранится в
`deploy/sshd-hardening.conf`.

## 9. Последовательность деплоя

Рекомендуемый production workflow:

1. Выполнить backend и frontend проверки.
2. Собрать release-артефакт из точного commit SHA.
3. Сформировать архив и checksum.
4. Загрузить их на VPS по SSH через порт `49401`.
5. Получить блокировку деплоя через `flock`.
6. Проверить checksum.
7. Распаковать релиз в `/opt/tarot-app/releases/<commit-sha>`.
8. Проверить структуру и права файлов.
9. Создать ссылку `backend/uploads` на `/var/lib/tarot-app/uploads`.
10. Запомнить предыдущую цель ссылки `current`.
11. Переключить `current` на новый релиз.
12. Перезапустить `tarot-app.service`.
13. Выполнить локальный и внешний health check.
14. При ошибке вернуть предыдущий релиз и повторно запустить сервис.
15. После успеха удалить временный архив.
16. Оставить несколько последних релизов для отката.

Пример проверок:

```bash
curl --fail --silent http://127.0.0.1:8888/api/health
curl --fail --silent https://nebula-tarot.ru/api/health
```

Серверный deployment script должен завершаться с ошибкой при любом неуспешном шаге и не должен удалять предыдущий рабочий релиз до завершения smoke test.

## 10. Systemd

После перехода на release-структуру рекомендуемая основа сервиса:

```ini
[Unit]
Description=TarotApp Go Backend
After=network.target mongod.service

[Service]
Type=simple
User=tarot-user
Group=tarot-user
WorkingDirectory=/opt/tarot-app/current/backend
EnvironmentFile=/etc/tarot-app/app.env
ExecStart=/opt/tarot-app/current/backend/tarot-app
Restart=always
RestartSec=5

NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/tarot-app/uploads

[Install]
WantedBy=multi-user.target
```

Рекомендуемые права:

```text
/etc/tarot-app/app.env             root:tarot-user 0640
/opt/tarot-app/releases            root или tarot-deploy, read-only для tarot-user
/var/lib/tarot-app/uploads         tarot-user:tarot-user
```

Пользователь приложения должен иметь возможность изменять только runtime-данные, но не бинарники и env-файл.

## 11. Nginx

После перехода на release-структуру frontend root следует направить на активный релиз:

```nginx
root /opt/tarot-app/current/frontend;
index index.html;

location /api/ {
    proxy_pass http://127.0.0.1:8888/api/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

location /uploads/ {
    proxy_pass http://127.0.0.1:8888/uploads/;
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

location / {
    try_files $uri $uri/ /index.html;
}
```

Дополнительно следует:

- настроить перенаправление HTTP на HTTPS;
- проверить SSL-конфигурацию именно для `nebula-tarot.ru`;
- выполнять `nginx -t` перед reload при изменении конфигурации;
- настроить `client_max_body_size` с учётом загрузки ZIP-архивов колод;
- закрыть внешний доступ к backend-порту `8888`;
- оставить backend доступным Nginx через loopback-интерфейс.

Изменение Nginx не требуется при каждом обычном релизе. Конфигурация должна обновляться отдельной административной процедурой.

## 12. MongoDB и seed-команды

Обычный деплой не должен автоматически:

- удалять базу;
- очищать коллекции;
- повторно выполнять seed;
- менять пользователей MongoDB;
- изменять `/etc/mongod.conf`.

Seed следует запускать отдельной ручной процедурой только для релизов, которым он необходим. Перед операцией нужно создавать резервную копию MongoDB.

Аудит активных соединений, systemd-сервисов, Docker и конфигурационных файлов
не выявил внешних потребителей MongoDB. MongoDB переведена на
`127.0.0.1:27017`, а `security.authorization` включена. Backend использует
отдельного пользователя с ролью `readWrite` только для `tarot_daily`;
backup использует отдельного пользователя с ролью `backup`.

Историческая база `servitor_trade` сохранена в полном backup и не удалялась.
Если для неё появится отдельный потребитель, ему потребуется собственный
минимально привилегированный MongoDB-пользователь и локальный доступ.

## 13. Резервные копии

Минимальная схема резервного копирования:

- ежедневный `mongodump`;
- ежедневный или инкрементальный backup `/var/lib/tarot-app/uploads`;
- хранение нескольких поколений копий;
- периодическое копирование backup за пределы VPS;
- регулярная тестовая процедура восстановления.

Production env можно хранить в отдельном защищённом backup, но он не должен попадать в общий архив приложения или GitHub artifact.

На VPS ежедневная процедура реализуется скриптом `/usr/local/sbin/backup-tarot`
и таймером `tarot-backup.timer`. Она создаёт архив всех баз MongoDB и архив
uploads, проверяет SHA-256, запускает `mongorestore --dryRun` и только после
успеха атомарно публикует каталог с backup. Хранятся семь ежедневных поколений;
именованные каталоги `pre-deploy-*` автоматической ротацией не удаляются.
Для подключения используется отдельная MongoDB-учётная запись с ролью
`backup`; URI и пароль хранятся в root-only
`/etc/tarot-app/mongodump.yml` и не передаются через командную строку.

Проверка и ручной запуск:

```bash
systemctl start tarot-backup.service
systemctl status tarot-backup.service
systemctl list-timers tarot-backup.timer
journalctl -u tarot-backup.service
```

## 14. Откат

Откат выполняется переключением `current` на предыдущий релиз:

```bash
ln -sfn /opt/tarot-app/releases/<previous-sha> /opt/tarot-app/current
systemctl restart tarot-app
curl --fail https://nebula-tarot.ru/api/health
```

Рекомендуется хранить последние 5–10 успешно установленных релизов.

При откате:

- `/etc/tarot-app/app.env` не заменяется;
- `/var/lib/tarot-app/uploads` не заменяется;
- MongoDB автоматически не откатывается;
- `gtwapi` не перезапускается и не изменяется.

Если новый релиз требует несовместимого изменения данных, для него должен существовать отдельный план миграции и обратной совместимости.

## 15. Smoke tests после деплоя

Минимальный автоматический набор:

1. `systemctl is-active tarot-app` возвращает `active`.
2. `http://127.0.0.1:8888/api/health` отвечает успешно.
3. `https://nebula-tarot.ru/api/health` отвечает успешно через Nginx.
4. Главная страница возвращает HTML.
5. Один статический asset frontend доступен.
6. Один существующий `/uploads/...` файл доступен.

Расширенный ручной набор:

- регистрация и вход;
- восстановление сессии после обновления страницы;
- профиль пользователя;
- создание расклада;
- запрос AI-интерпретации;
- административный вход;
- загрузка колоды и рубашки;
- смена языка и темы.

## 16. Наблюдаемость

Backend уже пишет логи в journald через systemd. Основные команды диагностики:

```bash
journalctl -u tarot-app --since "30 minutes ago"
systemctl status tarot-app
systemctl status mongod
nginx -t
```

Рекомендуется добавить внешний uptime-monitoring для:

```text
https://nebula-tarot.ru/
https://nebula-tarot.ru/api/health
```

Локальный таймер `tarot-health.timer` каждые пять минут проверяет сервисы,
локальный и публичный health endpoint и наличие минимум 512 MiB свободного
места. Результаты доступны через `journalctl -u tarot-health.service`. Этот
таймер обнаруживает сбой, но без внешней системы оповещений сам уведомление не
отправляет.

Полезно также отслеживать:

- свободное место на диске;
- размер MongoDB;
- размер `/var/lib/tarot-app/uploads`;
- частоту рестартов `tarot-app.service`;
- HTTP 5xx в Nginx;
- срок действия TLS-сертификата.

## 17. Рекомендуемый порядок внедрения

1. Очистить репозиторий от бинарника, runtime uploads и возможных секретов.
2. Добавить CI-проверки pull request.
3. Зафиксировать версии Node.js и Go.
4. Реализовать сборку единого release-архива.
5. Подготовить deploy-пользователя и каталоги на VPS.
6. Перенести uploads в постоянный каталог с предварительным backup.
7. Установить root-owned deployment script.
8. Перевести systemd и Nginx на ссылку `current`.
9. Настроить GitHub Environment и SSH deploy key.
10. Выполнить первый ручной deployment workflow.
11. Проверить автоматический rollback на тестовом нерабочем релизе.
12. Включить production deployment по tag или после ручного подтверждения.
13. Настроить backup и uptime-monitoring.

Production-деплой лучше запускать по release tag или вручную из GitHub Actions. Автоматический деплой каждого commit в `main` допустим только после стабилизации CI, smoke tests и rollback-процедуры.
