# Variables
APP_NAME = okbio-api
DB_FILE = database.db

# Go commands
help:
	@echo "make build - build the application"
	@echo "make run - run the application"
	@echo "make clean - remove the binary file"
	@echo "make test - run tests"
	@echo "make fmt - format the code"
	@echo "make lint - run linter"
	@echo "make docker-run - run the application in docker"
	@echo "make docker-stop - stop the application in docker"
	@echo "make logs - show logs"
build:
	go build -o $(APP_NAME) ./cmd/api/

run: build
	./$(APP_NAME)

clean:
	rm -f $(APP_NAME) $(DB_FILE)

test:
	go test ./... -v

fmt:
	go fmt ./...

lint:
	golangci-lint run

docker-run:
	docker compose up --build -d

docker-stop:
	docker compose down
logs:
	docker compose logs -f
.PHONY: build run clean test fmt lint docker-build docker-run docker-compose-up docker-compose-down
