package postgres

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nemopss/go-posts-comments-system/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

// Инициализация базы данных для тестирования
func init() {
	var err error
	testDB, err = sql.Open("postgres", "postgres://gosuper:Ukflbkby2004@localhost:5432/go-posts-comments-db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

// cleanDatabase очищает все таблицы перед запуском тестов
func cleanDatabase(db *sql.DB) {
	tables := []string{"pairs", "comments", "posts"}
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE " + table + " CASCADE")
		if err != nil {
			log.Fatalf("failed to clean table %s: %v", table, err)
		}
	}
}

func TestPostgresRepository(t *testing.T) {
	repo := postgres.NewPostgresRepository(testDB)

	// Тест GetCommentsByPostID
	t.Run("TestGetCommentsByPostID_Postgres", func(t *testing.T) {
		cleanDatabase(testDB)
		//Создание поста
		post, err := repo.CreatePost("Title", "Content", false)
		assert.NoError(t, err)
		// Создание комментариев
		_, err = repo.CreateComment(post.ID, "", "Comment 1")
		assert.NoError(t, err)
		_, err = repo.CreateComment(post.ID, "", "Comment 2")
		assert.NoError(t, err)
		// Проверка получения комментариев
		comments, err := repo.GetCommentsByPostID(post.ID, 10, nil)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(comments))
	})
	// Тест GetCommentsByParentID
	t.Run("TestGetCommentsByParentID_Postgres", func(t *testing.T) {
		cleanDatabase(testDB)
		//Создание поста
		post, err := repo.CreatePost("Title", "Content", false)
		assert.NoError(t, err)

		// Создание комментариев
		comment1, err := repo.CreateComment(post.ID, "", "Comment 1")
		assert.NoError(t, err)
		_, err = repo.CreateComment(post.ID, comment1.ID, "Comment 1.1")
		assert.NoError(t, err)
		_, err = repo.CreateComment(post.ID, comment1.ID, "Comment 1.2")
		assert.NoError(t, err)

		// Проверка получения комментариев
		comments, err := repo.GetCommentsByParentID(comment1.ID, 10, nil)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(comments))
	})

	// Проверка удаления поста
	t.Run("TestDeletePost_Postgres", func(t *testing.T) {
		cleanDatabase(testDB)
		//Создание поста
		post, err := repo.CreatePost("Title", "Content", false)
		assert.NoError(t, err)

		// Создание комментариев
		_, err = repo.CreateComment(post.ID, "", "Comment 1")
		assert.NoError(t, err)
		_, err = repo.CreateComment(post.ID, "", "Comment 2")
		assert.NoError(t, err)

		// Удаление поста
		err = repo.DeletePost(post.ID)
		assert.NoError(t, err)

		_, err = repo.GetPost(post.ID)
		assert.Error(t, err)
		// Проверка на отсутствие поста
		comments, err := repo.GetCommentsByPostID(post.ID, 10, nil)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(comments))
	})

	t.Run("TestDeleteComment_Postgres", func(t *testing.T) {
		cleanDatabase(testDB)

		//Создание поста
		post, err := repo.CreatePost("Title", "Content", false)
		assert.NoError(t, err)

		// Создание комментариев
		comment1, err := repo.CreateComment(post.ID, "", "Comment 1")
		assert.NoError(t, err)
		comment2, err := repo.CreateComment(post.ID, comment1.ID, "Comment 1.1")
		assert.NoError(t, err)
		_ = comment2
		_, err = repo.CreateComment(post.ID, comment1.ID, "Comment 1.2")
		assert.NoError(t, err)

		err = repo.DeleteComment(comment1.ID)
		assert.NoError(t, err)

		// Проверка на удаление комментария
		comments, err := repo.GetCommentsByPostID(post.ID, 10, nil)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(comments))
	})
}
