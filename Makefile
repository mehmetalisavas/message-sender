BINARY_NAME=message-sender

GO_CMD=go
SWAG_CMD=swag
SWAG_VERSION=v1.8.4


# Generate Swagger Docs using swag inside a Docker container
docs:
	@echo "Generating Swagger Docs using Docker..."
	@docker run --rm -v $(PWD):/code ghcr.io/swaggo/swag:latest init -g cmd/main.go # Adjust the path to your main Go file if needed
	@echo "Swagger documentation generated!"


build:
	go build -o bin/$(BINARY_NAME) cmd/main.go


docker-build:
	docker-compose  --env-file .env  -f ./infra/docker-compose.yml up --build 

build-test:
	docker-compose --env-file .env_test -f ./infra/docker-compose.test.yml up --build --exit-code-from app-test
	docker-compose --env-file .env_test -f ./infra/docker-compose.test.yml down --volumes --remove-orphans

test-docker: 
	docker-compose --env-file .env_test -f ./infra/docker-compose.test.yml down --volumes --remove-orphans
	docker-compose --env-file .env_test -f ./infra/docker-compose.test.yml up migrations
	docker-compose --env-file .env_test -f ./infra/docker-compose.test.yml up --build app-test
	docker-compose --env-file .env_test -f ./infra/docker-compose.test.yml run app-test go test -v ./...
	docker-compose --env-file .env_test -f ./infra/docker-compose.test.yml down --volumes --remove-orphans