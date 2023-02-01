FROM golang:1.18-alpine as builder
RUN apk --no-cache add ca-certificates git build-base
WORKDIR /app

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN go build -v ./cmd/promo-bot 

EXPOSE 8080
CMD ["./promo-bot"]