package repository

import "github.com/nemopss/go-posts-comments-system/internal/models"

// Repository представляет интерфейс для работы с постами и комментариями
// и позволяет абстрагироваться от конкретной реализации хранилища данных
// будь то Postgres или in-memory хранилище
type Repository interface {

	// GetPosts возвращает список всех постов
	// Возвращается слайс указателей на модели Post и ошибку в случае неудачи
	GetPosts() ([]*models.Post, error)

	// GetPost возвращает пост по его идентификатору uuid
	// Принимает строковый идентификатор поста и возвращает указатель на модель Post и ошибку в случае неудачи
	GetPost(id string) (*models.Post, error)

	// CreatePost создаёт новый пост
	// Принимает заголовок (title), содержание (content) и флаг отключения комментариев (commentsDisabled).
	// Возвращает указатель на созданную модель Post и ошибку в случае неудачи.
	CreatePost(title, content string, commentsDisabled bool) (*models.Post, error)

	// CreateComment создает новый комментарий к посту.
	// Принимает идентификатор поста (postId), идентификатор родительского комментария (parentId)
	// и содержание комментария (content). ParentID может быть пустым, если комментарий не является ответом.
	// Возвращает указатель на созданную модель Comment и ошибку в случае неудачи.
	CreateComment(postId, parentId, content string) (*models.Comment, error)

	// GetCommentsByPostID возвращает список комментариев для указанного поста.
	GetCommentsByPostID(postId string) ([]*models.Comment, error)

	// GetCommentsByParentID возвращает список дочерних комментариев для указанного комментария.
	GetCommentsByParentID(parentId string) ([]*models.Comment, error)

	DeletePost(id string) error

	DeleteComment(id string) error
}
