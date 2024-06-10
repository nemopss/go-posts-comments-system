package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/nemopss/go-posts-comments-system/internal/repository"
)

// Resolver отвечает за реализацию функций, которые будут вызываться.
// для разрешения запросов и мутаций GraphQL.
type Resolver struct {
	repo repository.Repository
}

// NewResolver создаёт новый экземпляр Resolver с заданным репозиторием.
func NewResolver(repo repository.Repository) *Resolver {
	return &Resolver{repo: repo}
}

// QueryPosts возвращает список всех постов.
// Этот метод вызывается при запросе поля `posts` в схеме GraphQL.
func (r *Resolver) QueryPosts(params graphql.ResolveParams) (interface{}, error) {
	return r.repo.GetPosts()
}

// QueryPost возвращает пост по его идентификатору.
// Этот метод вызывается при запросе поля `post` с идентификатором `id` в схеме GraphQL.
func (r *Resolver) QueryPost(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	return r.repo.GetPost(id)
}

// Create post создаёт новый пост.
// Этот метод вызывается при выполнении мутации `createPost` в схеме GraphQl c аргументами `title`, `content`, `commentsDisabled`.
func (r *Resolver) CreatePost(params graphql.ResolveParams) (interface{}, error) {
	title := params.Args["title"].(string)
	content := params.Args["content"].(string)
	commentsDisabled := params.Args["commentsDisabled"].(bool)
	return r.repo.CreatePost(title, content, commentsDisabled)
}

// CreateComment создаёт новый комментарий.
// Этот метод вызывается при выполнении мутации `createComment` в схеме GraphQL с аргументами `postId`, `parentId`, `content`
func (r *Resolver) CreateComment(params graphql.ResolveParams) (interface{}, error) {
	postId := params.Args["postId"].(string)
	parentId := params.Args["parentId"].(string)
	content := params.Args["content"].(string)
	return r.repo.CreateComment(postId, parentId, content)
}
