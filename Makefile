.PHONY: build
build:
	docker build -t app .

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: migrate-up
migrate-up:
	migrate -path ./db/migration -database "postgresql://app:app@localhost:5432/app?sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path ./db/migration -database "postgresql://app:app@localhost:5432/app?sslmode=disable" -verbose down
