package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/nemopss/go-posts-comments-system/internal/repository"
	"github.com/nemopss/go-posts-comments-system/internal/repository/inmemory"
	"github.com/nemopss/go-posts-comments-system/internal/repository/postgres"
	"github.com/nemopss/go-posts-comments-system/internal/server"
)

func main() {
	// Определение флага командной строки для выбора хранилища. true -> Postgres, false -> in-memory storage
	postgresStorageFlag := flag.Bool("postgres", false, "A flag to determine which storage you want to use. true -> Postgres, false -> in-memory storage")
	flag.Parse()
	var rep repository.Repository
	// Если флаг postgresStorageFlag установлен в false, используем хранилище в памяти
	if !*postgresStorageFlag {
		rep = inmemory.NewInMemoryRepository()
		log.Println("In-memory storage active...")
	} else {
		// Формирование строки подключения к базе данных PostgreSQL с использованием переменных окружения
		connStr := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOST") + ":5432/" + os.Getenv("POSTGRES_DB") + "?sslmode=disable"

		// Открытие соединения с базой данных PostgreSQL
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Error opening db: %v", err)
		}
		// Инициализация репозитория с использованием PostgreSQL
		rep = postgres.NewPostgresRepository(db)
		log.Println("PostgreSQL storage active...")
	}

	// Создание нового сервера GraphQL
	srv := server.NewServer(rep)

	// Регистрация обработчика для маршрута /graphql
	http.Handle("/graphql", srv.Handler())

	log.Println("Server is running on http://localhost:8080/graphql")

	// Запуск HTTP сервера на порту 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
