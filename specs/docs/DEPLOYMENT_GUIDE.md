# Production deployment Nebula Tarot

**Статус:** действующая операционная инструкция
**Проверено:** 19 октября 2026 года

## 1. Схема деплоя

Источник истины — точный commit ветки origin/main на GitHub. Деплой запускается на VPS от пользователя tarot-deploy. Сервер получает commit, выполняет проверки и сборку, затем запускает существующий release installer.

Локальная машина не собирает и не загружает файлы релиза. Терминал нужен только для подключения к VPS; те же команды можно выполнить через консоль VPS у хостера.

Обычный деплой не заменяет:

- /etc/tarot-app/app.env;
- MongoDB и /etc/mongod.conf;
- /var/lib/tarot-app/uploads;
- Nginx и TLS;
- gtwapi;
- OpenRouter relay и VPN.

## 2. Production-структура

| Ресурс | Расположение |
|---|---|
| Публичный сайт | https://nebula-tarot.ru |
| Пользователь деплоя | tarot-deploy |
| Клон репозитория | /home/tarot-deploy/nebula-tarot |
| Активный релиз | symlink /opt/tarot-app/current |
| Каталоги релизов | /opt/tarot-app/releases/<commit-sha> |
| Backend service | tarot-app.service, запускается от tarot-user |
| Runtime env | /etc/tarot-app/app.env |
| Постоянные uploads | /var/lib/tarot-app/uploads |
| Входящие релизы | /var/lib/tarot-deploy/incoming |
| Привилегированный installer | /usr/local/sbin/install-tarot-release |
| Smoke tests | /usr/local/sbin/smoke-test-tarot |

На VPS требуются Linux x86_64, Git, Go и Node.js 20.19+ или 22.12+ с npm. Скрипт проверяет версию Node до установки зависимостей. У tarot-deploy настроен ограниченный sudo-доступ к installer. Репозиторий приватный, поэтому до первого клонирования нужно настроить read-only SSH-доступ к GitHub.

## 3. Первоначальная настройка GitHub и клона

На проверенном VPS уже создан отдельный Ed25519 deploy key пользователя tarot-deploy:

- Публичный ключ расположен в `/home/tarot-deploy/.ssh/github_deploy_key.pub`; вывести его можно командой `cat ~/.ssh/github_deploy_key.pub` от tarot-deploy.
- Fingerprint ключа: `SHA256:ZJoH/oWRkErc2aPiO82B7sfrLu/SFZgzY3J1v5f1yck`.
- Приватный ключ остаётся на VPS в `/home/tarot-deploy/.ssh/github_deploy_key`, с правами `0600`.

Добавьте содержимое `.pub` в GitHub: Settings репозитория → Deploy keys, с правом только на чтение. Не загружайте приватный ключ.

SSH config уже указывает на этот ключ и включает `StrictHostKeyChecking`. GitHub host key закреплён в `/home/tarot-deploy/.ssh/known_hosts`; его Ed25519 fingerprint — `SHA256:+DiY3wvvV6TuJJhbpZisF/zLDA0zPMSvHdkr4UvCOqU`, сверенный с [официальным списком fingerprints GitHub](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/githubs-ssh-key-fingerprints). Не отключайте проверку host key.

Проверьте доступ к репозиторию:

    git ls-remote git@github.com:xamesLab/tarot-project-fullstack-vue3-go.git refs/heads/main

Один раз клонируйте репозиторий:

    git clone --branch main --single-branch git@github.com:xamesLab/tarot-project-fullstack-vue3-go.git ~/nebula-tarot

Deploy key и проверка host key на VPS настроены; пока публичный deploy key не добавлен в GitHub, доступ к приватному репозиторию и клонирование не сработают. После добавления ключа выполните `git ls-remote`, затем клонируйте репозиторий.

## 4. Деплой

После merge изменений в GitHub main войдите на VPS и выполните:

    sudo -iu tarot-deploy
    cd ~/nebula-tarot
    bash deploy.sh

Скрипт:

