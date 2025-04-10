ROOT 			:=	$(shell pwd)
BUILD_DIR 		:=	$(ROOT)/build
BIN_NAME		:=	myapp
LOG_FILE		:= 	$(ROOT)/log
CGO_ENABLED 	?=	0
LDFLAGS     	:= -w -s
PID 			:= $(shell ps -ef | grep $(BIN_NAME) | grep -v grep | awk '{print $$2}')

.DEFAULT_GOAL := restart

.PHONY: restart
restart: stop run

.PHONY: run
run:
	@mkdir -p $(LOG_FILE)
	nohup $(BUILD_DIR)/$(BIN_NAME) > $(LOG_FILE)/runtime.log 2>&1 &

.PHONY: stop
stop:
ifneq ($(strip $(PID)),)
	@echo "Process PID: $(PID)"
	kill $(PID)
endif

.PHONY: build
build:
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) \
	go build -ldflags='$(LDFLAGS)' -o '$(BUILD_DIR)/$(BIN_NAME)' $(ROOT)/cmd/app/

.PHONY: test
test:
	go test ./test/...