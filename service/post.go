package service

import (
	"errors"
	"time"
	"unicode/utf8"
	
	"github.com/novadinn/social-network/model"
)

var (
	ErrInvalidPostContent = errors.New("invalid post content")
	ErrUnauthenticated = errors.New("unauthenticated")
)

func (svc *Service) CreatePost(content, userID string, loggedIn bool) error {
	if content == "" || utf8.RuneCountInString(content) > 1023 {
		return ErrInvalidPostContent
	}

	if !loggedIn {
		// TODO: this is not working
		return ErrUnauthenticated
	}

	id := genID()
	createdAt := time.Now()
	
	if err := svc.Queries.CreatePost(id, userID, content, createdAt); err != nil {
		return err
	}
	
	
	return nil
}

func (svc *Service) GetPosts() ([]model.Post, error) {
	return svc.Queries.GetPosts()
}

func (svc *Service) GetFollowingPosts(ids []interface{}) ([]model.Post, error) {
	return svc.Queries.GetFollowingPosts(ids)
}

func (svc *Service) GetPostsByUsername(username string) ([]model.Post, error) {
	return svc.Queries.GetPostsByUsername(username)
}

func (svc *Service) GetPostByID(id string) (model.Post, error) {
	return svc.Queries.GetPostByID(id)
}
