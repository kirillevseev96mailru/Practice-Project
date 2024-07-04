# siteforlogs/Dockerfile
FROM golang:1.20 AS builder

WORKDIR /app

COPY New-logs.json .
# Копируем исходный код в образ
COPY . .

# Собираем Go приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o siteforlogs main.go

# Финальный образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/New-logs.json .
COPY --from=builder /app/siteforlogs .
COPY --from=builder  /app/static/ /app/static/ 
COPY --from=builder  /app/templates/ /app/templates/
CMD ["./siteforlogs"]
