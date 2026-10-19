# Готово к независимому ревью

Дата: 2026-10-19

## Реализовано

- Dev-запуск поддерживает `DEV_DATA_MODE=auto|mongo|mock`. В `auto` при отсутствии `backend/.env` или недоступности MongoDB frontend работает с моками; `mock` запускает frontend с мок API и авторизацией. `stop-dev.sh` завершает дерево процессов.
- Frontend-моки и фикстуры вынесены в корневой каталог `mocks/`. Реализован мок API карты дня для гостя и пользователя.
- Блок «Послание дня» отделён от краткого `general`. `general` ограничен 20 словами; `message` и четыре сферы (`relationships`, `finances`, `emotional_state`, `career_business`) должны содержать не менее 50 слов каждая. Backend валидирует ограничения и запрашивает исправление ответа при нарушении; версия промпта обновлена для обхода старого кэша.
- Блок «Послание дня» расположен над сферами и увеличен с максимума 620 px до 930 px (+50%). Сферы занимают ту же ширину; на узком экране используют доступное место.
- Интерпретация автоматически раскрывается после появления данных и показа карты. Под кнопкой сворачивания/разворачивания задан отступ 30 px.
- Oracle chat отвечает на каждый непустой вопрос локализованным текстом: «Оракул сейчас молчит» на русском и «The Oracle is silent right now.» на английском. Ответ формируется на клиенте без обращения к API.
- Русские и английские мок-данные обновлены под новую схему. Тесты покрывают ограничения длины, раздельные `general`/`message`, автоматическое раскрытие и ответы Oracle на обоих языках.

## Исправление по результатам ревью

- Замечание P2 о русскоязычном ответе при английском интерфейсе устранено: добавлен ключ `oracleSilenceReply` в обе локали, компонент использует перевод, smoke-тест проверяет русский и английский ответы.

## Затронутые области

- Запуск и документация: `start-dev.sh`, `stop-dev.sh`, `README.md`, `mocks/README.md`, `frontend/IMPLEMENTATION_PLAN.md`.
- Backend-схема, промпт, валидация и тесты: `backend/internal/models/ask.go`, `backend/internal/services/ai_prompt.go`, `backend/internal/services/ask.go`, `backend/internal/services/ask_test.go`.
- Frontend-интерпретация и Oracle: `frontend/src/widgets/Interpretation/Interpretation.tsx`, `frontend/src/features/oracle-chat/OracleChat.tsx`, `frontend/src/pages/home/home.css`.
- API, типы, интерфейс, локализация и сборка: `frontend/src/api/`, `frontend/src/shared/`, `frontend/src/features/i18n/locales/{ru,en}/home.json`, `frontend/tsconfig.app.json`, `frontend/vite.config.js`, `frontend/vite.config.ts`.
- Моки и тесты: `mocks/frontend/`, `frontend/src/test/`; прежние моки из `frontend/src/api/mock/` перенесены в `mocks/frontend/`.

## Проверки

Команды ниже запускались из указанных каталогов.

- Из `backend/`: `go test ./...` — успешно; пакеты `internal/handlers`, `internal/services` и `internal/utils` прошли.
- Из `frontend/`: `$env:VITE_API_MODE='mock'; $env:VITE_AUTH_API_MODE='http'; $env:VITE_API_URL='http://localhost:8888/api'; npm.cmd test` — успешно, 15 файлов и 57 тестов; smoke-тест проверил ответ на русском и английском.
- Из `frontend/`: `npm.cmd run typecheck` — успешно.
- Из `frontend/`: `npm.cmd run lint` — успешно, ошибок и предупреждений нет.
- Из `frontend/`: `npm.cmd run build` — успешно, Vite преобразовал 194 модуля.
- Из `frontend/`: `npx.cmd eslint ../mocks/frontend --ext .ts --max-warnings 0 --config ./.eslintrc.cjs` — успешно.
- Из корня: `git diff --check` — успешно, замечаний о пробелах нет.
- Ранее выполнена проверка синтаксиса скриптов: `& 'C:\Program Files\Git\usr\bin\bash.exe' -n ./start-dev.sh ./stop-dev.sh` — успешно.
- Ранее вручную проверен dev fallback без `backend/.env` и с недоступной MongoDB (`MONGO_URI=mongodb://127.0.0.1:1/tarot_daily`): frontend переключился на моки, `Invoke-WebRequest -UseBasicParsing http://localhost:5173/` вернул HTTP 200. Режим `DEV_DATA_MODE=mock` также запускался отдельно; процессы остановлены через `stop-dev.sh`.

## Осталось учесть

- Подключение к рабочему экземпляру MongoDB не проверялось; fallback проверен при недоступном локальном порте.
- Mock-режим не запускает Go API и не сохраняет моковые данные в MongoDB.
- При предыдущем `npm ci` npm сообщил о 29 предупреждениях аудита зависимостей (5 moderate, 22 high, 2 critical); зависимости в рамках этой работы не обновлялись.
- Ответ-заглушка Oracle добавляется локально в состояние интерфейса и не сохраняется после перезагрузки страницы.

## Git

- Ветка: `Nebula-tarot`.
- В рабочей копии есть незакоммиченные изменения; `REVIEW_READY.md` включён в них.
- Коммит не создавался. Текущий `HEAD`: `b10487b682ea682b58b3be17e68961e5ee5a0afe`.
