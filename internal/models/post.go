package models

import "time"

// Post представляет собой структуру поста
type Post struct {
	ID               string     // Уникальный идентификатор поста
	Title            string     // Заголовок поста
	Content          string     // Содержимое поста
	Comments         []*Comment // Список комментариев к посту
	CommentsDisabled bool       // Флаг, указывающий, отключены ли комментарии к посту
	CreatedAt        time.Time  // Время создания поста
}
