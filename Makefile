.PHONY: migrate-up migrate-down migrate-create migrate-force test run

# Database connection for migrations
DB_URL=mysql://root:rootpassword@tcp(localhost:3306)/store

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path migrations -database "$(DB_URL)" force $$version

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

test:
	go test ./... -v

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

# Docker commands
docker-build:
	docker build -t store-api .

docker-up:
	docker compose up -d

docker-down:
	docker-compose down 

docker-logs:
	docker-compose logs -f api 

docker-restart:
	docker-compose restart api

docker-clean:
	docker-compose down -v
	docker rmi store-api

docker-shell:
	docker exec -it store_api sh

docker-mysql:
	docker exec -it store_mysql mysql -u root -p

docker-redis:
	docker exec -it store_redis redis-cli

redis-cli:
	redis-cli