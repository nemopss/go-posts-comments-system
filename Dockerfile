# Используем официальный образ Go
FROM golang:1.22.3

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект в рабочую директорию контейнера
COPY . .

# Устанавливаем необходимые пакеты
RUN go get -d -v ./...
RUN go install -v ./...

# Стандартная команда запуска, может быть переопределена в docker-compose.yml
CMD ["go", "run", "cmd/main.go"]
