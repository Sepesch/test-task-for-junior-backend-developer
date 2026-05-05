# Task Service

Сервис для управления задачами с HTTP API на Go.

## Требования

- Go `1.23+`
- Docker и Docker Compose

## Быстрый запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу `http://localhost:8080`.

Если `postgres` уже запускался ранее со старой схемой, пересоздай volume:

```bash
docker compose down -v
docker compose up --build
```

Причина в том, что SQL-файл из `migrations/0001_create_tasks.up.sql` монтируется в `docker-entrypoint-initdb.d` и применяется только при инициализации пустого data volume.

## Swagger

Swagger UI:

```text
http://localhost:8080/swagger/
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/openapi.json
```

## API

Базовый префикс API:

```text
/api/v1
```

Основные маршруты:

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`

# Task Service

## Recurring Tasks

You can attach recurrence rules to existing tasks. When a recurring task is marked as `done`, the scheduler will automatically reset it to `new` on the next occurrence date.

### API Endpoints for Recurrence

- `GET /api/v1/tasks/{id}/recurrence` – get recurrence settings
- `PUT /api/v1/tasks/{id}/recurrence` – set or replace recurrence
- `DELETE /api/v1/tasks/{id}/recurrence` – remove recurrence

### Supported recurrence types

| type          | fields                        | description                        |
|---------------|-------------------------------|------------------------------------|
| daily         | `interval` (int >0)           | every N days                       |
| monthly       | `day_of_month` (1-30)         | each month on the given day        |
| weekly        | `week_days` (1..7, 1=Monday)  | on selected weekdays               |
| specific_dates| `dates` (array of YYYY-MM-DD) | only on exact dates                |
| even_odd      | `parity` ("even"/"odd")       | on even or odd days of the month   |

### Example requests

```bash
# Set daily recurrence every 2 days
curl -X PUT /api/v1/tasks/1/recurrence -H "Content-Type: application/json" \
  -d '{"type":"daily","interval":2}'

# Set monthly on 15th
curl -X PUT /api/v1/tasks/1/recurrence -d '{"type":"monthly","day_of_month":15}'

# Set weekly on Mon/Wed/Fri
curl -X PUT /api/v1/tasks/1/recurrence -d '{"type":"weekly","week_days":[1,3,5]}'

# Set specific dates
curl -X PUT /api/v1/tasks/1/recurrence -d '{"type":"specific_dates","dates":["2026-06-01","2026-07-01"]}'

# Set even days
curl -X PUT /api/v1/tasks/1/recurrence -d '{"type":"even_odd","parity":"even"}'