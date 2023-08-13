package http

import (
	"douyin-microservice/app/gateway/rpc"
	"douyin-microservice/app/gateway/utils"
	"douyin-microservice/idl/pb"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	utils.Response
	VideoList []*pb.VideoDVO `json:"video_list,omitempty"`
	NextTime  int64          `json:"next_time,omitempty"`
}

func FeedHandler(c *gin.Context) {
	var feedReq pb.FeedRequest
	latestTimeStr := c.Query("latest_time")
	token := c.Query("token")
	var userId int64 = -1

	log.Printf("时间戳", latestTimeStr)
	var latestTime time.Time
	if latestTimeStr != "0" {
		me, _ := strconv.ParseInt(latestTimeStr, 10, 64)
		latestTime = time.Unix(me, 0)
		// 前端传入的可能是毫秒级
		if latestTime.Year() > 9999 {
			latestTime = time.Unix(me/1000, 0)
		}
	} else {
		latestTime = time.Now()
	}
	log.Printf("获取到的时间 %v", latestTime)

	if token != "" {
		userClaims, err0 := utils.AnalyseToken(token)
		if err0 != nil {
			log.Println("解析token失败")
		}
		userId = userClaims.CommonEntity.Id
	}

	feedReq.LatestTime = latestTimeStr
	feedReq.UserId = userId
	feedResp, err := rpc.Feed(c, &feedReq)
	if err != nil {
		c.JSON(http.StatusOK, utils.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  utils.Response{StatusCode: 0},
		VideoList: feedResp.VideoList,
		NextTime:  feedResp.NextTime,
	})
}
