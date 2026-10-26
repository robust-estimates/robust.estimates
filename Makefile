SHELL := /bin/sh
.DEFAULT_GOAL := help
.NOTPARALLEL:

ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
PYTHON ?= $(shell for candidate in python3.13 python3.12 python3.11 python3.10 python3; do command -v $$candidate 2>/dev/null && break; done)
GO ?= go
DOCKER ?= docker
COMPOSE ?= $(DOCKER) compose
VENV ?= $(ROOT)/.venv
VPY := $(VENV)/bin/python
API_URL ?= http://localhost:8080/api/v1
WAIT_TIMEOUT ?= 180
REPORT ?= $(ROOT)/build/api-report.json
E2E_FLAGS ?= --extended
TARGET_ARCH ?= amd64
export GOTOOLCHAIN ?= go1.27.3

GO_PACKAGES := ./platform/... ./company-service/... ./application-service/... ./gateway/...
GO_SERVICES := gateway company-service application-service
GO_DIRS := gateway company-service application-service platform contracts scripts deploy/minio
RUN := $(MAKE) --no-print-directory -C "$(ROOT)"

.PHONY: help doctor doctor-host doctor-docker env setup deps proto-deps fmt fmt-check vet test test-go test-python build build-linux api compose-check check check-go up start down restart status logs ready e2e e2e-contract proto

help:
	@printf '%s\n' \
	  'make check          Финальная проверка: код, тесты, Docker, все 34 операции API' \
	  'make check-go       Только Go: форматирование, vet, тесты с race и сборка' \
	  'make setup          Создать отсутствующий .env и установить тестовые зависимости в .venv' \
	  'make up             Собрать образы и запустить Docker с ожиданием готовности' \
	  'make start          Запустить Docker без пересборки образов' \
	  'make e2e            Проверить работающий API, сохранить JSON-отчёт' \
	  'make e2e-contract   Проверить E2E со сверкой ответов по OpenAPI' \
	  'make test           Тесты Go и Python' \
	  'make test-go        Тесты Go с детектором гонок' \
	  'make api            Валидация OpenAPI' \
	  'make fmt            Отформатировать Go-код' \
	  'make build          Собрать три Go-сервиса в build/' \
	  'make build-linux    Собрать Linux-бинарники в build/linux/' \
	  'make proto          Перегенерировать Go/Python protobuf bindings' \
	  'make status         Состояние контейнеров' \
	  'make logs           Последние логи сервисов' \
	  'make down           Остановить контейнеры, сохранить volumes' \
	  '' \
	  'Пример: make e2e E2E_FLAGS="--verbose --extended"' \
	  'Другой адрес: make e2e API_URL=http://localhost:8081/api/v1'

doctor-host:
	@command -v "$(GO)" >/dev/null 2>&1 || { printf '%s\n' 'Не найден Go. Установите Go 1.27.3 или укажите GO=/path/to/go.'; exit 1; }
	@test -n "$(PYTHON)" && test -x "$(PYTHON)" || command -v "$(PYTHON)" >/dev/null 2>&1 || { printf '%s\n' 'Не найден Python. Укажите PYTHON=/path/to/python3.'; exit 1; }
	@"$(GO)" version
	@"$(PYTHON)" -c 'import sys; print("Python " + sys.version.split()[0]); sys.exit(0 if sys.version_info >= (3, 10) else "Нужен Python >= 3.10; укажите PYTHON=/path/to/python3.13")'

doctor-docker:
	@command -v "$(DOCKER)" >/dev/null 2>&1 || { printf '%s\n' 'Не найден Docker.'; exit 1; }
	@$(COMPOSE) version
	@$(DOCKER) info >/dev/null 2>&1 || { printf '%s\n' 'Docker Engine недоступен. Запустите Docker Desktop и повторите команду.'; exit 1; }

doctor:
	@$(RUN) doctor-host
	@$(RUN) doctor-docker

env:
	@cd "$(ROOT)" && if [ ! -f .env ]; then cp .env.example .env && printf '%s\n' 'Создан .env из .env.example'; fi

$(VENV)/.backend-test-deps: $(ROOT)/scripts/requirements-test.txt $(ROOT)/scoring/requirements.txt
	@test -n "$(PYTHON)" || { printf '%s\n' 'Не найден Python.'; exit 1; }
	@test -x "$(VPY)" || "$(PYTHON)" -m venv "$(VENV)"
	@"$(VPY)" -c 'import sys; sys.exit(0 if sys.version_info >= (3, 10) else "В .venv нужен Python >= 3.10; укажите другую папку VENV")'
	@"$(VPY)" -m pip install --disable-pip-version-check --no-cache-dir -r "$(ROOT)/scripts/requirements-test.txt" -r "$(ROOT)/scoring/requirements.txt"
	@touch "$@"

deps: $(VENV)/.backend-test-deps

