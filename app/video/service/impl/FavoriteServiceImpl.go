package impl

import (
	"context"
	"douyin-microservice/app/gateway/utils"
	"douyin-microservice/app/video/models"
	"douyin-microservice/config"
	"strconv"
)

type FavoriteServiceImpl struct {
}

func (favoriteService FavoriteServiceImpl) FindIsFavouriteByUserIdAndVideoId(userId int64, videoId int64) bool {
	//tx := utils.GetMysqlDB()
	likeKey := config.LikeKey + strconv.FormatInt(userId, 10)
	videoIdStr := strconv.FormatInt(videoId, 10)
	exists, _ := utils.GetRedisDB().Exists(context.Background(), likeKey).Result()
	if exists != 0 {
		videoExists, _ := utils.GetRedisDB().SIsMember(context.Background(), likeKey, videoIdStr).Result()
		return videoExists
	}
	like := models.Like{
		UserId:  userId,
		VideoId: videoId,
	}

	isLike, _ := like.FindByUserIdAndVedioId()

	if isLike.Id != 0 {
		return true
	} else {
		return false
	}
}
