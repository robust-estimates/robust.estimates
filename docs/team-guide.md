# Как работать с проектом командой

В репозиторий передаётся весь исходный проект. Frontend обращается только к gateway по HTTP. Gateway вызывает Go-сервисы по gRPC. Application service вызывает Python scoring по gRPC и сохраняет результат. Компания, заявки и авторизация имеют отдельные базы; документы лежат в MinIO.

## Карта файлов

Все пути ниже указаны от корня репозитория.

### Запуск и общие контракты

| Файл или папка | Назначение |
|---|---|
| `README.md` | Запуск, учётные записи, соглашения API, ограничения MVP |
| `Makefile` | Команды запуска, сборки, проверок и генерации protobuf |
| `compose.yaml` | Контейнеры, сеть, volumes, проверки готовности и зависимости запуска |
| `.env.example` | Шаблон локальной конфигурации; копируется в `.env` |
| `.gitignore` | Исключает локальные секреты, окружения, сборки и отчёты |
| `.dockerignore` | Исключает лишние файлы из Docker build context |
| `.gitattributes` | Правила перевода строк |
| `go.work`, `go.work.sum` | Совместная локальная разработка Go-модулей и контрольные суммы |
| `api/public/openapi.yaml` | Публичный HTTP-контракт для frontend и Swagger |
| `api/proto/credit/v1/credit.proto` | Внутренний gRPC-контракт Go-сервисов и Python scoring |
| `contracts/credit/v1/credit.pb.go` | Сгенерированные Go-типы protobuf |
| `contracts/credit/v1/credit_grpc.pb.go` | Сгенерированные Go-клиенты и серверные интерфейсы gRPC |
| `THIRD_PARTY_NOTICES` | Сведения о сторонних компонентах |

Сгенерированные protobuf-файлы коммитим, но вручную не редактируем. При изменении `.proto` выполняем `make proto` и включаем обновлённые Go- и Python-файлы в тот же коммит. `go.mod` и `go.sum`, если присутствуют в модуле, тоже коммитим.

### Gateway

| Путь внутри `gateway/` | Назначение |
|---|---|
| `cmd/server/main.go` | Подключение БД, миграции, gRPC-клиенты, HTTP-сервер и остановка |
| `internal/domain/user.go` | Пользователь, роли и ошибки учётных записей |
| `internal/usecase/auth.go` | Регистрация, проверка пароля, выдача сессии; регистрация создаёт CLIENT |
| `internal/adapter/http/server.go` | REST-маршруты, JWT, CORS, заголовки команд, преобразование HTTP ↔ gRPC, файлы и Swagger |
| `internal/adapter/postgres/users.go` | Хранение пользователей и создание демонстрационных учёток |
| `internal/usecase/auth_test.go` | Проверки регистрации и входа на уровне сценариев |
| `internal/adapter/http/auth_test.go` | Проверки HTTP регистрации, ограничений ввода и частоты запросов |
| `migrations/001_init.sql` | Начальная схема auth_db |
| `migrations/002_login_case_insensitive.sql` | Уникальность логина без учёта регистра |
| `go.mod`, `go.sum`, `Dockerfile` | Зависимости модуля и сборка контейнера |

### Company service

| Путь внутри `company-service/` | Назначение |
|---|---|
| `cmd/server/main.go` | Сборка зависимостей и запуск gRPC-сервера |
| `internal/domain/company.go` | Типы компании и продуктов, проверки данных, расчёт платежей |
| `internal/domain/company_test.go` | Проверки расчётов и правил предметной области |
| `internal/usecase/company.go` | Профили компаний, проверка, каталог, подбор и калькулятор |
| `internal/adapter/grpc/server.go` | Приём gRPC-запросов и преобразование моделей |
| `internal/adapter/postgres/repository.go` | Чтение и запись компаний и каталога |
| `migrations/001_init.sql` | Схема company_db |
| `go.mod`, `go.sum`, `Dockerfile` | Зависимости и контейнер |

### Application service

