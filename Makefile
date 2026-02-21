# Variables
DB_URL=postgres://postgres:password@localhost:5432/football_db?sslmode=disable

.PHONY: db-up db-down db-force docker-up docker-down docker-start docker-stop

# -------------
#  Migrations
# -------------

## Run database migrations up
migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

## Run database migrations down (1 step by default)
migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down 1

## Force database migration version
migrate-force:
	migrate -path migrations -database "$(DB_URL)" -verbose force $(version)


# -------------
# Docker Compose
# -------------

## Start docker containers in the background (build if necessary)
docker-up:
	docker-compose up -d --build

## Stop and remove docker containers, networks, and volumes
docker-down:
	docker-compose down

## Start existing docker containers
docker-start:
	docker-compose start

## Stop running docker containers without removing them
docker-stop:
	docker-compose stop
