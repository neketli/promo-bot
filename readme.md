# promo-bot

Implementation of simple calculator bot in telegram

## Build and run app

```bash
make install
make build
./promo-bot
```

### If you use docker

```bash
make docker-build # build an image
# remember to setup .env
make docker-run # runs app
```

---

> Написать сервис придерживаясь принципов чистой архитектуры. В сервисе нужно реализовать интеграцию с telegram ботом, который умеет вычислять простейшие арифметические операции (+,-,\*,:). При добавлении бота в telegram, сохранять отдаваемую информацию о пользователе (можно использовать любую БД). Для показа количества пользователей, в сервисе добавить поддержку HTTP-сервера, в котором будет единственный метод GET /info