1. Проверяет Linux x86_64, пользователя tarot-deploy, ожидаемый GitHub origin и чистое состояние ветки main.
2. Блокирует параллельный деплой, получает origin/main и переключает отдельный клон на точный SHA этого commit.
3. Запускает `CGO_ENABLED=0 go test ./...` и собирает backend для Linux amd64 также без CGO. Системные C-заголовки для этих шагов не нужны.
4. Выполняет `npm ci --include=dev`, frontend lint и `npm run test:deploy`, затем production build. Перед tests и build скрипт очищает унаследованные переменные VITE_*. Vitest использует frontend/.env.test, production bundle — frontend/.env.production (оба API-адаптера в production работают через HTTP). Команды явно выбирают `vite.config.ts`. Весь набор тестов выполняется одним worker без параллельных файлов, с лимитом 30 секунд на тест и hook. Отдельные асинхронные DOM-запросы ожидаются до 5 секунд; изоляция файлов сохраняется.
5. Формирует архив релиза и SHA-256 checksum на VPS и помещает их в закрытый каталог incoming.
6. Запускает root-owned installer через существующее ограниченное правило sudo.

Работающий service не перезапускается до завершения проверок и сборки. Installer проверяет checksum и состав архива, атомарно переключает symlink current, перезапускает только tarot-app.service и возвращает предыдущий релиз при неуспешных smoke tests. Артефакты сборки и зависимости остаются в deployment-клоне или временном каталоге; production secrets не используются при сборке frontend.

Перед сборкой скрипт останавливается, если в frontend найдены локальные dotenv-файлы, которые могли бы переопределить committed test или production конфигурацию.

Сборка на сервере использует CPU и место на диске, пока текущая версия продолжает работать. Перед крупными обновлениями зависимостей проверяйте свободное место. Скрипт печатает версию Node, доступную память, swap, load average и свободное место перед сборкой и перед frontend-тестами. Для диагностики используются стандартные Linux-команды `df` и `tee`, данные памяти и нагрузки читаются из `/proc`.

Подробный журнал frontend-тестов (длительность каждого теста и JS heap, который не равен общей памяти процесса) вместе с ресурсами перед прогоном сохраняется в `$HOME/.cache/tarot-deploy-tests-<SHA>.<случайный суффикс>.log`; точный путь выводится перед запуском. Журнал остаётся после ошибки и удаления временной сборки, не загрязняя Git checkout. Ошибка тестов или записи журнала останавливает деплой; повторы и пропуск тестов не используются. При необходимости старые журналы удаляются вручную. Обычный локальный `npm test` сохраняет лимит 10 секунд и максимум два worker.

## 5. Проверка после деплоя

При успехе выводятся сообщения:

    smoke tests passed
    release <commit-sha> installed successfully
    Deployment complete: <commit-sha>

Installer повторяет локальную проверку API до 30 секунд, пока service запускается. В это время может появляться connection refused; деплой считается успешным только после успешного завершения всех smoke tests.

Ручная проверка:

    systemctl status tarot-app --no-pager
    readlink -f /opt/tarot-app/current
    curl --fail --silent --show-error http://127.0.0.1:8888/api/health
    curl --fail --silent --show-error https://nebula-tarot.ru/api/health
    curl --fail --silent --show-error https://nebula-tarot.ru/ >/dev/null

Если релиз меняет поведение AI, выполните по одному контролируемому запросу на каждом затронутом языке. Генерация может тарифицироваться провайдером.

## 6. Откат

Installer автоматически возвращает предыдущий symlink и перезапускает backend, если health checks не проходят. Для ручного отката приложения администратор сервера выбирает сохранённый релиз:

    readlink -f /opt/tarot-app/current
    ls -1 /opt/tarot-app/releases
    ln -sfn /opt/tarot-app/releases/<previous-commit-sha> /opt/tarot-app/current
    systemctl restart tarot-app
    curl --fail https://nebula-tarot.ru/api/health

Откат приложения не меняет MongoDB, runtime env и uploads. Для несовместимого изменения данных нужен отдельный план миграции, backup и отката.

## 7. Безопасность и автоматизация

- GitHub deploy key на VPS доступен только для чтения и принадлежит пользователю tarot-deploy.
- Проверка host key GitHub включена.
- Пользователь деплоя собирает релиз и имеет ограниченный sudo-доступ к installer; само приложение продолжает работать от tarot-user.
- Не храните production secrets в клоне, release archive, shell history или логах.
- GitHub Actions deployment не настроен. git push не разворачивает релиз автоматически; оператор запускает скрипт на VPS.
