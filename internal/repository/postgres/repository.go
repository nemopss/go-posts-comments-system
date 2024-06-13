package postgres

import (
	"database/sql"
	"errors"
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
	log.Println("Querying posts...")
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

		// Получение комментариев для данного поста
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
	log.Println("Querying post with ID:", id)
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
	id := uuid.New().String() // Генерация нового уникального ID для поста
	log.Println("Creating post with ID:", id)
	createdAt := time.Now() // Текущее время как время создания поста
	_, err := repo.db.Exec("INSERT INTO posts (id, title, content, comments_disabled, created_at) VALUES ($1, $2, $3, $4, $5)", id, title, content, commentsDisabled, createdAt)
	if err != nil {
		return nil, err
	}
	log.Println("Created post", id)
	return &models.Post{ID: id, Title: title, Content: content, CommentsDisabled: commentsDisabled, CreatedAt: createdAt}, nil
}

// CreateComment создает новый комментарий
func (repo *PostgresRepository) CreateComment(postId, parentId, content string) (*models.Comment, error) {
	if len(content) > 2000 {
		return nil, errors.New("комментарий не может превышать 2000 символов")
	}
	var commentsDisabled bool
	err := repo.db.QueryRow("SELECT comments_disabled FROM posts WHERE id = $1", postId).Scan(&commentsDisabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Post not found!")
		}
		return nil, err
	}

	// Проверка, отключены ли комментарии для поста
	if commentsDisabled {
		return nil, errors.New("Comments are disabled on this post!")
	}
	id := uuid.New().String() // Генерация нового уникального ID для комментария
	log.Println("Creating comment with ID:", id)
	createdAt := time.Now() // Текущее время как время создания комментария
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
	log.Printf("Created comment on post %v...\n", postId)
	return &models.Comment{ID: id, PostID: postId, ParentID: &parentId, Content: content, CreatedAt: createdAt}, nil
}

// GetCommentsByPostID возвращает список комментариев для указанного поста с пагинацией
func (repo *PostgresRepository) GetCommentsByPostID(postId string, first int64, after *string) ([]*models.Comment, error) {
	log.Println("Getting comments on post with ID:", postId)
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
        WHERE c.post_id = $1 AND c.parent_id IS NULL AND c.id > $3
        ORDER BY c.created_at
        LIMIT $2
    `
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
		var parentId sql.NullString
		err := rows.Scan(&comment.ID, &comment.PostID, &parentId, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		if parentId.Valid {
			comment.ParentID = &parentId.String
		} else {
			comment.ParentID = nil
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// GetCommentsByParentID возвращает список дочерних комментариев для указанного комментария с пагинацией
func (repo *PostgresRepository) GetCommentsByParentID(parentId string, first int64, after *string) ([]*models.Comment, error) {
	log.Println("Getting comments from parent with ID:", parentId)
	query := `
		SELECT c.id, c.post_id, c.parent_id, c.content, c.created_at
		FROM comments c
		WHERE c.parent_id = $1
		ORDER BY c.created_at
		LIMIT $2
	`
	args := []interface{}{parentId, first}

	if after != nil {
		query += "AND c.id > $3::uuid"
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
	log.Println("Deleting post with ID:", id)
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Сначала удаляем все связи комментариев этого поста
	_, err = tx.Exec("DELETE FROM pairs WHERE child_id IN (SELECT id FROM comments WHERE post_id = $1)", id)
	if err != nil {
		return err
	}

	// Сначала удаляем все комментарии к этому посту
	_, err = tx.Exec("DELETE FROM comments WHERE post_id = $1", id)
	if err != nil {
		return err
	}

	// Теперь удаляем сам пост
	_, err = tx.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteComment удаляет комментарий по его ID
func (repo *PostgresRepository) DeleteComment(id string) error {

	log.Println("Deleting comment with ID:", id)
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Удаление всех вложенных комментариев
	err = repo.deleteChildComments(tx, id)
	if err != nil {
		return err
	}

	// Сначала удаляем связи этого комментария с дочерними комментариями
	_, err = tx.Exec("DELETE FROM pairs WHERE parent_id = $1 OR child_id = $1", id)
	if err != nil {
		return err
	}

	// Удаление самого комментарий
	_, err = tx.Exec("DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// deleteChildComments рекурсивно удаляет все дочерние комментарии
func (repo *PostgresRepository) deleteChildComments(tx *sql.Tx, parentId string) error {
	log.Println("Deleting child comments from parent with ID:", parentId)
	childComments := []string{}
	rows, err := tx.Query("SELECT id FROM comments WHERE parent_id = $1", parentId)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var childId string
		err := rows.Scan(&childId)
		if err != nil {
			return err
		}
		childComments = append(childComments, childId)
	}

	for _, childId := range childComments {
		err := repo.deleteChildComments(tx, childId)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM pairs WHERE parent_id = $1 OR child_id = $1", childId)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM comments WHERE id = $1", childId)
		if err != nil {
			return err
		}
	}

	return nil
}
