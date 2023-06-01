#
# Project-specific variables
#
# Service name. Used for binary name, docker-compose service name, etc...
SERVICE=fasttrack-service
# Enable Go Modules.
GO111MODULE=on

#
# General variables
#
# Path to Docker file
PATH_DOCKER_FILE=$(realpath ./docker/Dockerfile)
# Path to docker-compose file
PATH_DOCKER_COMPOSE_FILE=$(realpath ./docker/docker-compose.yml)
# Docker compose starting options.
DOCKER_COMPOSE_OPTIONS= -f $(PATH_DOCKER_COMPOSE_FILE)

#
# Go targets.
#
.PHONY: go-get
go-get: ## get go modules.
	@echo '>>> Getting go modules.'
	@go mod download

.PHONY: go-build
go-build: ## build service binary.
	@echo '>>> Building go binary.'
	@go build -ldflags="-s -w" -o $(SERVICE) ./main.go

#
# Tests targets.
#
.PHONY: tests-run-unit
tests-run-unit: ## run unit tests.
	@echo ">>> Running unit tests."
	go test -v ./...

.PHONY: tests-run-integration
tests-run-integration: ## run integration tests.
	@echo ">>> Running integration tests."
	go test -v -p 1 -tags="integration" ./tests/integration/golang/...


#
# Service test targets
#
.PHONY: service-build
service-build: ## build service and all it's dependencies
	@docker-compose $(DOCKER_COMPOSE_OPTIONS) build --no-cache

.PHONY: start-service-dependencies
service-start-dependencies: ## start service dependencies in docker.
	@echo ">>> Start all Service dependencies."
	@docker-compose $(DOCKER_COMPOSE_OPTIONS) up \
	-d \
	fasttrack-postgres

.PHONY: service-start
service-start: service-build service-start-dependencies ## start service in docker.
	@echo ">>> Sleeping 5 seconds until dependencies start."
	@sleep 5
	@echo ">>> Starting service."
	@echo ">>> Starting up service container."
	@docker-compose $(DOCKER_COMPOSE_OPTIONS) up -d $(SERVICE)

.PHONY: service-stop
service-stop: ## stop service in docker.
	@echo ">>> Stopping service."
	@docker-compose $(DOCKER_COMPOSE_OPTIONS) down -v --remove-orphans

.PHONY: service-restart
service-restart: service-stop service-start ## restart service in docker

.PHONY: service-test
service-test: service-stop service-start ## run tests over the service in docker.
	@echo ">>> Running tests over service."
	@docker-compose $(DOCKER_COMPOSE_OPTIONS) \
		run fasttrack-integration-tests

#
# Mockery targets.
#
.PHONY: mocks-clean
mocks-clean: ## cleans old mocks.
	find . -name "mock_*.go" -type f -print0 | xargs -0 /bin/rm -f

.PHONY: mocks-generate
mocks-generate: mocks-clean ## generate mock based on all project interfaces.
	mockery --all --dir "./pkg/api/mlflow" --inpackage --case underscore

.PHONY: help
help: ## display this help
	@ echo "Please use \`make <target>' where <target> is one of:"
	@ echo
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-10s\033[0m - %s\n", $$1, $$2}'
	@ echo
