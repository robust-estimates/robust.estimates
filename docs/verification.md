# Результаты проверки backend

Изменён репозиторий `/Users/cofffeen/Desktop/credit-api/robust.estimates`.

Реализованы gateway, company-service, application-service на Go, gRPC-контракты и Python-сервис rules-v1 для подключения будущей ML-модели. Код без комментариев. Docker Compose включает PostgreSQL с раздельными базами/ролями и MinIO, собираемый из официальных исходников на фиксированном коммите.

## Выполнено

- Go 1.27.3: go test -race и go vet — успешно.
- Linux amd64, CGO_ENABLED=0: все три Go-сервиса собраны с GOWORK=off.
- Python: 3 unit-теста скоринга и компиляция исходников — успешно.
- OpenAPI 3.0.3, 33 операции: формальная валидация — успешно.
- docker compose config — успешно.
- E2E: реальные процессы Go + Python, настоящий MinIO из зафиксированного исходного кода и PGlite через PostgreSQL wire protocol. Все проверки пройдены, включая сверку JSON-ответов с OpenAPI.
- Проверены одобрение/выдача, альтернативные условия, ручная проверка/запрос документов/повторная отправка/отказ, ручное предложение, отмена, права доступа, PDF roundtrip, повторы команд и конфликт содержимого.

## Ограничение проверки окружения

Полный Docker Compose не запускался: Docker Engine недоступен, приложение Docker Desktop имеет отсутствующий executable. Обычный локальный PostgreSQL не запустился из-за запрета shared memory в песочнице. Для интеграционной проверки SQL использован PGlite; его результаты не заменяют проверку конкурентных транзакций и отказоустойчивости на PostgreSQL 17 в Compose. Workflow CI содержит полный запуск Compose и E2E для окружения с Docker.

## Запуск у команды

```sh
cd /Users/cofffeen/Desktop/credit-api/robust.estimates
cp .env.example .env
docker compose up --build -d --wait
python3 scripts/e2e.py
```

Swagger: http://localhost:8080/docs/ . Пароль тестовых учёток — DEMO_PASSWORD из .env.example. Подробности в README репозитория. Frontend не изменён. Redis/Kafka/Patroni/etcd/HAProxy/Prometheus/Grafana — будущие этапы, в этот Compose не входят.
