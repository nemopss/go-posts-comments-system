package models

type Comment struct {
	ID       string
	PostID   string
	ParentID *string
	Content  string
	Children []*Comment
}
