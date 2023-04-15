package model

import (
	"time"
)

type Post struct {
	ID string
	UserID string
	Content string
	CreatedAt time.Time
	
	Username string
}
