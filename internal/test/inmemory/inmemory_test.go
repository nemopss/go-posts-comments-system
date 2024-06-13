package test

import (
	"log"
	"testing"

	"github.com/nemopss/go-posts-comments-system/internal/models"
	"github.com/nemopss/go-posts-comments-system/internal/repository/inmemory"
)

// Utility functions to create posts and comments
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

// Test for GetPosts
func TestGetPosts_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Create some posts
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

// Test for GetPost
func TestGetPost_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Create a post
	post := createPost(repo, "Test Post", "This is a test post", false)

	fetchedPost, err := repo.GetPost(post.ID)
	if err != nil {
		t.Errorf("failed to get post: %v", err)
	}
	if fetchedPost.ID != post.ID {
		t.Errorf("expected post ID %s, got %s", post.ID, fetchedPost.ID)
	}
}

// Test for CreatePost
func TestCreatePost_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Create a post
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

// Test for CreateComment
func TestCreateComment_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Create a post and a comment
	post := createPost(repo, "Test Post", "This is a test post", false)
	comment := createComment(repo, post.ID, "", "This is a test comment")

	if comment.Content != "This is a test comment" {
		t.Errorf("expected content 'This is a test comment', got %s", comment.Content)
	}
	if *comment.ParentID != "" {
		t.Errorf("expected parentId to be empty, got %s", *comment.ParentID)
	}
}

// Test for GetCommentsByPostID
func TestGetCommentsByPostID_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Create a post and some comments
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

// Test for GetCommentsByParentID with pagination
func TestGetCommentsByParentID_InMemory(t *testing.T) {
	repo := inmemory.NewInMemoryRepository()

	// Create a post and some nested comments
	post := createPost(repo, "Test Post", "This is a test post", false)
	comment1 := createComment(repo, post.ID, "", "Comment 1")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 2")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 3")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 4")
	_ = createComment(repo, post.ID, comment1.ID, "Comment 5")

	// Check fetching comments with pagination
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
