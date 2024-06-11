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
	log.Println("GetPosts called")
	posts := []*models.Post{}
	for _, post := range repo.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

// GetPost возвращает пост по его ID. Если пост не найден, возвращает ошибку.
func (repo *InMemoryRepository) GetPost(id string) (*models.Post, error) {
	post, ok := repo.posts[id]
	if !ok {
		return nil, errors.New("Post not found")
	}
	return post, nil
}

// CreatePost создает новый пост и добавляет его в репозиторий.
func (repo *InMemoryRepository) CreatePost(title, content string, commentsDisabled bool) (*models.Post, error) {
	log.Println("CreatePost called")
	id := uuid.New().String()
	createdAt := time.Now()
	post := &models.Post{
		ID:               id,
		Title:            title,
		Content:          content,
		CommentsDisabled: commentsDisabled,
		CreatedAt:        createdAt,
	}
	repo.posts[id] = post
	log.Println("Created post with ID:", id)
	return post, nil
}

// CreateComment создает новый комментарий и добавляет его в репозиторий.
func (repo *InMemoryRepository) CreateComment(postId, parentId, content string) (*models.Comment, error) {
	id := uuid.New().String()
	log.Printf("Created comment with ID: %v\n", id)
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
	postComments := []*models.Comment{}
	for _, comment := range repo.comments {
		if comment.PostID == postId && comment.ParentID == nil {
			postComments = append(postComments, comment)
		}
	}

	// Сортируем комментарии по времени создания
	sort.Slice(postComments, func(i, j int) bool {
		return postComments[i].CreatedAt.Before(postComments[j].CreatedAt)
	})

	// Применяем пагинацию
	startIndex := 0
	if after != nil {
		for i, comment := range postComments {
			if comment.ID > *after {
				startIndex = i + 1
				break
			}
		}
	}

	endIndex := int64(startIndex) + first
	if endIndex > int64(len(postComments)) {
		endIndex = int64(len(postComments))
	}

	return postComments[startIndex:endIndex], nil
}

// GetCommentsByParentID возвращает список дочерних комментариев для указанного комментария с пагинацией
func (repo *InMemoryRepository) GetCommentsByParentID(parentId string, first int64, after *string) ([]*models.Comment, error) {
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

	log.Printf("Range of children: %v", len(parentComment.Children))
	log.Printf("Returning comments by parent ID: %v. SI: %v, EI: %v, first: %v\n", parentId, startIndex, endIndex, first)
	for i, comm := range parentComment.Children[startIndex:endIndex] {
		log.Printf("#%v Comment: %v", i, comm.Content)
	}
	return parentComment.Children[startIndex:endIndex], nil
}

// DeletePost удаляет пост по его ID
func (repo *InMemoryRepository) DeletePost(id string) error {
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
	comment, ok := repo.comments[id]
	if !ok {
		return errors.New("Comment not found")
	}

	// Удаляем связи этого комментария с дочерними комментариями
	for _, childComment := range repo.comments {
		if childComment.ParentID != nil && *childComment.ParentID == id {
			childComment.ParentID = nil
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
