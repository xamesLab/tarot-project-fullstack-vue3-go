# Production deployment Nebula Tarot

**Статус:** действующая инструкция

**Проверено:** 21 июля 2026 года

## 1. Назначение

Документ описывает фактически настроенный ручной production-деплой после `git push`. Источник релиза — точный commit SHA; frontend и backend собираются локально, на VPS отправляется только готовый архив.

Обычный деплой не изменяет:

- `/etc/tarot-app/app.env`;
- MongoDB и `/etc/mongod.conf`;
- `/var/lib/tarot-app/uploads`;
- Nginx и TLS;
- `gtwapi`;
- OpenRouter relay-прокси и VPN на relay-сервере.

## 2. Текущая production-структура

| Ресурс | Фактическое расположение |
|---|---|
| Публичный сайт | `https://nebula-tarot.ru` |
| Nginx site | `/etc/nginx/sites-available/nebula-tarot.ru` |
| Releases | `/opt/tarot-app/releases/<commit-sha>` |
| Активный релиз | `/opt/tarot-app/current` — symlink на release |
| Backend binary | `/opt/tarot-app/current/backend/tarot-app` |
| Frontend root | `/opt/tarot-app/current/frontend` |
| Backend service | `tarot-app.service` |
| Runtime user | `tarot-user` |
| Production env | `/etc/tarot-app/app.env` |
| Persistent uploads | `/var/lib/tarot-app/uploads` |
| Incoming artifacts | `/var/lib/tarot-deploy/incoming` |
| Server installer | `/usr/local/sbin/install-tarot-release` |
| Smoke test | `/usr/local/sbin/smoke-test-tarot` |
| OpenRouter relay | `tarot-openrouter-relay.service` |

Внутри релиза находятся только готовый backend binary и статическая frontend-сборка:

```text
/opt/tarot-app/releases/<commit-sha>/
├── backend/
│   ├── tarot-app
│   └── uploads -> /var/lib/tarot-app/uploads
└── frontend/
    ├── index.html
    └── assets/
```

## 3. Деплой после push

Из корня репозитория в Git Bash или WSL:

```bash
bash deploy.sh
```

Скрипт [`deploy.sh`](../../deploy.sh):

1. разрешает production-деплой только из `main`;
2. останавливается при незакоммиченных изменениях;
3. выполняет `git fetch` и проверяет, что `HEAD` совпадает с upstream;
4. запускает `go test ./...`;
5. собирает статический Linux `amd64` backend binary;
6. выполняет `npm ci`, lint, frontend tests и production build;
7. формирует архив `backend/ + frontend/` и SHA-256 checksum;
8. загружает оба файла на VPS через обязательный нестандартный SSH-порт;
9. запускает серверный установщик;
10. удаляет локальный временный каталог при завершении.

По умолчанию скрипт использует SSH-ключ `~/.ssh/id_rsa`. Параметры можно переопределить без изменения файла:

```bash
DEPLOY_KEY=/path/to/key \
DEPLOY_HOST=production-host \
DEPLOY_PORT=ssh-port \
bash deploy.sh
```

Не добавляйте реальные адреса, пароли и приватные ключи в репозиторий.

## 4. Что делает серверный установщик

`/usr/local/sbin/install-tarot-release` принимает commit SHA, архив и checksum. Он:

1. проверяет формат 40-символьного commit SHA;
2. принимает файлы только из `/var/lib/tarot-deploy/incoming`;
3. блокирует параллельный деплой через `flock`;
4. сверяет SHA-256;
5. отклоняет архивы с небезопасными путями, symlink или неожиданными файлами;
6. распаковывает релиз во временный staging-каталог;
7. подключает persistent uploads;
8. атомарно переключает `/opt/tarot-app/current`;
9. перезапускает только `tarot-app.service`;
10. выполняет локальный и публичный smoke test;
11. автоматически возвращает предыдущий symlink и перезапускает backend при ошибке;
12. удаляет входной архив после успешной установки.

Сообщение об успешном завершении:

```text
smoke tests passed
release <commit-sha> installed successfully
```

## 5. Production env

Единственный runtime env приложения:

```text
/etc/tarot-app/app.env
```

Права:

```text
root:tarot-user 0640
```

Env не хранится в Git, не включается в release-архив и не заменяется при деплое. После ручного изменения env:

```bash
chown root:tarot-user /etc/tarot-app/app.env
chmod 640 /etc/tarot-app/app.env
systemctl restart tarot-app
systemctl status tarot-app --no-pager
```

Основные группы переменных:

- приложение и HTTP timeouts;
- MongoDB;
- JWT;
- frontend/public URLs;
- `OPENROUTER_*`;
- `HTTPS_PROXY` и `NO_PROXY` для OpenRouter relay;
- сохранённые `GATEWAY_*`, которые не делают `gtwapi` частью deploy flow.

Значения секретов нельзя выводить в логи или копировать в документацию.

## 6. Проверки после деплоя

Установщик выполняет их автоматически, но их можно повторить вручную:

```bash
systemctl status tarot-app --no-pager
systemctl status tarot-openrouter-relay --no-pager

curl --fail --silent http://127.0.0.1:8888/api/health
curl --fail --silent https://nebula-tarot.ru/api/health
curl --fail --silent https://nebula-tarot.ru/ >/dev/null

journalctl -u tarot-app -n 100 --no-pager
```

При изменениях AI дополнительно выполняется контролируемый запрос интерпретации. Такой smoke test платный и может создать обычную гостевую запись в MongoDB.

## 7. Ручной откат

Автоматический откат уже встроен в installer. Для ручного отката выберите предыдущий существующий release:

```bash
readlink -f /opt/tarot-app/current
ls -1 /opt/tarot-app/releases

ln -sfn /opt/tarot-app/releases/<previous-commit-sha> /opt/tarot-app/current
systemctl restart tarot-app
curl --fail https://nebula-tarot.ru/api/health
```

При откате приложения:

- env не откатывается;
- MongoDB автоматически не откатывается;
- uploads сохраняются;
- `gtwapi`, relay-прокси, VPN, Nginx и TLS не перезапускаются.

Если релиз требует несовместимого изменения данных, для него нужен отдельный migration/rollback plan и резервная копия базы.

## 8. GitHub Actions

Автоматического deployment workflow в репозитории пока нет. `git push` сам по себе production не обновляет; после push нужно запустить `bash deploy.sh`.

Будущая автоматизация должна переиспользовать тот же формат release-архива и серверный installer. Для CI следует использовать отдельного `tarot-deploy` с ограниченным `sudo`, а не root-доступ. Production environment рекомендуется защищать ручным подтверждением.

## 9. Связанные документы

- [OpenRouter relay-прокси](./OPENROUTER_PROXY.md)
- [Карта проекта](./AI_PROJECT_MAP.md)
- [Индекс документации](../README.md)
