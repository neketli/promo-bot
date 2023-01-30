.PHONY: build
build:
	go build -v ./cmd/promo-bot

.PHONY: install
install:
	go mod download

.PHONY: docker-build
docker-build:
	docker compose build
.PHONY: db-init


.DEFAULT_GOAL := build