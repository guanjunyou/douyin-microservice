package service

import (
	"douyin-microservice/app/video/models"
	"time"
)

type VideoService interface {
	// Publish
	// 将传入的视频流保存在文件服务器中，并存储在mysql表中
	// 5.23 加入title
	Publish(data []byte, userId int64, title string, filename string) error

	// PublishList
	// 通过userId来查询对应用户发布的视频，并返回对应的视频切片数组
	PublishList(userId int64) ([]models.VideoDVO, error)

	// GetVideoList
	GetVideoListByLastTime(latestTime time.Time, userId int64) ([]models.VideoDVO, time.Time, error)
}
