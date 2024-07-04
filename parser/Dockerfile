# Используем образ Golang как базовый образ для сборки
FROM golang:1.20 AS builder

WORKDIR /app

# Копируем зависимости Go модулей
#COPY go.mod .
#COPY go.sum .
#COPY config.json .

# Устанавливаем зависимости Go
#RUN go mod download

# Копируем остальные файлы в контейнер (включая main.go и другие)
COPY . .

# Собираем Go приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o parser main.go

# Финальный образ без лишних зависимостей
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app .

# CMD устанавливается внутри Docker Compose
CMD ["./parser"]
