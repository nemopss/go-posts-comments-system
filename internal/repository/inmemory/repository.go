package inmemory

import (
	"errors"

	"github.com/google/uuid"
	"github.com/nemopss/go-posts-comments-system/internal/models"
)

type InMemoryRepository struct {
	posts    map[string]*models.Post
	comments map[string]*models.Comment
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		posts:    make(map[string]*models.Post),
		comments: make(map[string]*models.Comment),
	}
}

func (repo *InMemoryRepository) GetPosts() ([]*models.Post, error) {
	posts := []*models.Post{}
	for _, post := range repo.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *InMemoryRepository) GetPost(id string) (*models.Post, error) {
	post, ok := repo.posts[id]
	if !ok {
		return nil, errors.New("Post not found")
	}
	return post, nil
}

func (repo *InMemoryRepository) CreatePost(title, content string, commentsDisabled bool) (*models.Post, error) {
	id := uuid.New().String()
	post := &models.Post{
		ID:               id,
		Title:            title,
		Content:          content,
		CommentsDisabled: commentsDisabled,
	}
	repo.posts[id] = post
	return post, nil
}

func (repo *InMemoryRepository) CreateComment(postId, parentId, content string) (*models.Comment, error) {
	id := uuid.New().String()
	comment := &models.Comment{
		ID:       id,
		PostID:   postId,
		ParentID: &parentId,
		Content:  content,
	}
	repo.comments[id] = comment
	return comment, nil
}

// GetCommentsByPostID возвращает список комментариев для указанного поста
func (repo *InMemoryRepository) GetCommentsByPostID(postId string) ([]*models.Comment, error) {
	comments := []*models.Comment{}
	for _, comment := range repo.comments {
		if comment.PostID == postId && (comment.ParentID == nil || *comment.ParentID == "") {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

// GetCommentsByParentID возвращает список дочерних комментариев для указанного комментария
func (repo *InMemoryRepository) GetCommentsByParentID(parentId string) ([]*models.Comment, error) {
	comments := []*models.Comment{}
	for _, comment := range repo.comments {
		if comment.ParentID != nil && *comment.ParentID == parentId {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}
