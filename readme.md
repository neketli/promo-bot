# Promo bot

[![Build and Deploy](https://github.com/neketli/promo-bot/actions/workflows/deploy.yml/badge.svg?branch=master)](https://github.com/neketli/promo-bot/actions/workflows/deploy.yml) [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/neketli/promo-bot.svg)](https://github.com/neketli/promo-bot)

Телеграм бот для модерирования каналов с промокодами релаизован на основе [GO telegram bot api](github.com/go-telegram-bot-api)

## Для запуска

Настройте `.env` файл и сделайте

```bash
make install
make build
./promo-bot
```

___
Или можно использовать Docker

```bash
docker build -t promo-bot github.com/neketli/promo-bot
docker run promo-bot
```
