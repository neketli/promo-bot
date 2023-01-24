.PHONY: build
build:
	go build -v ./cmd/promo-bot

.PHONY: install
install:
	go mod download

.PHONY: docker-build
docker-build:
	docker build -t "promo-bot" .

.PHONY: docker-run
docker-run:
	docker run --env-file .env promo-bot

.DEFAULT_GOAL := build