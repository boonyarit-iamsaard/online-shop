APP := api
BIN := ./bin/$(APP)
COMPOSE_LOCAL := docker compose
COMPOSE_INFRA := docker compose -f compose.yml

ifeq ($(BUILD),1)
UP_FLAGS := --build
endif

.PHONY: build test run dev up down up-infra down-infra

build:
	go build -o $(BIN) ./cmd/api

test:
	go test ./...

run:
	go run ./cmd/api

dev:
	air

up:
	$(COMPOSE_LOCAL) up -d --wait $(UP_FLAGS)

down:
	$(COMPOSE_LOCAL) down

up-infra:
	$(COMPOSE_INFRA) up -d --wait

down-infra:
	$(COMPOSE_INFRA) down
