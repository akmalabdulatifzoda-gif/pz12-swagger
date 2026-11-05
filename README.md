# ПЗ №12 — Swagger/OpenAPI для Notes API

**Студент:** Абдулатифзода Акмал, группа ЭФМО-0126
**Дисциплина:** Технологии индустриального программирования (преподаватель Адышкин С. С.)

Продолжение [ПЗ №11](https://github.com/akmalabdulatifzoda-gif/pz11-notes-api): к CRUD API заметок подключена
автогенерация документации по подходу **code-first** — аннотации в коде → `swag init` → Swagger UI.

| Адрес | Что открывается |
|---|---|
| http://localhost:8080/docs/index.html | Swagger UI (интерактивная документация, Try it out) |
| http://localhost:8080/docs/doc.json   | сгенерированная спецификация |
| http://localhost:8080/redoc           | ReDoc (статическая страница, скрипт встроен в бинарник) |

## Быстрый старт

```bash
make tools        # один раз: установить swag CLI (v1.16.4)
make swagger      # сгенерировать docs/ и экспортировать спецификацию в api/
make run          # запустить сервер на :8080
```

С авторизацией (POST/PATCH/DELETE требуют `Authorization: Bearer <token>`):

```bash
API_TOKEN=secret make run
```

В Swagger UI нажмите **Authorize** и введите `Bearer secret`. Если `API_TOKEN` не задан, авторизация выключена.

## Как обновлять документацию

1. Изменили хэндлер — поправьте аннотации над ним (`@Summary`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security`).
2. Выполните `make swagger` — перегенерируются `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`,
   а также `api/swagger.yaml` (Swagger 2.0) и `api/openapi.yaml` (конвертация в OpenAPI 3.0 через swagger2openapi).
3. `make lint` — проверка `api/openapi.yaml` линтером Redocly (правила в `redocly.yaml`).
4. `make all` = swagger + lint + test.

Каталог `docs/` генерируется автоматически — руками его не редактируем.

## Makefile

| Цель | Действие |
|---|---|
| `make tools`   | установить `swag` |
| `make swagger` | `swag fmt` + `swag init -g cmd/api/main.go -o docs --parseInternal` + экспорт в `api/` |
| `make lint`    | `redocly lint api/openapi.yaml` |
| `make run`     | `go run ./cmd/api` |
| `make test`    | `go test ./...` (CRUD, пагинация, Bearer, доступность /docs и /redoc) |

## Эндпоинты

| Метод | Путь | Описание | Авторизация |
|---|---|---|---|
| GET    | /api/v1/notes?page=&limit=&q= | список, пагинация, поиск по title, заголовок `X-Total-Count` | — |
| POST   | /api/v1/notes      | создать (201) | Bearer |
| GET    | /api/v1/notes/{id} | получить (200/404) | — |
| PATCH  | /api/v1/notes/{id} | частично обновить (200/400/404) | Bearer |
| DELETE | /api/v1/notes/{id} | удалить (204/404) | Bearer |

Модели в спецификации: `core.Note`, `core.NoteCreate`, `core.NoteUpdate`, `core.ErrorResponse` (с примерами значений).

## Структура

```
├── Makefile, redocly.yaml
├── api/            openapi.yaml (OAS 3.0), swagger.yaml (экспорт)
├── cmd/api/        main.go — общие аннотации (@title, @BasePath, @securityDefinitions)
├── docs/           сгенерировано swag: docs.go, swagger.json, swagger.yaml
└── internal/
    ├── core/       модели и DTO, service/ — бизнес-логика
    ├── http/       router.go, auth.go (Bearer), handlers/ (аннотации), redoc/
    └── repo/       in-memory хранилище
```

## Отчёт

[report/](report) — отчёт в Word и скриншоты.
