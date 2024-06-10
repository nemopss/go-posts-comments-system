package models

type Post struct {
	ID               string
	Title            string
	Content          string
	Comments         []*Comment
	CommentsDisabled bool
}
