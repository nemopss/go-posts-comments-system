package gql

import (
	"log"

	"github.com/graphql-go/graphql"
	"github.com/nemopss/go-posts-comments-system/internal/models"
	"github.com/nemopss/go-posts-comments-system/internal/repository"
)

// Resolver отвечает за реализацию функций, которые будут вызываться.
// для разрешения запросов и мутаций GraphQL.
type Resolver struct {
	repo repository.Repository
}

// NewResolver создаёт новый экземпляр Resolver с заданным репозиторием.
func NewResolver(repo repository.Repository) *Resolver {
	log.Println("Creating resolver...")
	return &Resolver{repo: repo}
}

// QueryPosts возвращает список всех постов.
// Этот метод вызывается при запросе поля `posts` в схеме GraphQL.
func (r *Resolver) QueryPosts(params graphql.ResolveParams) (interface{}, error) {
	log.Println("Quering posts...")
	return r.repo.GetPosts()
}

// QueryPost возвращает пост по его идентификатору.
// Этот метод вызывается при запросе поля `post` с идентификатором `id` в схеме GraphQL.
func (r *Resolver) QueryPost(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)

	log.Println("Quering post...")
	return r.repo.GetPost(id)
}

// Create post создаёт новый пост.
// Этот метод вызывается при выполнении мутации `createPost` в схеме GraphQl c аргументами `title`, `content`, `commentsDisabled`.
func (r *Resolver) CreatePost(params graphql.ResolveParams) (interface{}, error) {
	title := params.Args["title"].(string)
	content := params.Args["content"].(string)
	commentsDisabled := params.Args["commentsDisabled"].(bool)
	log.Println("Creating post...")
	return r.repo.CreatePost(title, content, commentsDisabled)
}

// CreateComment создаёт новый комментарий.
// Этот метод вызывается при выполнении мутации `createComment` в схеме GraphQL с аргументами `postId`, `parentId`, `content`
func (r *Resolver) CreateComment(params graphql.ResolveParams) (interface{}, error) {
	postId := params.Args["postId"].(string)
	parentId := params.Args["parentId"].(string)
	content := params.Args["content"].(string)
	log.Println("Creating comment...")
	return r.repo.CreateComment(postId, parentId, content)
}

// ResolvePostComments возвращает список комментариев для заданного поста с пагинацией
func (r *Resolver) ResolvePostComments(p graphql.ResolveParams) (interface{}, error) {
	post := p.Source.(*models.Post)
	first, _ := p.Args["first"].(int64)
	after, _ := p.Args["after"].(string)
	return r.repo.GetCommentsByPostID(post.ID, first, &after)
}

// ResolveCommentChildren возвращает список дочерних комментариев для заданного комментария с пагинацией
func (r *Resolver) ResolveCommentChildren(p graphql.ResolveParams) (interface{}, error) {
	comment := p.Source.(*models.Comment)
	first, _ := p.Args["first"].(int64)
	after, _ := p.Args["after"].(string)
	return r.repo.GetCommentsByParentID(comment.ID, first, &after)
}

// DeletePost удаляет пост по его ID
func (r *Resolver) DeletePost(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	err := r.repo.DeletePost(id)
	if err != nil {
		return nil, err
	}
	return true, nil
}

// DeleteComment удаляет комментарий по его ID
func (r *Resolver) DeleteComment(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	err := r.repo.DeleteComment(id)
	if err != nil {
		return nil, err
	}
	return true, nil
}
