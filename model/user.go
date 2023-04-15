package model

import (
	"time"
)

type User struct {
	ID string
	Email string
	Username string
	CreatedAt time.Time
	
	Following bool
}
