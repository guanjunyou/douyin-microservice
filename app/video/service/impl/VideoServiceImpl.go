package impl

import (
	"context"
	"douyin-microservice/app/video/models"
	"douyin-microservice/app/video/rpc"
	"douyin-microservice/app/video/service"
	"douyin-microservice/config"
	"douyin-microservice/idl/pb"
	"douyin-microservice/pkg/utils"
	"github.com/jinzhu/copier"
	"github.com/zhenghaoz/gorse/client"
	"log"
	"strconv"
	"sync"
	"time"
)

type VideoServiceImpl struct {
	//service.UserService
	service.FavoriteService
}

func (videoService VideoServiceImpl) GetVideoListByLastTime(latestTime time.Time, userId int64) ([]models.VideoDVO, time.Time, error) {
	videolist, err := models.GetVideoListByLastTime(latestTime)
	if userId != -1 {
		// 用户登录，可以推荐视频
		recommendList := getRecommendVideos(userId)
		videolist = append(recommendList, videolist...)
	}

	size := len(videolist)
	var wg sync.WaitGroup
	VideoDVOList := make([]models.VideoDVO, size)
	if err != nil {
		return nil, time.Time{}, err
	}
	var err0 error
	for i := range videolist {
		var authorId = videolist[i].AuthorId
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// 通过 videoService 来调用 userService
			//user, err1 := videoService.UserService.GetUserById(authorId)
			//if err1 != nil {
			//	err0 = err1
			//	return
			//}
			var userReq pb.UserRequest
			userReq.UserId = authorId
			userResp, err1 := rpc.UserClient.GetUserById(context.Background(), &userReq)
			if err1 != nil {
				err0 = err1
				return
			}
			var videoDVO models.VideoDVO
			err2 := copier.Copy(&videoDVO, &videolist[i])
			if err2 != nil {
				err0 = err2
				return
			}
			videoDVO.Author = BuildUser(userResp.User)

			if userId != -1 {
				videoDVO.IsFavorite = favoriteService.FindIsFavouriteByUserIdAndVideoId(userId, videoDVO.Id)
			} else {
				videoDVO.IsFavorite = false
			}
			VideoDVOList[i] = videoDVO
		}(i)
	}

	wg.Wait()
	if err0 != nil {
		return nil, time.Time{}, err0
	}
	nextTime := time.Now()
	if len(videolist) != 0 {
		nextTime = videolist[len(videolist)-1].CreateDate
	}
	return VideoDVOList, nextTime, nil
}

func getRecommendVideos(userId int64) []models.Video {
	var recommend []models.Video
	ids, _ := utils.GorseClient.GetRecommend(context.Background(), strconv.FormatInt(userId, 10), "", 10)
	for i := range ids {
		videoId, _ := strconv.ParseInt(ids[i], 10, 64)
		video, _ := models.GetVideoById(videoId)
		recommend = append(recommend, video)
	}
	return recommend
}

// Publish 投稿接口
// TODO 借助redis协助实现feed流
func (videoService VideoServiceImpl) Publish(data []byte, userId int64, title string, filename string) error {
	//从title中过滤敏感词汇
	replaceTitle := utils.Filter.Replace(title, '#')
	//文件名
	//filename := filepath.Base(data.Filename)
	////将文件名拼接用户id
	//finalName := fmt.Sprintf("%d_%s", userId, filename)
	////保存文件的路径，暂时保存在本队public文件夹下
	//saveFile := filepath.Join("./public/", finalName)
	//保存视频在本地中
	// if err = c.SaveUploadedFile(data, saveFile); err != nil {
	coverName, err := utils.UploadToServer(data, filename)
	if err != nil {
		return err
	}
	user, err1 := models.GetUserById(userId)
	if err1 != nil {
		return nil
	}
	//将扩展名修改为.png并返回新的string作为封面文件名
	//ext := filepath.Ext(filename)
	//name := filename[:len(filename)-len(ext)]
	//coverName := name + ".png"
	//保存视频在数据库中
	video := models.Video{
		CommonEntity: utils.NewCommonEntity(),
		AuthorId:     userId,
		PlayUrl:      "http://" + config.Config.VideoServer.Addr2 + "/videos/" + filename,
		CoverUrl:     "http://" + config.Config.VideoServer.Addr2 + "/photos/" + coverName,
		Title:        replaceTitle,
	}
	err2 := models.SaveVideo(&video)
	if err2 != nil {
		return err2
	}
	//用户发布作品数加1
	user.WorkCount = user.WorkCount + 1
	err = models.UpdateUser(utils.GetMysqlDB(), user)
	if err != nil {
		return err
	}
	go InsertGorse(&video)
	return nil
}

func InsertGorse(video *models.Video) {
	//title := "合成高温超导体"
	videoId := video.Id
	timestamp := time.Unix(time.Now().Unix(), 0).UTC().Format(time.RFC3339)
	//x := gojieba.NewJieba()
	//defer x.Free()
	//labels := x.Cut(video.Title, true)
	item := client.Item{
		ItemId:   strconv.FormatInt(videoId, 10),
		IsHidden: false,
		//Labels:     labels,
		Categories: []string{"video"},
		Timestamp:  timestamp,
		Comment:    "",
	}
	utils.GorseClient.InsertItem(context.Background(), item)
}

// PublishList  发布列表
func (videoService VideoServiceImpl) PublishList(userId int64) ([]models.VideoDVO, error) {
	videoList, err := models.GetVediosByUserId(userId)
	if err != nil {
		return nil, err
	}
	size := len(videoList)
	VideoDVOList := make([]models.VideoDVO, size)
	//创建多个协程并发更新
	var wg sync.WaitGroup
	//接收协程产生的错误
	var err0 error
	for i := range videoList {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var userId = videoList[i].AuthorId
			var userReq pb.UserRequest
			userReq.UserId = userId
			userResp, err1 := rpc.UserClient.GetUserById(context.Background(), &userReq)
			if err1 != nil {
				err0 = err1
				return
			}
			var videoDVO models.VideoDVO
			err := copier.Copy(&videoDVO, &videoList[i])
			if err != nil {
				err0 = err1
			}
			videoDVO.Author = BuildUser(userResp.User)
			VideoDVOList[i] = videoDVO
		}(i)
	}
	wg.Wait()
	//处理协程内的错误
	if err0 != nil {
		return nil, err0
	}
	return VideoDVOList, nil
}

func BuildUser(userPb *pb.User) models.User {
	var user models.User
	err := copier.Copy(&user, &userPb)
	if err != nil {
		log.Println(err.Error())
	}
	return user
}

//以下是点赞功能
