# Promo bot

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
