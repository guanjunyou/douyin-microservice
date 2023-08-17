package service

import "douyin-microservice/app/relation/models"

type MessageService interface {

	// SendMessage SendMsg 发送消息
	SendMessage(token string, toUserId int64, content string) error

	// GetHistoryOfChat 查看消息记录
	GetHistoryOfChat(token string, toUserId int64) ([]models.MessageDVO, error)
}
