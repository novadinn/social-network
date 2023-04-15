package service

import (
	"errors"
	"regexp"
	"time"

	"github.com/novadinn/social-network/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUsernameTaken = errors.New("username taken")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidUsername = errors.New("invalid username")
)

var (
	reEmail = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	reUsername = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{0,17}$`)
)

func (svc *Service) Login(email, username string) (model.User, error) {
	valid := reEmail.MatchString(email)
	if !valid {
		return model.User{}, ErrInvalidEmail
	}

	if username != "" {
		valid = reUsername.MatchString(username)
		if !valid {
			return model.User{}, ErrInvalidUsername
		}
	}
	
	exists, err := svc.Queries.UserExistsByEmail(email)
	if err != nil {
		return model.User{}, err
	}
	
	if exists {
		return svc.Queries.GetUserByEmail(email)
	}

	if username == "" {
		return model.User{}, ErrUserNotFound
	}

	exists, err = svc.Queries.UserExistsByUsername(username)
	if err != nil {
		return model.User{}, err
	}

	if exists {
		return model.User{}, ErrUsernameTaken
	}

	id := genID()
	createdAt := time.Now()

	if err = svc.Queries.CreateUser(id, email, username, createdAt); err != nil {
		return model.User{}, err
	}
	
	return model.User{ID: id, Email: email, Username: username, CreatedAt: createdAt}, nil
}
