package http

import (
	"douyin-microservice/app/gateway/rpc"
	"douyin-microservice/app/video/models"
	"douyin-microservice/config"
	"douyin-microservice/idl/pb"
	"douyin-microservice/pkg/utils"
	"douyin-microservice/pkg/utils/resultutil"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type FeedResponse struct {
	utils.Response
	VideoList []*pb.VideoDVO `json:"video_list,omitempty"`
	NextTime  int64          `json:"next_time,omitempty"`
}
type CommentListResponse struct {
	utils.Response
	CommentList []models.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	utils.Response
	Comment models.Comment `json:"comment,omitempty"`
}

// 接收点赞的结构体
type FavoriteActionReq struct {
	Token      string `form:"token"`
	VideoId    string `form:"video_id"`    // 视频id
	ActionType string `form:"action_type"` // 1-点赞，2-取消点赞
}
type VideoListResponse struct {
	utils.Response
	VideoList []*pb.VideoDVO `json:"video_list"`
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

func PublishHandler(c *gin.Context) {
	var publishReq pb.PublishRequest
	//1.获取token并解析出user_id、data、title
	token := c.PostForm("token")
	userClaims, _ := utils.AnalyseToken(token)
	userId := userClaims.CommonEntity.Id
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, utils.Response{
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
		c.JSON(http.StatusOK, utils.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
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
			Response: utils.Response{
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
			Response: utils.Response{
				StatusCode: 1,
				StatusMsg:  "数据库异常",
			},
			VideoList: nil,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: utils.Response{
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

func FavoriteActionHandler(c *gin.Context) {
	var faReq FavoriteActionReq
	if err := c.ShouldBind(&faReq); err != nil {
		log.Printf("点赞操作，绑定参数发生异常：%v \n", err)
		resultutil.GenFail(c, "参数错误")
		return
	}
	fmt.Printf("参数 %+v \n", faReq)

	videoId, err := strconv.ParseInt(faReq.VideoId, 10, 64)

	if err != nil {
		log.Printf("点赞操作，videoId字符串转换发生异常 = %v", err)
		resultutil.GenFail(c, "参数错误")
		return
	}

	// 从Token中获取Uid
	var userClaim *utils.UserClaims
	userClaim, err = utils.AnalyseToken(faReq.Token)

	if err != nil {
		log.Printf("解析token发生异常 = %v", err)
		return
	}
	userId := userClaim.CommonEntity.Id

	var actionType int
	actionType, err = strconv.Atoi(faReq.ActionType)

	if err != nil {
		log.Printf("点赞操作，actionType字符串转换发生异常 = %v", err)
		resultutil.GenFail(c, "参数错误")
		return
	}
	var req pb.LikeVideoRequest
	req.UserId = userId
	req.VideoId = strconv.FormatInt(videoId, 10)
	req.ActionType = int32(actionType)
	if err = rpc.LikeVideo(c, &req); err != nil {
		log.Printf("点赞发生异常 = %v", err)
		if err.Error() == "-1" {
			resultutil.GenFail(c, "该视频已点赞")
			return
		}

		if err.Error() == "-2" {
			resultutil.GenFail(c, "未找到要取消的点赞记录")
			return
		}

		resultutil.GenFail(c, "点赞发生错误")
		return
	}

	resultutil.GenSuccessWithMsg(c, "success")
}
func FavoriteListHandler(c *gin.Context) {
	userIdStr := c.Query("user_id")

	//userId, err := strconv.ParseInt(userIdStr, 10, 64)
	var req pb.QueryVideosOfLikeRequest
	req.UserId = userIdStr
	resp, err := rpc.QueryVideosOfLike(c, &req)
	if err != nil {
		log.Printf("获取喜欢列表，获取发生异常 = %v", err)
		resultutil.GenFail(c, "获取失败")
		return
	}
	var videoList []models.LikeVedioListDVO
	for i := range resp.VideoList {
		videoList = append(videoList, BuildLikeVideo(*resp.VideoList[i]))
	}
	c.JSON(http.StatusOK, models.VideoListResponse2{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videoList,
	})
}
func CommentActionHandler(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	userClaim, err := utils.AnalyseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, utils.Response{StatusCode: 1, StatusMsg: "Token is invalid"})
		return
	}
	var userReq pb.UserRequest
	userReq.Username = userClaim.Name
	resp, err := rpc.UserService.GetUserByName(c, &userReq)
	user := DbUser2User(*resp.User)
	if err != nil {
		c.JSON(http.StatusOK, utils.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	if actionType == "1" {
		text := c.Query("comment_text")
		utils.InitFilter()

		textAfterFilter := utils.Filter.Replace(text, '*')

		comment := models.Comment{
			CommonEntity: utils.NewCommonEntity(),
			//Id:         1,
			User:    user,
			Content: textAfterFilter,
		}

		video_id := c.Query("video_id")
		var postReq pb.PostCommentsRequest
		postReq.VideoId = video_id
		postReq.Comment = BuildCommentDb(comment)
		_, err := rpc.VideoService.PostComments(c, &postReq)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, utils.Response{StatusCode: 1, StatusMsg: "Comment failed"})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{Response: utils.Response{StatusCode: 0, StatusMsg: ""},
			Comment: comment})
		return
	} else if actionType == "2" {
		//comment_id := parseCommetId(c)
		commentId := c.Query("comment_id")
		var deleteReq pb.DeleteCommentsRequest
		deleteReq.CommentId, _ = strconv.ParseInt(commentId, 10, 64)
		_, err := rpc.VideoService.DeleteComments(c, &deleteReq)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, utils.Response{StatusCode: 1, StatusMsg: "Delete comment failed"})
			return
		}
		c.JSON(http.StatusOK, utils.Response{StatusCode: 0})
		return
	}
	c.JSON(http.StatusOK, utils.Response{StatusCode: 0})
}
func CommentListHandler(c *gin.Context) {
	videoIdStr := c.Query("video_id")
	var req pb.CommentListRequest
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)
	req.VideoId = videoId
	resp, _ := rpc.VideoService.CommentList(c, &req)
	list := resp.Comments
	var commentList []models.Comment
	for i := range list {
		commentList = append(commentList, BuildComment(list[i]))
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    utils.Response{StatusCode: 0, StatusMsg: "Comment list"},
		CommentList: commentList,
	})
}

func BuildLikeVideo(videoPb pb.VideoDVO) models.LikeVedioListDVO {
	user := DbUser2User(*videoPb.Author)
	likeVideo := models.LikeVedioListDVO{
		Video: models.Video{
			CommonEntity:  utils.CommonEntity{Id: videoPb.Id},
			AuthorId:      videoPb.Author.Id,
			PlayUrl:       videoPb.PlayUrl,
			CoverUrl:      videoPb.CoverUrl,
			FavoriteCount: videoPb.FavoriteCount,
			CommentCount:  videoPb.CommentCount,
			IsFavorite:    videoPb.IsFavorite,
			Title:         videoPb.Title,
		},
		Author: &user,
	}
	return likeVideo
}

func BuildComment(commentDb *pb.Comment) models.Comment {
	createDate, _ := time.Parse(config.DateLayout, commentDb.CreateDate)
	comment := models.Comment{
		CommonEntity: utils.CommonEntity{Id: commentDb.Id, CreateDate: createDate, IsDeleted: 0},
		User:         DbUser2User(*commentDb.User),
		Content:      commentDb.Content,
	}
	return comment
}

func BuildCommentDb(comment models.Comment) *pb.Comment {
	commentPb := pb.Comment{
		Id:         comment.Id,
		User:       BuildUserPb(&comment.User),
		Content:    comment.Content,
		CreateDate: comment.CreateDate.Format(config.DateLayout),
	}
	return &commentPb
}

func BuildUserPb(user *models.User) *pb.User {
	var userPb pb.User
	err := copier.Copy(&userPb, &user)
	if err != nil {
		log.Println(err.Error())
	}
	return &userPb
}
