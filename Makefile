.PHONY: build
build:
	go build -v ./cmd/promo-bot

.PHONY: install
install:
	go mod download

.DEFAULT_GOAL := build