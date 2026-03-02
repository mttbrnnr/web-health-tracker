.PHONY: dev build build-pi clean deploy

BINARY_NAME := health-tracker
BUILD_DIR := ./build

dev:
	go run ./cmd/health-tracker

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/health-tracker

build-pi:
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-arm64 ./cmd/health-tracker

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)

deploy:
	./deploy/deploy-pi.sh
