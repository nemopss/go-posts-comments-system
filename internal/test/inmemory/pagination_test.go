package test

import (
	"testing"

	"github.com/nemopss/go-posts-comments-system/internal/models"
	"github.com/nemopss/go-posts-comments-system/internal/repository/inmemory"
)

func createPost(repo *inmemory.InMemoryRepository, title, content string, commentsDisabled bool) *models.Post {
	post, err := repo.CreatePost(title, content, commentsDisabled)
	if err != nil {
		panic(err)
	}
	return post
}

func createComment(repo *inmemory.InMemoryRepository, postId, parentId, content string) *models.Comment {
	comment, err := repo.CreateComment(postId, parentId, content)
	if err != nil {
		panic(err)
	}
	return comment
}

func TestGetCommentsByParentID_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Создаем пост и несколько комментариев
	post := createPost(repo, "Test Post", "This is a test post", false)
	comment1 := createComment(repo, post.ID, "", "Comment 1")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 2")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 3")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 4")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 5")

	// Проверяем получение комментариев с пагинацией
	comments, err := repo.GetCommentsByParentID(comment1.ID, 2, nil)
	if err != nil {
		t.Errorf("failed to get comments: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}

	after := comments[len(comments)-1].ID
	comments, err = repo.GetCommentsByParentID(comment1.ID, 2, &after)
	if err != nil {
		t.Errorf("failed to get comments: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}

}
