package impl

import (
	"douyin-microservice/app/relation/models"
	"douyin-microservice/app/relation/rpc"
	"douyin-microservice/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

var chatConnMap = sync.Map{}

type MessageServiceImpl struct {
}

func (messageService MessageServiceImpl) SendMessage(token string, toUserId int64, content string) error {
	userClaim, err := utils.AnalyseToken(token)
	if err != nil {
		return err
	}
	user, err0 := rpc.GetUserById(userClaim.CommonEntity.Id)

	if err0 != nil {
		return err0
	}
	userId := user.Id

	err = models.SaveMessage(&models.Message{
		CommonEntity: utils.NewCommonEntity(),
		Content:      content,
	})
	if err != nil {
		return err
	}
	err = models.SaveMessageSendEvent(&models.MessageSendEvent{
		CommonEntity: utils.NewCommonEntity(),
		UserId:       userId,
		ToUserId:     toUserId,
		MsgContent:   content,
	})
	if err != nil {
		return err
	}
	err = models.SaveMessagePushEvent(&models.MessagePushEvent{
		CommonEntity: utils.NewCommonEntity(),
		FromUserId:   userId,
		MsgContent:   content,
	})
	if err != nil {
		return err
	}
	return nil
}

func (messageService MessageServiceImpl) GetHistoryOfChat(token string, toUserId int64, preMsgTime string) ([]models.MessageDVO, error) {
	userClaim, err := utils.AnalyseToken(token)
	if err != nil {
		return []models.MessageDVO{}, err
	}
	user, err0 := rpc.GetUserById(userClaim.CommonEntity.Id)

	if err0 != nil {
		return []models.MessageDVO{}, err0
	}
	userId := user.Id

	var preTime time.Time
	if preMsgTime != "0" {
		me, _ := strconv.ParseInt(preMsgTime, 10, 64)
		preTime = time.Unix(me, 0)
		if preTime.Year() > 9999 {
			preTime = time.Unix(me/1000, 0)
		}
	} else {
		preTime = time.Unix(0, 0)
	}

	//find from meesageSendEvent table
	var messageSendEvents []models.MessageSendEvent
	if preMsgTime == "0" {
		messageSendEvents, err = models.FindMessageSendEventByUserIdAndToUserId(userId, toUserId, preTime)
	}
	if err != nil {
		return nil, err
	}
	messageSendEventsOpposite, err := models.FindMessageSendEventByUserIdAndToUserId(toUserId, userId, preTime)
	if err != nil {
		return nil, err
	}
	messageSendEvents = append(messageSendEvents, messageSendEventsOpposite...)
	//sort.Sort(models.ByCreateTime(messageSendEvents))

	var messages []models.MessageDVO
	var wg sync.WaitGroup
	for _, messageSendEvent := range messageSendEvents {
		wg.Add(1)
		go func(messageSendEvent models.MessageSendEvent) {
			defer wg.Done()
			messages = append(messages, models.MessageDVO{
				Id:         messageSendEvent.Id,
				UserId:     messageSendEvent.UserId,
				ToUserId:   messageSendEvent.ToUserId,
				Content:    messageSendEvent.MsgContent,
				CreateTime: messageSendEvent.CreateDate.Unix(),
			})
		}(messageSendEvent)
	}
	wg.Wait()

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreateTime < messages[j].CreateTime
	})
	return messages, nil
}

func RunMessageServer() {
	listen, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Printf("Run message sever failed: %v\n", err)
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("Accept conn failed: %v\n", err)
			continue
		}

		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()

	var buf [256]byte
	for {
		n, err := conn.Read(buf[:])
		if n == 0 {
			if err == io.EOF {
				break
			}
			fmt.Printf("Read message failed: %v\n", err)
			continue
		}

		var event = models.MessageSendEvent{}
		_ = json.Unmarshal(buf[:n], &event)
		fmt.Printf("Receive Message：%+v\n", event)

		fromChatKey := fmt.Sprintf("%d_%d", event.UserId, event.ToUserId)
		if len(event.MsgContent) == 0 {
			chatConnMap.Store(fromChatKey, conn)
			continue
		}

		toChatKey := fmt.Sprintf("%d_%d", event.ToUserId, event.UserId)
		writeConn, exist := chatConnMap.Load(toChatKey)
		if !exist {
			fmt.Printf("User %d offline\n", event.ToUserId)
			continue
		}

		pushEvent := models.MessagePushEvent{
			FromUserId: event.UserId,
			MsgContent: event.MsgContent,
		}
		pushData, _ := json.Marshal(pushEvent)
		_, err = writeConn.(net.Conn).Write(pushData)
		if err != nil {
			fmt.Printf("Push message failed: %v\n", err)
		}
	}
}
