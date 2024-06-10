package server

import (
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/nemopss/go-posts-comments-system/internal/gql"
	"github.com/nemopss/go-posts-comments-system/internal/repository"
)

// Server представляет сервер GraphQL
type Server struct {
	schema *graphql.Schema
}

// NewServer создает новый экземпляр Server
func NewServer(repo repository.Repository) *Server {
	schema, err := gql.NewSchema(repo)
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}
	log.Println("Starting server...")
	return &Server{schema: &schema}
}

// Handler возвращает обработчик HTTP для GraphQL запросов
func (s *Server) Handler() http.Handler {
	log.Println("Handling server...")
	return handler.New(&handler.Config{
		Schema:   s.schema,
		Pretty:   true,
		GraphiQL: true,
	})
}
