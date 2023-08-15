package http

import (
	"douyin-microservice/app/gateway/rpc"
	"douyin-microservice/idl/pb"
	utils2 "douyin-microservice/pkg/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
)

type FeedResponse struct {
	utils2.Response
	VideoList []*pb.VideoDVO `json:"video_list,omitempty"`
	NextTime  int64          `json:"next_time,omitempty"`
}

type VideoListResponse struct {
	utils2.Response
	VideoList []*pb.VideoDVO `json:"video_list"`
}

func FeedHandler(c *gin.Context) {
	var feedReq pb.FeedRequest
	latestTimeStr := c.Query("latest_time")
	token := c.Query("token")
	var userId int64 = -1

	log.Printf("时间戳", latestTimeStr)

	if token != "" {
		userClaims, err0 := utils2.AnalyseToken(token)
		if err0 != nil {
			log.Println("解析token失败")
		}
		userId = userClaims.CommonEntity.Id
	}

	feedReq.LatestTime = latestTimeStr
	feedReq.UserId = userId
	feedResp, err := rpc.Feed(c, &feedReq)
	if err != nil {
		c.JSON(http.StatusOK, utils2.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  utils2.Response{StatusCode: 0},
		VideoList: feedResp.VideoList,
		NextTime:  feedResp.NextTime,
	})
}

func PublishHandler(c *gin.Context) {
	var publishReq pb.PublishRequest
	//1.获取token并解析出user_id、data、title
	token := c.PostForm("token")
	userClaims, _ := utils2.AnalyseToken(token)
	userId := userClaims.CommonEntity.Id
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, utils2.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	title := c.PostForm("title")
	publishReq.Title = title
	publishReq.UserId = userId
	dataBytes := FileHeaderToBytes(data)
	publishReq.Data = dataBytes
	publishReq.FileName = filepath.Base(data.Filename)
	//2. 调用service层处理业务逻辑
	err = rpc.Pubilsh(c, &publishReq)
	if err != nil {
		c.JSON(http.StatusOK, utils2.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, utils2.Response{
		StatusCode: 0,
		StatusMsg:  "投稿成功！",
	})
	return
}

func PublishListHandler(c *gin.Context) {
	var publishListReq pb.PublishListRequest
	//获取用户id
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: utils2.Response{
				StatusCode: 1,
				StatusMsg:  "类型转换错误",
			},
			VideoList: nil,
		})
	}
	publishListReq.UserId = userId
	publishListResp, err := rpc.PublishList(c, &publishListReq)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: utils2.Response{
				StatusCode: 1,
				StatusMsg:  "数据库异常",
			},
			VideoList: nil,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: utils2.Response{
			StatusCode: 0,
			StatusMsg:  "查询成功",
		},
		VideoList: publishListResp.VideoList,
	})
}

func FileHeaderToBytes(data *multipart.FileHeader) []byte {
	dataFile, _ := data.Open()
	dataBytes, _ := ioutil.ReadAll(dataFile)
	return dataBytes
}
