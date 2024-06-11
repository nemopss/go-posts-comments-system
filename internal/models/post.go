package models

import "time"

type Post struct {
	ID               string
	Title            string
	Content          string
	Comments         []*Comment
	CommentsDisabled bool
	CreatedAt        time.Time
}
