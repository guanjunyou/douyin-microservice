package test

import (
	"context"
	"github.com/zhenghaoz/gorse/client"
	"log"
	"testing"
	"time"
)

func TestGorseTest(t *testing.T) {
	SendLikeToGorse(7090306410939941888, 7092363100510225408)
}

func TestInsertGorseTest(t *testing.T) {
	Client := client.NewGorseClient("http://127.0.0.1:8088", "")
	//title := "合成高温超导体"
	videoId := "7099396439897277448"
	timestamp := time.Unix(time.Now().Unix(), 0).UTC().Format(time.RFC3339)
	//x := gojieba.NewJieba()
	//defer x.Free()
	//labels := x.Cut(title, true)
	item := client.Item{
		ItemId:     videoId,
		IsHidden:   false,
		Labels:     []string{"研究", "高温", "超导体"},
		Categories: []string{"video"},
		Timestamp:  timestamp,
		Comment:    "",
	}
	Client.InsertItem(context.Background(), item)
}

func SendLikeToGorse(userId int64, videoId int64) {
	Client := client.NewGorseClient("http://127.0.0.1:8088", "")
	//video, err := models.GetVideoById(videoId)
	//if err != nil {
	//	log.Println(err.Error())
	//}
	//timestamp := time.Unix(time.Now().Unix(), 0).UTC().Format(time.RFC3339)
	//Client.InsertFeedback(context.Background(), []client.Feedback{{
	//	FeedbackType: "like",
	//	UserId:       strconv.FormatInt(userId, 10),
	//	Timestamp:    timestamp,
	//	ItemId:       video.Title,
	//}})
	recommend, _ := Client.GetRecommend(context.Background(), "7090306410939941888", "", 10)
	for i := range recommend {
		log.Println(recommend[i])
	}
}