| Путь внутри `application-service/` | Назначение |
|---|---|
| `cmd/server/main.go` | БД, S3, gRPC-клиенты, сервер и запуск фоновой обработки |
| `internal/domain/application.go` | Заявка, статусы, документы, предложения, решения и связанные типы |
| `internal/usecase/service.go` | Жизненный цикл заявки, проверки доступа, документы, решения, история, лента и worker скоринга |
| `internal/usecase/service_test.go` | Проверки бизнес-сценариев заявки |
| `internal/adapter/grpc/server.go` | Серверный gRPC API, включая потоковую передачу документов |
| `internal/adapter/grpc/clients.go` | Вызовы company-service и scoring; проверка соответствия ответа скоринга запросу |
| `internal/adapter/postgres/repository.go` | Агрегаты заявок, транзакции, очередь фоновых заданий и состояния S3-объектов |
| `internal/adapter/s3/storage.go` | Запись, чтение и удаление файлов в S3/MinIO |
| `migrations/001_init.sql` | Схема application_db |
| `go.mod`, `go.sum`, `Dockerfile` | Зависимости и контейнер |

`domain` описывает предметную область, `usecase` выполняет сценарии, `adapter` связывает их с транспортом и хранилищами. `cmd/server` собирает зависимости. Сценарии не должны зависеть от конкретного HTTP-сервера или S3 SDK. Применённые SQL-миграции не меняем: добавляем следующий файл.

### Python scoring

| Путь внутри `scoring/` | Назначение |
|---|---|
| `app/model.py` | Текущая модель `RuleModel` и тип результата `Risk`; основная точка подключения ML |
| `app/server.py` | gRPC-сервер, проверка служебного JWT, проверка входа, вызов модели и healthcheck |
| `generated/credit/v1/credit_pb2.py` | Сгенерированные protobuf-сообщения |
| `generated/credit/v1/credit_pb2_grpc.py` | Сгенерированные gRPC-интерфейсы |
| `tests/test_model.py` | Проверки модели |
| `requirements.txt` | Зависимости исполнения |
| `requirements-dev.txt` | Зависимости разработки и генерации |
| `Dockerfile` | Контейнер Python-сервиса |

### Общая инфраструктура и проверки

| Файл или папка | Назначение |
|---|---|
| `platform/core/core.go` | Общие типы идентичности, ошибок и команд без инфраструктурных зависимостей |
| `platform/platform.go` | Конфигурация, JWT, gRPC, БД, healthcheck и преобразование protobuf |
| `platform/migrations.go` | Выполнение и контроль SQL-миграций |
| `platform/transaction.go` | Транзакции и сохранение результатов идемпотентных команд |
| `platform/platform_test.go` | Проверки общей инфраструктуры |
| `deploy/postgres/init/001_databases.sh` | Создание баз и ролей PostgreSQL, настройка паролей |
| `deploy/minio/Dockerfile` | Сборка MinIO из зафиксированных исходников |
| `deploy/minio/healthcheck.go` | Проверка готовности MinIO |
| `scripts/check_api.py` | Консольная проверка API, в том числе загрузки и скачивания файлов; используется `make e2e` |
| `scripts/e2e.py` | Сквозные сценарии кредитного конвейера |
| `scripts/validate_contract.py` | Проверка OpenAPI и соответствия ответов контракту |
| `scripts/generate.py` | Генерация protobuf Go/Python |
| `scripts/strip-comments.go` | Удаление комментариев из сгенерированных Go-файлов |
| `scripts/requirements-test.txt` | Python-зависимости проверки контрактов |
| `testdata/scenarios.json` | Синтетические тестовые сценарии |
| `.github/` | Настройки автоматических проверок репозитория |
| `docs/swagger.html` | Страница Swagger UI, загружающая актуальный OpenAPI |
| `docs/architecture.md` | Архитектура, границы ответственности и схемы |
| `docs/registration.md` | Регистрация, ограничения и результаты её проверки |
| `docs/verification.md`, `docs/console-verification.md` | Записи о выполненных проверках; даты и условия относятся к конкретным запускам |

## Что делает frontend

Исходники интерфейса кладём в корневую папку `frontend/`. Фронтендер выбирает свой стек и коммитит код, package manifest и lock-файл. `node_modules` не коммитим. Текущий Compose запускает backend и scoring; контейнер frontend при необходимости добавляется отдельно.

