# Используем официальный образ Golang 1.24.3 для сборки
FROM golang:1.24.3 AS build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum в контейнер
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем все исходники проекта
COPY . .

# Собираем приложение
RUN go build -o app ./cmd/app

# Используем более новый образ с обновленной версией glibc (Debian Bookworm)
FROM debian:bookworm-slim

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y libc6

# Копируем скомпилированное приложение
COPY --from=build /app/app /usr/local/bin/app

# Копируем файл конфигурации (если нужно)
COPY config/config.yml /config/config.yml

# Порт, на котором будет слушать ваше приложение
EXPOSE 8080

# Команда для запуска приложения
CMD ["/usr/local/bin/app"]
