package models

type FollowMQToUser struct {
	UserId       int64 `json:"user_id"`
	FollowUserId int64 `json:"follow_user_id"`
	ActionType   int   `json:"action_type"`
}
