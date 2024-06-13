package inmemory

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/nemopss/go-posts-comments-system/internal/models"
)

// InMemoryRepository представляет репозиторий, хранящий данные в памяти.
type InMemoryRepository struct {
	posts    map[string]*models.Post
	comments map[string]*models.Comment
}

// NewInMemoryRepository создает новый репозиторий в памяти.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		posts:    make(map[string]*models.Post),
		comments: make(map[string]*models.Comment),
	}
}

// GetPosts возвращает все посты из репозитория.
func (repo *InMemoryRepository) GetPosts() ([]*models.Post, error) {
	log.Println("Querying posts...")
	posts := []*models.Post{}
	for _, post := range repo.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

// GetPost возвращает пост по его ID. Если пост не найден, возвращает ошибку.
func (repo *InMemoryRepository) GetPost(id string) (*models.Post, error) {
	log.Println("Querying post with ID:", id)
	post, ok := repo.posts[id]
	if !ok {
		return nil, errors.New("Post not found")
	}
	return post, nil
}

// CreatePost создает новый пост и добавляет его в репозиторий.
func (repo *InMemoryRepository) CreatePost(title, content string, commentsDisabled bool) (*models.Post, error) {
	id := uuid.New().String()
	log.Println("Creating post with ID:", id)
	createdAt := time.Now()
	post := &models.Post{
		ID:               id,
		Title:            title,
		Content:          content,
		CommentsDisabled: commentsDisabled,
		CreatedAt:        createdAt,
	}
	repo.posts[id] = post
	return post, nil
}

// CreateComment создает новый комментарий и добавляет его в репозиторий.
func (repo *InMemoryRepository) CreateComment(postId, parentId, content string) (*models.Comment, error) {
	if repo.posts[postId].CommentsDisabled == true {
		return nil, errors.New("Comments are disabled on this post!")
	}
	// Ограничение в 2000 символов на комментарий
	if len(content) > 2000 {
		return nil, errors.New("комментарий не может превышать 2000 символов")
	}
	id := uuid.New().String()
	log.Println("Creating comment with ID:", id)
	createdAt := time.Now()
	comment := &models.Comment{
		ID:        id,
		PostID:    postId,
		ParentID:  &parentId,
		Content:   content,
		CreatedAt: createdAt,
	}
	repo.comments[id] = comment
	if parentId != "" {
		repo.comments[parentId].Children = append(repo.comments[parentId].Children, comment)
	}
	return comment, nil
}

// GetCommentsByPostID возвращает список комментариев для указанного поста с пагинацией
func (repo *InMemoryRepository) GetCommentsByPostID(postId string, first int64, after *string) ([]*models.Comment, error) {
	log.Println("Getting comments on post with ID:", postId)
	var allComments []*models.Comment

	comments := []*models.Comment{}
	for _, comment := range repo.comments {
		if comment.PostID == postId && (comment.ParentID == nil || *comment.ParentID == "") {
			comments = append(comments, comment)
		}
	}

	// Сортируем комментарии по времени создания
	sort.Slice(allComments, func(i, j int) bool {
		return allComments[i].CreatedAt.Before(allComments[j].CreatedAt)
	})

	// Применяем пагинацию
	startIndex := 0
	if after != nil {
		for i, comment := range comments {
			if comment.ID > *after {
				startIndex = i + 1
				break
			}
		}
	}

	endIndex := int64(startIndex) + first
	if endIndex > int64(len(comments)) {
		endIndex = int64(len(comments))
	}

	return comments[startIndex:endIndex], nil

}

// GetCommentsByParentID возвращает список дочерних комментариев для указанного комментария с пагинацией
func (repo *InMemoryRepository) GetCommentsByParentID(parentId string, first int64, after *string) ([]*models.Comment, error) {
	log.Println("Getting comments from parent with ID:", parentId)
	parentComment, ok := repo.comments[parentId]
	if !ok {
		return nil, fmt.Errorf("comment with id %s not found", parentId)
	}

	startIndex := 0
	if after != nil {
		for i, comment := range parentComment.Children {
			if comment.ID > *after {
				startIndex = i + 1
				break
			}
		}
	}

	endIndex := int64(startIndex) + first
	if int(endIndex) > len(parentComment.Children) {
		endIndex = int64(len(parentComment.Children))
	}

	return parentComment.Children[startIndex:endIndex], nil
}

// DeletePost удаляет пост по его ID
func (repo *InMemoryRepository) DeletePost(id string) error {
	log.Println("Deleting post with ID:", id)
	_, ok := repo.posts[id]
	if !ok {
		return errors.New("Post not found")
	}

	// Удаляем все комментарии к этому посту
	for commentID, comment := range repo.comments {
		if comment.PostID == id {
			delete(repo.comments, commentID)
		}
	}

	// Теперь удаляем сам пост
	delete(repo.posts, id)
	return nil
}

// DeleteComment удаляет комментарий по его ID
func (repo *InMemoryRepository) DeleteComment(id string) error {
	log.Println("Deleting comment with ID:", id)
	comment, ok := repo.comments[id]
	if !ok {
		return errors.New("Comment not found")
	}

	// Удаляем связи этого комментария с дочерними комментариями
	for _, childComment := range repo.comments {
		if childComment.ParentID != nil && *childComment.ParentID == id {
			delete(repo.comments, childComment.ID)
		}
	}

	// Теперь удаляем сам комментарий
	delete(repo.comments, id)

	// Если у комментария был родитель, обновляем его список детей
	if comment.ParentID != nil {
		parentComment, ok := repo.comments[*comment.ParentID]
		if ok {
			for i, child := range parentComment.Children {
				if child.ID == id {
					parentComment.Children = append(parentComment.Children[:i], parentComment.Children[i+1:]...)
					break
				}
			}
		}
	}

	return nil
}
