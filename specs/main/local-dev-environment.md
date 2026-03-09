# Локальное окружение разработки (Step 3)

Дата проверки: 8 марта 2026.

## Цель шага
Подготовить и проверить локальное окружение для старта разработки frontend/backend и подключения к MongoDB по SSH-туннелю.

## Обязательные инструменты
- Go (для backend)
- Node.js + npm (для frontend)
- OpenSSH client (для туннеля к MongoDB на VPS)

## Опциональные инструменты
- Docker / Docker Compose (необязателен для текущего MVP flow)

## Факт проверки
- Go: установлен (`go1.26.1`), найден по пути `C:\Program Files\Go\bin\go.exe`.
- Node.js: установлен (`v18.16.0`).
- npm: установлен (`9.6.7`).
- SSH client: установлен (`OpenSSH_for_Windows_9.5p1`).
- Docker: не установлен (допустимо на текущем этапе).

## Автоматическая preflight-проверка
- Скрипт: `deploy/scripts/check-dev-env.ps1`
- Команда запуска:
  - `powershell -ExecutionPolicy Bypass -File deploy/scripts/check-dev-env.ps1`

## Результат шага 3
Локальная среда для разработки backend/frontend подготовлена на уровне инструментов.
Разрешено переходить к шагу 4 (конфигурация `.env` и `.env.example`).
