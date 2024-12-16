FROM golang:1.22.5-alpine AS builder

# Установка необходимых инструментов
RUN apk add  --no-cache git make curl bash gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка основного приложенияzw
RUN go build -ldflags="-s -w" -o /main ./cmd/main.go

# Финальный образ
FROM alpine:3.18

# Установка зависимостей для работы приложения
RUN apk add --no-cache ca-certificates curl bash make pandoc


# Копирование приложения, окружения и makefile
COPY --from=builder /main /app/main

WORKDIR /app

# Установка точки входа: выполнение makefile и запуск приложения
ENTRYPOINT ["sh", "-c", "/app/main || exit 1"]