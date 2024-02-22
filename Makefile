.PHONY: build
build:
	docker build -t app .

.PHONY: up
up:
	docker compose up -d --build

.PHONY: down
down:
	docker compose down

.PHONY: migrate-up
migrate-up:
	@go run cmd/main.go migrate up

.PHONY: migrate-down
migrate-down:
	@go run cmd/main.go migrate down