1. Получить проект и запустить `make up` из корня.
2. Открыть Swagger: `http://localhost:8080/docs/`.
3. Использовать REST base URL `http://localhost:8080/api/v1`. Например, для Vite можно завести `VITE_API_BASE_URL`; это настройка frontend, backend её не читает.
4. Запускать интерфейс на `http://localhost:5173` либо согласовать origin в корневом `.env` через `CORS_ORIGIN` и пересоздать gateway командой `docker compose up -d --force-recreate gateway`.
5. Реализовать вход или регистрацию и передавать `Authorization: Bearer <accessToken>` в защищённые методы. Токен не помещать в репозиторий.

Для POST/PATCH/DELETE нужен новый `Idempotency-Key` длиной 8–128 символов на новое пользовательское действие, кроме входа и регистрации. При сетевом повторе того же действия сохраняем ключ и тело. Для изменения существующей заявки передаём актуальный `If-Match`, например `"3"`; версию берём из ответа заявки. После загрузки или удаления документа запрашиваем заявку заново. При 409 перечитываем данные, а не повторяем изменение бесконечно.

Суммы — целые копейки, ставки — базисные пункты: `100000000` = 1 млн ₽, `1690` = 16,9%. Имена JSON-полей и допустимые значения берём из OpenAPI. Файлы отправляем через `FormData`, без ручной установки multipart Content-Type: браузер добавит boundary. Документы скачиваем через авторизованный HTTP-запрос и Blob. Прямой доступ браузера к MinIO не требуется.

Статусы после отправки меняются асинхронно: frontend периодически перечитывает заявку и решение. Кнопки и переходы строим с учётом `status`, `nextAction` и ошибок API. При 401 потребуется новый вход: refresh пока отсутствует. Роль EMPLOYEE нужна для банковской очереди и ручных решений.

## Что делает ML-разработчик

Меняем `scoring/app/model.py`, необходимые зависимости и тесты в `scoring/`. Сейчас там правила `rules-v1`, а не обученная модель. Если меняется имя класса или способ загрузки, обновляем создание модели в `app/server.py`. Модель желательно загрузить один раз при старте, а не на каждом запросе; сервер использует несколько потоков, поэтому инференс и состояние модели должны быть безопасны для конкурентных вызовов.

Транспорт сохраняем: `ScoringService.Evaluate` из `api/proto/credit/v1/credit.proto`. Application service уже вызывает его. Входные признаки:

- `business_age_months` — возраст бизнеса в месяцах;
- `monthly_revenue_minor`, `monthly_expenses_minor`, `monthly_debt_payments_minor` — месячные суммы в копейках;
- `requested_amount_minor` — запрошенная сумма в копейках;
- `term_months` — срок в месяцах;
- `has_overdue_debt` — наличие просроченной задолженности.

Результат содержит `score` от 0 до 1000, `risk_level` (`LOW`, `MEDIUM`, `HIGH`), `factors`, `model_version` и `recommendation` (`APPROVE`, `MANUAL_REVIEW`, `REJECT`). Текущий `Risk` использует имя `level`, а `server.py` преобразует его в protobuf `risk_level`. Идентификаторы `scoring_request_id`, `application_id`, `application_version` должны вернуться без изменений; текущий сервер это уже делает.

Не используем score как вероятность без согласованной шкалы. При замене модели согласовываем соответствие оценки, уровня риска и рекомендации; итоговые статусы заявки и решение принадлежат application-service. Сейчас MEDIUM/HIGH направляются сотруднику, LOW проходит дополнительные правила продукта и доступной суммы.

На один gRPC-вызов Go-клиент отводит 10 секунд. Если нужен больший бюджет, меняем его согласованно. Валидацию, служебную авторизацию и healthcheck в сервере сохраняем. Нужны дополнительные признаки — сначала согласовываем их источник и меняем `.proto` и построение запроса в application-service, затем выполняем `make proto` и коммитим все изменения вместе.

Модельные артефакты можно хранить в `scoring/models/`, если они небольшие и разрешены к распространению. Большие веса кладём в отдельное хранилище, в Git сохраняем версию, checksum и инструкцию получения. При добавлении файлов модели нужно обновить Dockerfile: сейчас он копирует только `app/` и `generated/`. Датасеты с реальными персональными данными в проект не добавляем.

```sh
make test-python
docker compose up -d --build --wait scoring
make e2e
```

После смены реальной модели ожидаемые решения E2E могут отличаться от `rules-v1`. Согласовываем тестовые данные и ожидаемые исходы; сам транспортный контракт сохраняем.
