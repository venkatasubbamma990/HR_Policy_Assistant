APP_NAME := hrpolicy
BINARY   := bin/$(APP_NAME)
MAIN     := ./cmd/hrpolicy
INGEST   := ./cmd/ingest
INDEX    := ./cmd/index
ASK      := ./cmd/ask
COMPOSE  := docker compose

# Use the locally installed Go toolchain (avoids auto-download in WSL/offline envs)
export GOTOOLCHAIN := local

.PHONY: build ingest index ask test up down docker-build clean check-llm setup-wsl-llm

build:
	go build -o $(BINARY) $(MAIN)

ingest:
	go run $(INGEST)

index:
	go run $(INDEX)

ask:
	go run $(ASK) -q "$(Q)"

check-llm:
	bash scripts/check-llm.sh

setup-wsl-llm:
	powershell.exe -ExecutionPolicy Bypass -File scripts/setup-wsl-llm.ps1

test:
	go test ./...

docker-build:
	$(COMPOSE) build

up: docker-build
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

clean:
	rm -rf bin/
	$(COMPOSE) down --rmi local -v 2>/dev/null || true
