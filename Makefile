# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v ./...

## audit: tidy dependencies, format and vet code
.PHONY: audit
audit:
	@echo 'Tidying dependencies...'
	go mod tidy
	@echo 'Formatting code...'
	go fmt ./...
	golines . -w
	@echo 'Vetting code...'
	go vet ./...

.PHONY: vendor
vendor :
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api/docker-compose: run the app using docker compose
.PHONY: run/api/docker-compose
run/api/docker-compose:
	@echo 'Formatting code...'
	go fmt ./...
	golines . -w

	@echo 'Vetting code...'
	go vet ./...

	@echo 'Running tests...'
	go test -race -vet=off ./...

	@echo 'Creating swagger documentation...'
	swag init -g ./cmd/api/main.go

	@echo 'Running docker compose...'
	docker compose up --build

## run/rabbitmq/docker-compose
.PHONY: run/rabbitmq/docker-compose
run/rabbitmq/docker-compose:
	@echo 'Starting RabbitMQ'
	docker compose --profile rabbitmq up

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:

	@echo 'Creating swagger documentation...'
	swag init -g ./cmd/api/main.go

	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

##  run/bin/api:
.PHONY: run/bin/api
run/bin/api:
	@echo 'Creating swagger documentation...'
	swag init -g ./cmd/api/main.go

	@echo 'Running cmd/api...'
	./bin/api
