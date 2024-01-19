PG_DSN_TEST="postgresql://test:test@localhost:54321/postgres?sslmode=disable"

gen-api-types:
	oapi-codegen -generate "types" -package api ./api/api.json > ./internal/delivery/http/internal/api/types.gen.go

gen-api-server:
	oapi-codegen -generate "chi-server" -package api ./api/api.json > ./internal/delivery/http/internal/api/server.gen.go

install-dependencies:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: format
format:
	goimports -local github.com/Employee-s-file-cabinet/backend -w ./internal/

.PHONY: build-dev
build-dev:
	go build -o server -tags=dev  ./cmd/

.PHONY: new-migration
new-migration: ## Создание новой миграции (задать name=...)
	goose -dir migrations create $(name) sql

.PHONY: migrate-up
migrate-up: ## Запуск всех миграций (до последней версии)
	goose -dir migrations postgres "$(PG_DSN)" up

.PHONY: migrate-up-by-one
migrate-up-by-one: ## Запуск миграций по одной
	goose -dir migrations postgres "$(PG_DSN)" up-by-one

.PHONY: migrate-down
migrate-down: ## Откат на одну миграцию
	goose -dir migrations postgres "$(PG_DSN)" down

.PHONY: migrate-reset
migrate-reset: ## Откатить все миграции
	goose -dir migrations postgres "$(PG_DSN)" reset

.PHONY: migrate-status
migrate-status: ## Откатить все миграции
	goose -dir migrations postgres "$(PG_DSN)" status

.PHONY: test-migrate-up
test-migrate-up: ## Запуск миграции на тестовую базу
	migrate -path migrations -database "$(PG_DSN_TEST)" -verbose up

.PHONY: test-migrate-down
test-migrate-down: ## Откат миграций тестовой базы
	migrate -path migrations -database "$(PG_DSN_TEST)" -verbose down

.PHONY: test-env-up
test-env-up: ## Запуск тестового окружения.
	docker compose -f 'docker-compose-test.yaml' up --exit-code-from migrate migrate

.PHONY: test-env-down
test-env-down: ## Останов и очистка тестового окружения.
	docker compose -f 'docker-compose-test.yaml' down -v

.PHONY: run-tests
run-tests: ## Запуск unit-тестов
	go test -count 1 -coverpkg=./... -race -v ./...