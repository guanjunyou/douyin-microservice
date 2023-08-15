package models

import (
	"douyin-microservice/idl/pb"
	"douyin-microservice/pkg/utils"
)

type VideoDVO struct {
	utils.CommonEntity
	Author        pb.User `json:"author"`
	PlayUrl       string  `json:"play_url"`
	CoverUrl      string  `json:"cover_url"`
	FavoriteCount int64   `json:"favorite_count"`
	CommentCount  int64   `json:"comment_count"`
	IsFavorite    bool    `json:"is_favorite"`
	Title         string  `json:"title,omitempty"`
}
