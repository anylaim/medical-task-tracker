# Task Service

Сервис для управления задачами с HTTP API на Go с поддержкой периодичности (recurrence).

---

## Требования

* Go `1.23+`
* Docker и Docker Compose

---

## Быстрый запуск

```bash
docker compose up --build
```

Сервис будет доступен:

```
http://localhost:9000
```

---

### Важно

Если Postgres уже запускался ранее:

```bash
docker compose down -v
docker compose up --build
```

Миграции из `migrations/0001_create_tasks.up.sql` применяются только при первом создании volume.

---

## Swagger

* UI:
  http://localhost:9000/swagger/

* OpenAPI JSON:
  http://localhost:9000/swagger/openapi.json

---

## API

Базовый префикс:

```
/api/v1
```

### Основные эндпоинты

* `POST /api/v1/tasks`
* `GET /api/v1/tasks`
* `GET /api/v1/tasks/{id}`
* `PUT /api/v1/tasks/{id}`
* `DELETE /api/v1/tasks/{id}`

---

# Реализация периодичности

Добавлена поддержка периодических задач с автоматической генерацией серии.

## Поддерживаемые типы

### 1. Daily (каждые N дней)

```json
{
  "recurrence_type": "daily",
  "recurrence_value": 2
}
```

---

### 2. Monthly (ежемесячно)

```json
{
  "recurrence_type": "monthly",
  "due_date": "2025-01-10T00:00:00Z"
}
```

---

### 3. Specific Dates (конкретные даты)

```json
{
  "recurrence_type": "specific_dates",
  "specific_dates": [
    "2025-01-01T00:00:00Z",
    "2025-01-10T00:00:00Z"
  ]
}
```

---

### 4. Parity (чётные / нечётные дни)

```json
{
  "recurrence_type": "parity",
  "parity_type": "even"
}
```

---

## Как это работает

Используется стратегия **eager generation**:

* задачи создаются сразу при запросе
* генерируется серия задач вперёд

### Горизонт генерации

| Тип     | Горизонт   |
| ------- | ---------- |
| daily   | 30 дней    |
| monthly | 12 месяцев |
| parity  | 30 дней    |
| weekly  | 10 недель  |

---

## Архитектурные решения

* Чистая архитектура: handler → usecase → repository
* DTO отделены от domain-модели
* Бизнес-логика полностью находится в usecase
* Используется PostgreSQL + pgx

---

## Связь задач (серии)

Задачи, созданные в рамках одной периодичности:

* первая задача — родительская
* остальные содержат `parent_id`

Это позволяет:

* группировать задачи
* потенциально работать с сериями

---

## Валидация

Реализована базовая валидация:

* `daily` → требуется `recurrence_value > 0`
* `specific_dates` → список не пуст
* `parity` → только `even` или `odd`

---

## Ограничения (осознанные)

Для упрощения решения:

* recurrence хранится в каждой задаче (дублирование)
* нет отдельной сущности recurrence
* задачи генерируются только при создании
* нет обновления всей серии задач

---

## Возможные улучшения (production)

* выделить recurrence в отдельную таблицу
* генерировать задачи через background worker
* добавить редактирование серии
* добавить soft-delete

---

## Примеры запросов

### Создание ежедневной задачи

```bash
curl -X POST http://localhost:9000/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Обход пациентов",
    "description": "Ежедневный обход",
    "status": "new",
    "recurrence_type": "daily",
    "recurrence_value": 1
  }'
```

---

### Получение списка задач

```bash
curl http://localhost:9000/api/v1/tasks
```

---

### Получение задачи

```bash
curl http://localhost:9000/api/v1/tasks/1
```

---

## Тесты

```bash
go test ./...
```

---

## Итог

Реализована поддержка периодических задач с учётом различных сценариев использования, включая генерацию серий задач и базовую валидацию входных данных.

Решение ориентировано на простоту и расширяемость.
