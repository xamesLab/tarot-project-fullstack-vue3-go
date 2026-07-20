# OpenRouter relay-прокси

**Статус:** настроено в production

**Проверено:** 21 июля 2026 года

## 1. Зачем нужен relay

Исходящие обращения основного VPS к OpenRouter блокировались внешней security policy до обработки API-запроса. Ключ, модель и URL OpenRouter при этом были корректными. Для смены точки исходящего трафика настроен отдельный relay-сервер.

Схема:

```text
tarot-app
  -> socks5h://127.0.0.1:1080
  -> постоянный SSH dynamic forward
  -> relay-сервер
  -> openrouter.ai:443
```

Relay не является публичным HTTP/SOCKS proxy. На основном VPS порт доступен только через loopback; на relay-сервере не открывался отдельный proxy-порт.

## 2. Компоненты

На основном VPS:

| Компонент | Назначение |
|---|---|
| `tarot-openrouter-relay.service` | Поддерживает SSH dynamic forward |
| `127.0.0.1:1080` | Локальная SOCKS5-точка |
| `tarot-relay` | Непривилегированный локальный пользователь сервиса |
| `/etc/tarot-app/relay/id_ed25519` | Отдельный приватный ключ relay-сервиса |
| `/etc/tarot-app/relay/known_hosts` | Закреплённый host key relay-сервера |
| `/etc/tarot-app/app.env` | Подключает backend к локальному SOCKS5 |

На relay-сервере:

- отдельный пользователь `openrouter-relay`;
- пароль пользователя заблокирован;
- интерактивная shell-сессия не используется;
- ключ ограничен destination `openrouter.ai:443`;
- agent forwarding, X11, PTY и user rc отключены;
- существующий SSH-порт используется без открытия публичного proxy-порта.

Адрес relay-сервера, пароли, API-ключи, private key и полная строка `authorized_keys` в документации намеренно не фиксируются.

## 3. Backend env

В `/etc/tarot-app/app.env` настроено:

```env
HTTPS_PROXY=socks5h://127.0.0.1:1080
NO_PROXY=127.0.0.1,localhost,::1
```

Используется `socks5h`, поэтому имя `openrouter.ai` разрешается на стороне relay. Backend создаёт стандартный Go `http.Client`, а стандартный transport учитывает `HTTPS_PROXY`.

Переменная действует на все исходящие HTTPS-запросы процесса `tarot-app`. В текущем коде реальный внешний HTTPS provider call — OpenRouter; `GATEWAY_*` сохраняются как конфигурационный stub и не запускают отдельный HTTP flow.

`OPENROUTER_BASE_URL` по умолчанию остаётся `https://openrouter.ai/api/v1`. Название модели не меняет provider base URL: модель передаётся в JSON тела `POST /chat/completions`.

## 4. Systemd и восстановление соединения

Unit расположен здесь:

```text
/etc/systemd/system/tarot-openrouter-relay.service
```

Сервис:

- запускается как `tarot-relay`;
- слушает только `127.0.0.1:1080`;
- использует pinned host key и только выделенный identity file;
- завершается, если forwarding не удалось открыть;
- отправляет SSH keepalive;
- автоматически перезапускается через несколько секунд;
- включён в `multi-user.target`.

Backend и relay — разные сервисы. Обычный deployment перезапускает только `tarot-app`; туннель продолжает работать.

## 5. Проверки

Статус и listener:

```bash
systemctl status tarot-openrouter-relay --no-pager
systemctl is-enabled tarot-openrouter-relay
ss -lntp '( sport = :1080 )'
```

Бесплатная проверка OpenRouter без model generation:

```bash
curl --proxy socks5h://127.0.0.1:1080 \
  --fail --silent --show-error \
  https://openrouter.ai/api/v1/models >/dev/null
```

Проверка backend:

```bash
systemctl status tarot-app --no-pager
curl --fail --silent http://127.0.0.1:8888/api/health
curl --fail --silent https://nebula-tarot.ru/api/health
journalctl -u tarot-app -n 100 --no-pager
```

Ограничение ключа проверялось отдельным запросом: OpenRouter разрешён, произвольный внешний домен через этот relay блокируется.

## 6. Задержка

Контрольные замеры на момент настройки:

- новый TCP/TLS connection: примерно `+0,5–0,6 с`;
- reused keep-alive connection: примерно `+0,1–0,25 с`;
- основной вклад в пользовательское время ответа обычно вносит генерация модели.

Значения являются ориентировочными и зависят от маршрута и нагрузки обоих серверов.

## 7. VPN на relay-сервере

На relay-сервере уже работает отдельный VPN. При настройке прокси не изменялись:

- VPN interface и его конфигурация;
- UDP listener VPN;
- IP routes и policy rules;
- forwarding/NAT rules;
- firewall policy.

Контроль после настройки подтвердил активный VPN service, свежие handshakes и рост счётчиков трафика. При обслуживании relay нельзя менять VPN, маршрутизацию или NAT ради OpenRouter-туннеля: для relay достаточно существующего SSH и обычного исходящего HTTPS.

## 8. Изменение и откат

После изменения `/etc/tarot-app/app.env`:

```bash
chown root:tarot-user /etc/tarot-app/app.env
chmod 640 /etc/tarot-app/app.env
systemctl restart tarot-app
```

Чтобы временно отключить relay для backend:

1. удалить `HTTPS_PROXY` и связанный `NO_PROXY` из production env либо восстановить проверенную backup-копию;
2. перезапустить `tarot-app`;
3. при необходимости остановить `tarot-openrouter-relay.service`.

Отключение вернёт исходящий запрос на основной VPS, где OpenRouter может снова блокироваться.

## 9. Правила безопасности

- Не публиковать IP relay-сервера без операционной необходимости.
- Не хранить SSH/API credentials в Git, документации и shell history.
- Не превращать relay в публичный SOCKS/HTTP proxy.
- Не использовать root для постоянного tunnel connection.
- Не снимать ограничение `openrouter.ai:443` без отдельного security review.
- После раскрытия пароля в любом канале менять его и переходить на key-only access.
- При смене host key сначала проверять fingerprint по доверенному каналу, затем обновлять pinned `known_hosts`.

## 10. Связанные документы

- [Production deployment](./DEPLOYMENT_GUIDE.md)
- [Карта проекта](./AI_PROJECT_MAP.md)
