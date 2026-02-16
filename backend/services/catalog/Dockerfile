# 👇 ИЗМЕНЕНИЕ ЗДЕСЬ: меняем 1.22 на 1.24
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

# Копируем папку с библиотеками (чтобы не качать из сети)
COPY vendor ./vendor

COPY . .

# Собираем из локальных библиотек
RUN go build -mod=vendor -o main cmd/app/main.go


# --- Вторая часть без изменений ---
FROM alpine:3.20

WORKDIR /root/

RUN apk --no-cache add ca-certificates

# ... (остальное без изменений) ...
COPY --from=builder /app/main .
COPY .env .

COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./main"]
