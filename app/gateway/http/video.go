package http

import (
	"douyin-microservice/app/gateway/rpc"
	"douyin-microservice/app/gateway/utils"
	"douyin-microservice/idl/pb"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
