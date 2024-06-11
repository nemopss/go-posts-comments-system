package postgres

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nemopss/go-posts-comments-system/internal/models"
)

// PostgresRepository представляет собой хранилище данных в PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository создает новый экземпляр PostgresRepository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetPosts возвращает список всех постов
func (repo *PostgresRepository) GetPosts() ([]*models.Post, error) {
	rows, err := repo.db.Query("SELECT id, title, content, comments_disabled, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Получаем комментарии для данного поста
		comments, err := repo.GetCommentsByPostID(post.ID, 0, nil)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		posts = append(posts, post)
	}

	return posts, nil
}

// GetPost возвращает пост по его ID
func (repo *PostgresRepository) GetPost(id string) (*models.Post, error) {
	row := repo.db.QueryRow("SELECT id, title, content, comments_disabled FROM posts WHERE id=$1", id)
	post := &models.Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// CreatePost создает новый пост
func (repo *PostgresRepository) CreatePost(title, content string, commentsDisabled bool) (*models.Post, error) {
	id := uuid.New().String()
	createdAt := time.Now()
	_, err := repo.db.Exec("INSERT INTO posts (id, title, content, comments_disabled, created_at) VALUES ($1, $2, $3, $4, $5)", id, title, content, commentsDisabled, createdAt)
	if err != nil {
		return nil, err
	}
	return &models.Post{ID: id, Title: title, Content: content, CommentsDisabled: commentsDisabled, CreatedAt: createdAt}, nil
}

// CreateComment создает новый комментарий
func (repo *PostgresRepository) CreateComment(postId, parentId, content string) (*models.Comment, error) {
	id := uuid.New().String()
	createdAt := time.Now()
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var parentIdSQL interface{}
	if parentId == "" {
		parentIdSQL = nil
	} else {
		parentIdSQL = parentId
	}

	_, err = tx.Exec("INSERT INTO comments (id, post_id, parent_id, content, created_at) VALUES ($1, $2, $3, $4, $5)", id, postId, parentIdSQL, content, createdAt)
	if err != nil {
		return nil, err
	}

	if parentId != "" {
		_, err = tx.Exec("INSERT INTO pairs (parent_id, child_id) VALUES ($1, $2)", parentId, id)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &models.Comment{ID: id, PostID: postId, ParentID: &parentId, Content: content, CreatedAt: createdAt}, nil
}

// GetCommentsByPostID возвращает список комментариев для указанного поста с пагинацией
func (repo *PostgresRepository) GetCommentsByPostID(postId string, first int64, after *string) ([]*models.Comment, error) {
	query := `
        SELECT c.id, c.post_id, c.parent_id, c.content, c.created_at
        FROM comments c
        LEFT JOIN pairs p ON c.id = p.child_id
        WHERE c.post_id = $1 AND c.parent_id IS NULL
        ORDER BY c.created_at
        LIMIT $2
    `
	args := []interface{}{postId, first}

	if after != nil {
		query = `
        SELECT c.id, c.post_id, c.parent_id, c.content, c.created_at
        FROM comments c
        LEFT JOIN pairs p ON c.id = p.child_id
		WHERE c.post_id = $1 AND c.parent_id IS NULL AND c.id > $3::uuid
        ORDER BY c.created_at
        LIMIT $2
    `
		args = append(args, *after)
	}
	for i := range args {
		log.Printf("arg[%v] = %v", i, args[i])
	}
	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	// log.Printf("Returning comments...")
	// for _, comm := range comments {
	// 	log.Printf("comm: %v", comm.Content)
	// }
	return comments, nil
}

// GetCommentsByParentID возвращает список дочерних комментариев для указанного комментария с пагинацией
func (repo *PostgresRepository) GetCommentsByParentID(parentId string, first int64, after *string) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.parent_id, c.content, c.created_at
		FROM comments c
		WHERE c.parent_id = $1
		ORDER BY c.created_at
		LIMIT $2
	`
	args := []interface{}{parentId, first}

	if after != nil {
		query += "AND c.id > $3"
		args = append(args, *after)
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// DeletePost удаляет пост по его ID
func (repo *PostgresRepository) DeletePost(id string) error {
	// Сначала удаляем все комментарии к этому посту
	_, err := repo.db.Exec("DELETE FROM comments WHERE post_id = $1", id)
	if err != nil {
		return err
	}

	// Теперь удаляем сам пост
	_, err = repo.db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment удаляет комментарий по его ID
func (repo *PostgresRepository) DeleteComment(id string) error {
	// Сначала удаляем связи этого комментария с дочерними комментариями
	_, err := repo.db.Exec("DELETE FROM pairs WHERE parent_id = $1 OR child_id = $1", id)
	if err != nil {
		return err
	}

	// Теперь удаляем сам комментарий
	_, err = repo.db.Exec("DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
