package service

import (
	"errors"
	"time"
	"unicode/utf8"
	
	"github.com/novadinn/social-network/model"
)

var (
	ErrInvalidCommentContent = errors.New("invalid comment content")
)

func (svc *Service) CreateComment(userID, postID, content string, loggedIn bool) (model.Comment, error) {
	if content == "" || utf8.RuneCountInString(content) > 1023 {
		return model.Comment{}, ErrInvalidCommentContent
	}

	if !loggedIn {
		return model.Comment{}, ErrUnauthenticated
	}

	id := genID()
	createdAt := time.Now()

	svc.Queries.CreateComment(id, userID, postID, content, createdAt)
	
	return model.Comment{
		ID: id,
		UserID: userID,
		PostID: postID,
		Content: content,
		CreatedAt: createdAt,
	}, nil
}

func (svc *Service) GetComments(postID string) ([]model.Comment, error) {
	return svc.Queries.GetComments(postID)
}
