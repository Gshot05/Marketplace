# Используем официальный образ Go
FROM golang:1.23

# Установим git (нужен для go get)
RUN apt-get update && apt-get install -y git

# Рабочая директория
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Экспонируем порт приложения
EXPOSE 8080

# Запускаем сервер
CMD ["go", "run", "./cmd/server"]
