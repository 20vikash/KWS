include .env

up:
	docker compose up
down:
	docker compose down
stop:
	docker compose stop
start:
	docker compose start
	docker compose logs -f
dv:
	docker volume rm kws_postgres_db_data_kws
	docker volume rm kws_redis_db_data_kws
	docker volume rm kws_mq_kws

dvs:
	docker volume rm kws_postgres_db_service_data

create_migration:
	migrate create -ext=sql -dir=src/internal/database/migrations -seq init

migrate_up:
	migrate -path=src/internal/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DBNAME}?sslmode=disable" -verbose up

migrate_down-%:
	migrate -path=src/internal/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DBNAME}?sslmode=disable" -verbose down $*

migrate_down-all:
	migrate -path=src/internal/database/migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DBNAME}?sslmode=disable" -verbose down

.PHONY: create_migration migrate_up migrate_down
