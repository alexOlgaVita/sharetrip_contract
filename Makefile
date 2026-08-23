# Переменные
GO := go
GO_PKG := ./...
COVERAGE_PROFILE := coverage.out
GOOSE_MIGRATION_DIR := migrations
NAME_CONTRACT := create_table_contract
NAME_CONTRACT_SERVICE := create_table_contract_service
DOCKER_NAME := pg
BINARY_NAME := sharetripContract.exe
MAIN_PKG := ./cmd/sharetripContract
BIN_DIR := bin

# Цель по умолчанию
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  test        - Run all tests"
	@echo "  coverage    - Run tests and generate HTML coverage report"
	@echo "  cover       - Alias for coverage"
	@echo "  lint        - Run golangci-lint"
	@echo "  all         - Run lint, tests and coverage"
	@echo "  help        - Show this help"

# Запуск всех тестов
.PHONY: test
test:
	$(GO) test -v $(GO_PKG) -coverprofile=$(COVERAGE_PROFILE) $(GO_PKG)

# Генерация отчёта о покрытии в формате HTML
.PHONY: coverage cover
coverage cover:
	$(GO) test -coverprofile=$(COVERAGE_PROFILE) $(GO_PKG)
	$(GO) tool cover -html=$(COVERAGE_PROFILE) -o coverage.html
	@echo "Coverage report generated: %cd%\coverage.html"

# Вывод покрытия в терминал (опционально)
.PHONY: cover-report
cover-report:
	$(GO) test -cover $(GO_PKG)

# Запуск всех тестов
.PHONY: fmt
fmt:
	$(GO) fmt $(GO_PKG)

# Проверка кода с помощью golangci-lint
.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build:
	$(GO) build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PKG)

.PHONY: run
run:
	$(GO) run ./$(MAIN_PKG)/main.go

.PHONY: deps
deps:
	$(GO) mod tidy
	$(GO) mod download

.PHONY: migrate-status
migrate-status:
	goose -dir $(GOOSE_MIGRATION_DIR) postgres "postgres://postgres:password@localhost:6545/sharetripContracts?sslmode=disable" status

.PHONY: migrate-up
migrate-up:
	goose -dir $(GOOSE_MIGRATION_DIR) postgres "postgres://postgres:password@localhost:6545/sharetripContracts?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	goose -dir $(GOOSE_MIGRATION_DIR) postgres "postgres://postgres:password@localhost:6545/sharetripContracts?sslmode=disable" down

.PHONY: migrate-create
migrate-create:
	goose -dir $(GOOSE_MIGRATION_DIR) create -s $(NAME_CONTRACT) sql
	goose -dir $(GOOSE_MIGRATION_DIR) create -s $(NAME_CONTRACT_SERVICE) sql

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down --volumes

# Запуск всех проверок
.PHONY: check
check: fmt lint test coverage

.PHONY: e2eCheck
e2eCheck:
	curl http://localhost:8080/api/ready