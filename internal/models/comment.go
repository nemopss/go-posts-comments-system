package models

import "time"

// Comment представляет собой структуру комментария
type Comment struct {
	ID        string     // Уникальный идентификатор комментария
	PostID    string     // Идентификатор поста, к которому относится комментарий
	ParentID  *string    // Идентификатор родительского комментария (если есть)
	Content   string     // Содержимое комментария
	Children  []*Comment // Список дочерних комментариев
	CreatedAt time.Time  // Время создания комментария
}
