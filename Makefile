ROOT        := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BUILD_DIR   := $(ROOT)build
BIN_NAME    ?= quick-gin
BINARY      := $(BUILD_DIR)/$(BIN_NAME)
LOG_DIR     := $(ROOT)log
LOG_FILE    := $(LOG_DIR)/runtime.log
PID_FILE    := $(LOG_DIR)/$(BIN_NAME).pid
GO          ?= go
CGO_ENABLED ?= 0
LDFLAGS     ?= -w -s

.DEFAULT_GOAL := help

.PHONY: help
help:
	@printf "Available targets:\n"
	@printf "  make build    - Compile the application binary\n"
	@printf "  make run      - Build and start the application in the background\n"
	@printf "  make stop     - Stop the background application process\n"
	@printf "  make restart  - Restart the application\n"
	@printf "  make test     - Run Go unit tests\n"
	@printf "  make tidy     - Run go mod tidy\n"
	@printf "  make clean    - Remove build artifacts and PID file\n"

.PHONY: all
all: build

.PHONY: build
build: $(BINARY)

$(BINARY):
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -ldflags '$(LDFLAGS)' -o $(BINARY) ./cmd/app

.PHONY: run
run: build
	@mkdir -p $(LOG_DIR)
	@echo "Starting $(BIN_NAME)..."
	nohup $(BINARY) > $(LOG_FILE) 2>&1 & echo $$! > $(PID_FILE)

.PHONY: stop
stop:
	@if [ -f $(PID_FILE) ]; then \
		PID=$$(cat $(PID_FILE)); \
		if ps -p $$PID > /dev/null 2>&1; then \
			echo "Stopping $(BIN_NAME) ($$PID)"; \
			kill $$PID; \
		else \
			echo "Process $$PID not running"; \
		fi; \
		rm -f $(PID_FILE); \
	else \
		echo "No PID file found for $(BIN_NAME)."; \
	fi

.PHONY: restart
restart: stop run

.PHONY: test
test:
	$(GO) test ./...

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)
	@rm -f $(PID_FILE)
