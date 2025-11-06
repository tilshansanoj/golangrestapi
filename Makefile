APP_NAME = golang_rest_api
BIN_NAME = $(APP_NAME)_bin
VERSION = 1.0.0
BUILD_DIR = ./bin
GO_FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

run:
	@echo "Running $(APP_NAME) ..."
	@go run main.go || true

deps:
	@echo "Installing dependencies ..."
	@go mod tidy

fmt:
	@echo "Formatting Go files ..."
	@go fmt $(GO_FILES)

build:
	@echo "Building $(APP_NAME) ..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BIN_NAME) -ldflags "-X main.version=$(VERSION)" main.go
	@echo "$(APP_NAME) built successfully at $(BUILD_DIR)/$(BIN_NAME)"

clean:
	@echo "Cleaning build artifacts ..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned."

stop:
	@echo "Stopping $(APP_NAME) ..."
	@pkill -f "go run main.go" || echo "$(APP_NAME) is not running."

migrate-up:
	@echo "Applying database migrations ..."
	@dbmate up
	@echo "Database migrations applied."

migrate-down:
	@echo "Reverting database migrations ..."
	@dbmate down
	@echo "Database migrations reverted."

help:
	@echo "Makefile commands for $(APP_NAME):"
	@echo "  run            - Run the application"
	@echo "  deps           - Install dependencies"
	@echo "  fmt            - Format Go files"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  stop           - Stop the running application"
	@echo "  migrate-up     - Apply database migrations"
	@echo "  migrate-down   - Revert database migrations"
	@echo "  help           - Show this help message"
