package test

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // импорт драйвера PostgreSQL
	"github.com/nemopss/go-posts-comments-system/internal/models"
	"github.com/nemopss/go-posts-comments-system/internal/repository/postgres"
)

func createPost(repo *postgres.PostgresRepository, title, content string, commentsDisabled bool) *models.Post {
	post, err := repo.CreatePost(title, content, commentsDisabled)
	if err != nil {
		panic(err)
	}
	return post
}

func createComment(repo *postgres.PostgresRepository, postId, parentId, content string) *models.Comment {
	comment, err := repo.CreateComment(postId, parentId, content)
	if err != nil {
		panic(err)
	}
	return comment
}

func TestGetCommentsByPostID(t *testing.T) {
	connStr := "postgres://gosuper:Ukflbkby2004@localhost:5432/go-posts-comments-db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresRepository(db)

	// Создаем пост и несколько комментариев
	post := createPost(repo, "Test Post", "This is a test post", false)
	createComment(repo, post.ID, "", "Comment 1")
	createComment(repo, post.ID, "", "Comment 2")
	createComment(repo, post.ID, "", "Comment 3")
	createComment(repo, post.ID, "", "Comment 4")
	createComment(repo, post.ID, "", "Comment 5")

	// Проверяем получение комментариев с пагинацией
	comments, err := repo.GetCommentsByPostID(post.ID, 2, nil)
	if err != nil {
		t.Errorf("failed to get comments: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}

	after := comments[len(comments)-1].ID
	comments, err = repo.GetCommentsByPostID(post.ID, 2, &after)
	if err != nil {
		t.Errorf("failed to get comments: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}
}
