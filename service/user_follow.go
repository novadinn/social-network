package service

import (
	"errors"
	"time"
	
)

var (
	ErrCannotFollowSelf = errors.New("cannon follow self")
)

func (svc *Service) FollowUser(userID string, followUserID string) error {	
	if userID == followUserID {
		return ErrCannotFollowSelf
	}

	exists, err := svc.Queries.UserFollowExists(userID, followUserID)
	if err != nil {
		return err
	}
	
	if exists {
		return nil
	}

	exists, err = svc.Queries.UserExistsByID(followUserID)
	if err != nil {
		return err
	}

	if !exists {
		return ErrUserNotFound
	}

	if err = svc.Queries.CreateUserFollow(userID, followUserID, time.Now()); err != nil {
		return err
	}
	
	return nil
}

func (svc *Service) UnfollowUser(userID string, followUserID string) error {	
	if userID == followUserID {
		return ErrCannotFollowSelf
	}

	exists, err := svc.Queries.UserExistsByID(followUserID)
	if err != nil {
		return err
	}

	if !exists {
		return ErrUserNotFound
	}
	
	exists, err = svc.Queries.UserFollowExists(userID, followUserID)
	if err != nil {
		return err
	}
	
	if !exists {
		return nil
	}

	err = svc.Queries.DeleteUserFollow(userID, followUserID)
	if err != nil {
		return err
	}
	
	return nil
}
