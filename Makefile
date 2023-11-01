

.PHONY: up up-d down down-v

up:
	docker-compose up --build --force-recreate  --remove-orphans

up-d:
	docker-compose up --build --force-recreate  --remove-orphans -d

down:
	docker-compose down

down-v:
	docker-compose down -v
