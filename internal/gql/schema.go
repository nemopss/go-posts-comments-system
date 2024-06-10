package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/nemopss/go-posts-comments-system/internal/repository"
)

// NewSchema создаёт новую GraphQl схему, используя переданный репозиторий
func NewSchema(repo repository.Repository) (graphql.Schema, error) {
	// Создаём новый resolver
	resolver := NewResolver(repo)

	// Определяем тип commentType для последующей инициализации
	var commentType *graphql.Object

	// Определяем тип post для GraphQl схемы
	postType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"title": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"content": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"comments": &graphql.Field{
				Type:    graphql.NewList(commentType),
				Resolve: resolver.ResolvePostComments,
			},
			"commentsDisabled": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
	})

	// Определяем тип comment для GraphQl схемы
	commentType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"postId": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"parentId": &graphql.Field{
				Type: graphql.ID,
			},
			"content": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"children": &graphql.Field{
				Type:    graphql.NewList(commentType),
				Resolve: resolver.ResolveCommentChildren,
			},
		},
	})

	// Определяем корневой тип Query для GraphQl схемы
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"posts": &graphql.Field{
				Type: graphql.NewList(postType),
			},
			"post": &graphql.Field{
				Type: postType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: resolver.QueryPost,
			},
		},
	})

	// Определяем тип Mutation для GraphQl схемы
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createPost": &graphql.Field{
				Type: postType,
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"commentsDisabled": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Boolean),
					},
				},
			},
			"createComment": &graphql.Field{
				Type: commentType,
				Args: graphql.FieldConfigArgument{
					"postId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"parentId": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolver.CreateComment,
			},
		},
	})

	// Создаём конфиг схемы
	schemaConfig := graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}

	return graphql.NewSchema(schemaConfig)
}
