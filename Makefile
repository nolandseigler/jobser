

.PHONY: up up-d down down-v

up:
	docker-compose up --build

up-d:
	docker-compose up --build -d

down:
	docker-compose down

down-v:
	docker-compose down -v
