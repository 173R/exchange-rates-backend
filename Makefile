# Помещаем переменные окружения.
include .env

MIGRATIONS = "internal/db/migrations"
DB_URI = "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Создает новую миграцию.
migrations-create:
	@echo "Migration name: "
	@read -r name && \
	migrate create -ext sql -dir $(MIGRATIONS) $$name

# Запускает все миграции.
migrations-up:
	@migrate -database $(DB_URI) -path $(MIGRATIONS) up

# Откатывает указанное количество миграций.
migrations-down:
	@echo "Count to rollback (leave empty to rollback everything): "
	@read -r count && \
	migrate -database $(DB_URI) -path $(MIGRATIONS) down $$count
