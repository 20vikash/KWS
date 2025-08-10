include .env

ATTACH_SERVICES = \
	postgres.kws.services:lxdbr0:172.30.0.100/24 \
	adminer.kws.services:lxdbr0:172.30.0.101/24 \
	dnsmasq_kws:lxdbr0:172.30.0.102/24

define attach_services
	@echo "Attaching services to bridge..."
	@set -e; \
	for triple in $(ATTACH_SERVICES); do \
		container=$${triple%%:*}; \
		tmp=$${triple#*:}; \
		bridge=$${tmp%%:*}; \
		ipcidr=$${tmp#*:}; \
		echo " -> $$container to $$bridge with $$ipcidr"; \
		attach_to_bridge $$container $$bridge $$ipcidr; \
	done
endef

up:
	docker compose up -d
	$(call attach_services)
	docker compose logs -f

down:
	docker compose down

stop:
	docker compose stop

start:
	docker compose start
	$(call attach_services)
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
	migrate -path=src/internal/database/migrations \
		-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DBNAME}?sslmode=disable" \
		-verbose up

migrate_down-%:
	migrate -path=src/internal/database/migrations \
		-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DBNAME}?sslmode=disable" \
		-verbose down $*

migrate_down-all:
	migrate -path=src/internal/database/migrations \
		-database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DBNAME}?sslmode=disable" \
		-verbose down

.PHONY: up down stop start dv dvs create_migration migrate_up migrate_down migrate_down-all
