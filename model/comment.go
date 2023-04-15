package model

import (
	"time"
)

type Comment struct {
	ID string
	UserID string
	PostID string
	Content string
	CreatedAt time.Time
	
	Username string
}

