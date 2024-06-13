package test

import (
	"log"
	"testing"

	"github.com/nemopss/go-posts-comments-system/internal/models"
	"github.com/nemopss/go-posts-comments-system/internal/repository/inmemory"
)

// Функция для создания поста в in-memory хранилище
func createPost(repo *inmemory.InMemoryRepository, title, content string, commentsDisabled bool) *models.Post {
	post, err := repo.CreatePost(title, content, commentsDisabled)
	if err != nil {
		panic(err)
	}
	return post
}

// Функция для создания комментария в in-memory хранилище
func createComment(repo *inmemory.InMemoryRepository, postId, parentId, content string) *models.Comment {
	comment, err := repo.CreateComment(postId, parentId, content)
	if err != nil {
		panic(err)
	}
	return comment
}

// Тест GetPosts
func TestGetPosts_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Создание постов
	createPost(repo, "Test Post 1", "This is the first test post", false)
	createPost(repo, "Test Post 2", "This is the second test post", true)

	posts, err := repo.GetPosts()
	if err != nil {
		t.Errorf("failed to get posts: %v", err)
	}
	if len(posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(posts))
	}
}

// Тест GetPost
func TestGetPost_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Создание поста
	post := createPost(repo, "Test Post", "This is a test post", false)

	fetchedPost, err := repo.GetPost(post.ID)
	if err != nil {
		t.Errorf("failed to get post: %v", err)
	}
	if fetchedPost.ID != post.ID {
		t.Errorf("expected post ID %s, got %s", post.ID, fetchedPost.ID)
	}
}

// Тест CreatePost
func TestCreatePost_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Создание поста
	post := createPost(repo, "Test Post", "This is a test post", false)

	if post.Title != "Test Post" {
		t.Errorf("expected title 'Test Post', got %s", post.Title)
	}
	if post.Content != "This is a test post" {
		t.Errorf("expected content 'This is a test post', got %s", post.Content)
	}
	if post.CommentsDisabled != false {
		t.Errorf("expected commentsDisabled to be false, got %v", post.CommentsDisabled)
	}
}

// Тест CreateComment
func TestCreateComment_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Создание поста и комментария к нему
	post := createPost(repo, "Test Post", "This is a test post", false)
	comment := createComment(repo, post.ID, "", "This is a test comment")

	if comment.Content != "This is a test comment" {
		t.Errorf("expected content 'This is a test comment', got %s", comment.Content)
	}
	if *comment.ParentID != "" {
		t.Errorf("expected parentId to be empty, got %s", *comment.ParentID)
	}
}

// Тест GetCommentsByPostID
func TestGetCommentsByPostID_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Создание поста и комментариев к нему
	post := createPost(repo, "Test Post", "This is a test post", false)
	createComment(repo, post.ID, "", "Comment 1")
	createComment(repo, post.ID, "", "Comment 2")

	comments, err := repo.GetCommentsByPostID(post.ID, 2, nil)
	log.Println("LEN INMEM COMM:", len(comments))
	if err != nil {
		t.Errorf("failed to get comments: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}
}

// Тест GetCommentsByParentID with с пагинацией
func TestGetCommentsByParentID_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Создание поста и нескольких вложенных комментариев
	post := createPost(repo, "Test Post", "This is a test post", false)
	comment1 := createComment(repo, post.ID, "", "Comment 1")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 2")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 3")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 4")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 5")

	// Проверка полученных комментариев
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
