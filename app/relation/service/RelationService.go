package service

import "douyin-microservice/app/relation/models"

type RelationService interface {
	// FollowUser 关注用户
	FollowUser(userId int64, toUserId int64, actionType int) error

	// GetFollows 查询关注列表
	GetFollows(userId int64) ([]models.User, error)

	// GetFollowers 查询粉丝列表
	GetFollowers(userId int64) ([]models.User, error)

	// GetFriends 查询好友列表
	GetFriends(usrId int64) ([]models.User, error)

	//查询是否关注
	CheckFollowForUser(userId int64, toUserId int64) bool
}
