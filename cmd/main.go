package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/nemopss/go-posts-comments-system/internal/repository/inmemory"
	"github.com/nemopss/go-posts-comments-system/internal/server"
)

func main() {
	rep := inmemory.NewInMemoryRepository()

	srv := server.NewServer(rep)
	http.Handle("/graphql", srv.Handler())

	log.Println("Server is running on http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
