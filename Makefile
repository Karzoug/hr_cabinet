PG_DSN_TEST="postgresql://test:test@localhost:54321/postgres?sslmode=disable"

.PHONY: new-migration
new-migration: ## Создание новой миграции (задать name=...)
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-up
migrate-up: ## Запуск миграции
	migrate -path internal/db/migrations -database "$(PG_DSN)" -verbose up

.PHONY: migrate-down
migrate-down: ## Откат миграций
	migrate -path internal/db/migrations -database "$(PG_DSN)" -verbose down

.PHONY: test-migrate-up
test-migrate-up: ## Запуск миграции на тестовую базу
	migrate -path internal/db/migrations -database "$(PG_DSN_TEST)" -verbose up

.PHONY: test-migrate-down
test-migrate-down: ## Откат миграций тестовой базы
	migrate -path internal/db/migrations -database "$(PG_DSN_TEST)" -verbose down

.PHONY: test-env-up
test-env-up: ## Запуск тестового окружения.
	docker compose -f 'docker-compose-test.yaml' up --exit-code-from migrate migrate

.PHONY: test-env-down
test-env-down: ## Останов и очистка тестового окружения.
	docker compose -f 'docker-compose-test.yaml' down -v