setup:
	@$(RUN) env
	@$(RUN) deps

fmt:
	@cd "$(ROOT)" && find $(GO_DIRS) -type f -name '*.go' -exec gofmt -w {} +

fmt-check:
	@cd "$(ROOT)" && files=$$(find $(GO_DIRS) -type f -name '*.go' -exec gofmt -l {} +) && if [ -n "$$files" ]; then printf '%s\n' 'Требуется gofmt:' "$$files" 'Выполните make fmt.'; exit 1; fi

vet:
	@cd "$(ROOT)" && "$(GO)" vet $(GO_PACKAGES)

test-go:
	@cd "$(ROOT)" && "$(GO)" test -race -count=1 $(GO_PACKAGES)

test-python: deps
	@cd "$(ROOT)/scoring" && PYTHONPATH=generated:. "$(VPY)" -m unittest discover -s tests -v

test:
	@$(RUN) test-go
	@$(RUN) test-python

build:
	@mkdir -p "$(ROOT)/build"
	@set -e; for service in $(GO_SERVICES); do \
	  printf 'Сборка %s\n' "$$service"; \
	  (cd "$(ROOT)/$$service" && GOWORK=off CGO_ENABLED=0 "$(GO)" build -trimpath -o "$(ROOT)/build/$$service" ./cmd/server); \
	done

build-linux:
	@mkdir -p "$(ROOT)/build/linux"
	@set -e; for service in $(GO_SERVICES); do \
	  printf 'Сборка Linux/%s %s\n' "$(TARGET_ARCH)" "$$service"; \
	  (cd "$(ROOT)/$$service" && GOWORK=off CGO_ENABLED=0 GOOS=linux GOARCH="$(TARGET_ARCH)" "$(GO)" build -trimpath -o "$(ROOT)/build/linux/$$service" ./cmd/server); \
	done

api: deps
	@cd "$(ROOT)" && "$(VPY)" scripts/validate_contract.py

compose-check: env
	@cd "$(ROOT)" && $(COMPOSE) config --quiet

check-go:
	@$(RUN) fmt-check
	@$(RUN) vet
	@$(RUN) test-go
	@$(RUN) build

check:
	@$(RUN) doctor
	@$(RUN) setup
	@$(RUN) compose-check
	@$(RUN) check-go
	@$(RUN) test-python
	@$(RUN) api
	@$(RUN) up
	@$(RUN) ready
	@$(RUN) e2e
	@$(RUN) e2e-contract
	@printf '\n%s\n' 'Все проверки backend пройдены.' 'Контейнеры оставлены работающими; остановка: make down.' 'Отчёт: $(REPORT)'

up: env doctor-docker
	@cd "$(ROOT)" && if $(COMPOSE) up --build -d --wait --wait-timeout "$(WAIT_TIMEOUT)"; then :; else \
	  $(COMPOSE) ps -a; \
	  $(COMPOSE) logs --tail=60 database-init company-service application-service gateway; \
	  exit 1; \
	fi

start: env doctor-docker
	@cd "$(ROOT)" && $(COMPOSE) up -d --wait --wait-timeout "$(WAIT_TIMEOUT)"

down:
	@cd "$(ROOT)" && $(COMPOSE) down

restart: doctor-docker
	@cd "$(ROOT)" && $(COMPOSE) restart
	@$(RUN) start

status:
	@cd "$(ROOT)" && $(COMPOSE) ps -a

logs:
	@cd "$(ROOT)" && $(COMPOSE) logs --tail=100 database-init gateway company-service application-service scoring

ready:
	@"$(PYTHON)" -c 'import json,sys,urllib.error,urllib.request; url="$(API_URL)".rstrip("/"); url=url[:-7]+"/readyz" if url.endswith("/api/v1") else sys.exit("API_URL должен заканчиваться на /api/v1 для make ready"); response=urllib.request.urlopen(url,timeout=5); data=json.load(response); print("API готов:", data); sys.exit(0 if response.status==200 and data.get("status")=="ready" else 1)'

e2e:
	@cd "$(ROOT)" && "$(PYTHON)" scripts/check_api.py --base-url "$(API_URL)" $(E2E_FLAGS) --report "$(REPORT)"

e2e-contract: deps
	@cd "$(ROOT)" && CHECK_CONTRACT=1 API_URL="$(API_URL)" "$(VPY)" scripts/e2e.py

$(VENV)/.backend-proto-deps: $(VENV)/.backend-test-deps $(ROOT)/scoring/requirements-dev.txt
	@"$(VPY)" -m pip install --disable-pip-version-check --no-cache-dir -r "$(ROOT)/scoring/requirements-dev.txt"
	@touch "$@"

proto-deps: $(VENV)/.backend-proto-deps

proto: proto-deps
	@cd "$(ROOT)" && "$(VPY)" scripts/generate.py
