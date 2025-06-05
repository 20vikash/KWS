include .env

up:
	docker compose up
down:
	docker compose down
stop:
	docker compose stop

create_migration:
	migrate create -ext=sql -dir=internal/database/migrations -seq init

migrate_up:
	migrate -path=internal/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DBNAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=internal/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DBNAME}?sslmode=disable" -verbose down

.PHONY: create_migration migrate_up migrate_down
