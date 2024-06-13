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
	postgresStorageFlag := flag.Bool("postgres", false, "A flag to determine which storage you want to use. true -> Postgres, false -> in-memory storage")
	flag.Parse()
	var rep repository.Repository
	if !*postgresStorageFlag {
		rep = inmemory.NewInMemoryRepository()
	} else {
		//connStr := "postgres://gosuper:Ukflbkby2004@localhost:5432/go-posts-comments-db?sslmode=disable"
		connStr := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOST") + ":5432/" + os.Getenv("POSTGRES_DB") + "?sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Error opening db: %v", err)
		}
		rep = postgres.NewPostgresRepository(db)

	}
	srv := server.NewServer(rep)
	http.Handle("/graphql", srv.Handler())

	log.Println("Server is running on http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
