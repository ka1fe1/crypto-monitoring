BINARY_NAME=crypto-monitoring
MAIN_PATH=cmd/server/main.go

.PHONY: all build clean run deps help fmt sync-config swagger

all: fmt deps build

sync-config:
	go run cmd/sync-config/main.go

fmt: sync-config
	go fmt ./...

deps:
	go mod download
	go mod tidy
	go mod vendor

build: deps swagger
	mkdir -p bin
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

run: build
	./bin/$(BINARY_NAME)

clean:
	go clean
	rm -rf bin/

swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  deps         Download dependencies"
	@echo "  build        Build the project"
	@echo "  run          Run the project"
	@echo "  clean        Clean build artifacts"
	@echo "  fmt          Format the code"
	@echo "  sync-config  Generate config.yaml.temp from config.yaml"
	@echo "  swagger      Generate swagger documentation"
