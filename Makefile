.PHONY: build run up migrate_up migrate_down

up:
	docker-compose up --build app db migrate_up

migrate_up:
	docker-compose up -d migrate_up

migrate_down:
	docker-compose up -d migrate_down

build:
	docker-compose up --build app

run:
	docker-compose up app
