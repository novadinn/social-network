package service

import (
	"github.com/novadinn/social-network/model"
)

func (svc *Service) GetUserByUsername(username string) (model.User, error) {
	return svc.Queries.GetUserByUsername(username)
}

func (svc *Service) GetFollowingUser(username, followingID string) (model.User, error) {
	user, err := svc.Queries.GetUserByUsername(username)
	if err != nil {
		return model.User{}, err
	}

	exists, err := svc.Queries.UserFollowExists(followingID, user.ID)
	if err != nil {
		return model.User{}, err
	}

	user.Following = exists
	
	return user, nil
}
