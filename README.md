# Система постов и комментариев с использованием GraphQL

## Описание

Данный проект реализует систему для добавления и чтения постов и комментариев на языке Go с использованием GraphQL.

## Характеристики

### Система постов

- Можно просмотреть список постов.
- Можно просмотреть пост и комментарии под ним.
- Пользователь, создавая пост, может запретить оставление комментариев к своему посту.

### Система комментариев

- Комментарии организованы иерархически, позволяя вложенность без ограничений.
- Длина текста комментария ограничена до 2000 символов.

## Стек технологий

- Система написана на языке Go.
- Используется Docker для распространения сервиса в виде Docker-образа.
- Хранение данных может быть как в памяти (in-memory), так и в PostgreSQL. Выбор хранилища определяется параметром при запуске сервиса.

## Как запустить

1. Склонируйте репозиторий:

   ```bash
   git clone https://github.com/nemopss/go-posts-comments-system.git
   ```

2. Перейдите в директорию проекта:

   ```bash
   cd go-posts-comments-system
   ```

3. Запустите сервис с использованием Docker:

   ```bash
   docker-compose -f docker-compose-inmemory.yml up --build
   ```

   Если вы хотите использовать PostgreSQL в качестве хранилища данных, запустите сервис следующим образом:

   ```bash
   docker-compose -f docker-compose-postgres.yml up --build
   ```

4. Откройте GraphiQL в браузере по адресу `http://localhost:8080/graphql` и начните работу с API.

## Работа с API
Создать пост:
```graphql
mutation {
  createPost(title: "название_поста" content: "контент_поста", commentsDisabled:false) {
    id
    content
    commentsDisabled
  }
}
```
- ```title``` - название поста
- ```content``` - нонтент поста
- ```commentsDisabled``` -  флаг, определяющий возможность оставлять комментарии к посту

Удалить пост:
```graphql
mutation {
  deletePost(id: "айди_поста") 
}
```
- ```id``` - айди поста, который вы хотите удалить

Создать комментарий:
```graphql
mutation {
  createComment(parentId: "", content: "контент_коментария", postId: "айди_поста") {
    id
    parentId
    content
  }
}
```
- ```parentId``` - айди родительского комментария, оставьте пустым, если хотите оставить комментарий напрямую к посту
- ```content``` - контент комментария
- ```postId``` - айди поста, к которому вы хотите оставить комментарий

Удалить комментарий:
```graphql
mutation {
  deleteComment(id: "АЙДИ_КОММЕНТАРИЯ") 
}
```
- ```id``` - айди комментария, который вы хотите удалить

Вывести все посты: 
```graphql
fragment CommentFields on Comment {
  id
  content
  parentId
}

{
  posts {
    id
    title
    content
   comments(first:5, after:"") {
      ...CommentFields
      children (first:5, after:""){
        ...CommentFields
        children (first:5, after:""){
          ...CommentFields
        }
      }
    }
  }
}
```
Выводит все посты с глубиной комментариев 3.
GraphQL не поддерживает рекурсивный запрос дочерних комментариев, поэтому, если нужно вывести больне комментариев, то нужно дополнить структуру 
```graphql
    comments(first:5, after:"") {
      ...CommentFields
      children (first:5, after:""){
        ...CommentFields
        children (first:5, after:""){
          ...CommentFields
        }
      }
    }
```
до нужного количества вложенности

- ```first``` - параметр, отвечающий за количество выводимых комментариев на каждом уровне вложенности

Вывести отдельный пост
```graphql
fragment CommentFields on Comment {
  id
  content
  parentId
}

{
  post(id: "айди_поста") {
    id
    title
    content
   comments(first:5, after:"") {
      ...CommentFields
      children (first:5, after:""){
        ...CommentFields
        children (first:5, after:""){
          ...CommentFields
        }
      }
    }
  }
}
```
- ```id``` - айди поста, который вы хотите вывести

## Структура проекта 
```
graphql-comments-system/
├── cmd/
│   └── main.go
├── internal/
│   ├── gql/
│   │   ├── resolvers.go
│   │   ├── schema.graphql
│   │   └── schema.go
│   ├── models/
│   │   ├── comment.go
│   │   └── post.go
│   ├── repository/
│   │   ├── inmemory/
│   │   │   └── repository.go
│   │   ├── postgres/
│   │   │   └── repository.go
│   │   └── repository.go
│   ├── server/
│   │   └── server.go
│   └── test/
│       ├── inmemory/
│       │   └── inmemory_test.go
│       └── postgres/
│           └── postgres_test.go
├── Dockerfile
├── docker-compose-inmemory.yml
├── docker-compose-postgres.yml
├── init.sql
├── go.mod
└── go.sum
```
