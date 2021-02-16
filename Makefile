.PHONY: build run up down migrate-up migrate-down

PROJECT_NAME:="octo-server"
APP_CONTAINER_NAME=app
COMPOSE_RUN := docker-compose run --entrypoint "" $(APP_CONTAINER_NAME)
MODULE = $(shell go list -m)

-include .env

up: ## Starts the application containers
	docker-compose up -d

down: ## Stops the applications
	docker-compose down

migrate-up: ## Up the Database migrations
	$(COMPOSE_RUN) migrate -path=migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable up

migrate-down: ## Down the Database migrations
	$(COMPOSE_RUN) migrate -path=migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable down 1

build: ## Builds the application using Docker Compose
	docker-compose build $(APP_CONTAINER_NAME)

run: ## run the API server
	docker-compose up app
