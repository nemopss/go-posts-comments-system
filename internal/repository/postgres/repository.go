package postgres

import (
	"database/sql"
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
	rows, err := repo.db.Query("SELECT id, title, content, comments_disabled FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CommentsDisabled, &post.Comments)
		if err != nil {
			return nil, err
		}
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
	_, err := repo.db.Exec("INSERT INTO posts (id, title, content, comments_disabled) VALUES ($1, $2, $3, $4)", id, title, content, commentsDisabled)
	if err != nil {
		return nil, err
	}
	return &models.Post{ID: id, Title: title, Content: content, CommentsDisabled: commentsDisabled}, nil
}

// CreateComment создает новый комментарий
func (repo *PostgresRepository) CreateComment(postId, parentId, content string) (*models.Comment, error) {
	id := uuid.New().String()
	_, err := repo.db.Exec("INSERT INTO comments (id, post_id, parent_id, content) VALUES ($1, $2, $3, $4)", id, postId, parentId, content)
	if err != nil {
		return nil, err
	}
	return &models.Comment{ID: id, PostID: postId, ParentID: &parentId, Content: content}, nil
}

// GetCommentsByPostID возвращает список комментариев для указанного поста
func (repo *PostgresRepository) GetCommentsByPostID(postId string) ([]*models.Comment, error) {
	rows, err := repo.db.Query("SELECT id, post_id, parent_id, content FROM comments WHERE post_id=$1 AND parent_id IS NULL", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// GetCommentsByParentID возвращает список дочерних комментариев для указанного комментария
func (repo *PostgresRepository) GetCommentsByParentID(parentId string) ([]*models.Comment, error) {
	rows, err := repo.db.Query("SELECT id, post_id, parent_id, content FROM comments WHERE parent_id=$1", parentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
