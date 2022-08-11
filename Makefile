# Помещаем переменные окружения.
include .env

MIGRATIONS = "internal/db/migrations"
DB_URI = "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

migrations-create:
	@echo "Migration name: "
	@read -r name && \
	migrate create -ext sql -dir $(MIGRATIONS) $$name

migrations-up:
	@migrate -database $(DB_URI) -path $(MIGRATIONS) up