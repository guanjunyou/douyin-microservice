package models

import "douyin-microservice/app/gateway/utils"

type CommentMQToVideo struct {
	utils.CommonEntity
	ActionType int    `json:"action_type"`
	UserId     User   `json:"user"`
	VideoId    int64  `json:"video_id"`
	Content    string `json:"content"`
	CommentID  int64  `json:"id"`
}
