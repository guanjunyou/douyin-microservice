package impl

import (
	"context"
	"douyin-microservice/app/user/model"
	"douyin-microservice/app/user/mq"
	"douyin-microservice/app/user/rpc"
	"douyin-microservice/config"
	"douyin-microservice/idl/pb"
	"douyin-microservice/pkg/utils"
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"log"
	"sync/atomic"
	"time"
)

var userService UserServiceImpl
var once sync.Once

type UserServiceImpl struct {
}

func GetUserService() UserServiceImpl {
	once.Do(func() {
		userService = UserServiceImpl{}
	})
	return userService
}
func (userService UserServiceImpl) GetUserById(Id int64) (model.User, error) {
	result, err := model.GetUserById(Id)
	if err != nil {
		log.Printf("方法GetUserById() 失败 %v", err)
		return result, err
	}
	return result, nil
}

func (userService UserServiceImpl) GetUserByName(name string) (model.User, error) {
	result, err := model.GetUserByName(name)
	if err != nil {
		log.Printf("方法GetUserById() 失败 %v", err)
		return result, err
	}
	return result, nil
}

func (userService UserServiceImpl) Save(user model.User) error {
	return model.SaveUser(user)
}

/*
（
已完成
*/
func (userService UserServiceImpl) Register(username string, password string) (int64, string, error) {
	var userIdSequence = int64(1)
	_, errName := userService.GetUserByName(username)
	if errName == nil {
		//c.JSON(http.StatusBadRequest, UserLoginResponse{
		//	Response: models.Response{StatusCode: 1, StatusMsg: "用户名重复"},
		//})
		return -1, "", errors.New("用户名重复")
	}
	//var userRequest UserRequest
	//if err := c.ShouldBindJSON(&userRequest); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//username := userRequest.Username
	//password := userRequest.Password
	//加密
	encrypt, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	password = string(encrypt)

	atomic.AddInt64(&userIdSequence, 1)
	newUser := model.User{
		CommonEntity: utils.NewCommonEntity(),
		Name:         username,
		Password:     password,
	}

	err := userService.Save(newUser)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, UserLoginResponse{
		//	Response: models.Response{StatusCode: 1, StatusMsg: "Cant not save the User!"},
		//})
		return -1, "", err
	} else {
		token, err1 := utils.GenerateToken(username, newUser.CommonEntity)
		if err1 != nil {
			log.Printf("Can not get the token!")
			return -1, "", err1
		}
		err2 := utils.SaveTokenToRedis(newUser.Name, token, time.Duration(config.TokenTTL*float64(time.Second)))
		if err2 != nil {
			log.Printf("Fail : Save token in redis !")
			return -1, "", err2
		} else {
			//c.JSON(http.StatusOK, UserLoginResponse{
			//	Response: models.Response{StatusCode: 0},
			//	UserId:   newUser.Id,
			//	Token:    token,
			//})
			return newUser.Id, token, nil
		}
	}
}

/*
*
已完成
*/
func (userService UserServiceImpl) Login(username string, password string) (int64, string, error) {

	user, err := userService.GetUserByName(username)
	if err != nil {
		//c.JSON(http.StatusBadRequest, UserLoginResponse{
		//	Response: models.Response{StatusCode: 1, StatusMsg: "用户不存在，请注册!"},
		//})
		return -1, "", errors.New("用户不存在，请注册")
	}
	pwdErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if pwdErr != nil {
		//c.JSON(http.StatusOK, UserLoginResponse{
		//	Response: models.Response{StatusCode: 1, StatusMsg: "密码错误！"},
		//})
		return -1, "", errors.New("密码错误！")
	}

	token, err2 := utils.GenerateToken(username, user.CommonEntity)
	if err2 != nil {
		//c.JSON(http.StatusInternalServerError, UserLoginResponse{
		//	Response: models.Response{StatusCode: 1, StatusMsg: "生成token失败"},
		//})
		return -1, "", err2
	}

	err3 := utils.SaveTokenToRedis(user.Name, token, time.Duration(config.TokenTTL*float64(time.Second)))
	if err3 != nil {
		log.Printf("Fail : Save token in redis !")
		// TODO 开发完成后整理这个返回体 返回信息不能这么填
		//c.JSON(http.StatusInternalServerError, UserLoginResponse{
		//	Response: models.Response{StatusCode: 1, StatusMsg: "无法保存token 请检查redis连接"},
		//})
		return -1, "", err3
	}
	//
	//c.JSON(http.StatusOK, UserLoginResponse{
	//	Response: models.Response{StatusCode: 0, StatusMsg: "登录成功！"},
	//	UserId:   user.Id,
	//	Token:    token,
	//})
	return user.Id, token, nil
}

// LikeConsume  消费"userLikeMQ"中的消息
func (userService UserServiceImpl) LikeConsume(l *mq.LikeMQ) {
	_, err := l.Channel.QueueDeclare(l.QueueUserName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	//2、接收消息
	messages, err1 := l.Channel.Consume(
		l.QueueUserName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err1 != nil {
		panic(err1)
	}
	go userService.likeConsume(messages)
	//forever := make(chan bool)
	//log.Println(messages)

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")
}

// 点赞具体消费逻辑
func (userService UserServiceImpl) likeConsume(message <-chan amqp.Delivery) {
	for d := range message {
		jsonData := string(d.Body)
		log.Printf("user收到的消息为 %s\n", jsonData)
		data := mq.LikeMQToUser{}
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			panic(err)
		}
		userId := data.UserId
		tx := utils.GetMysqlDB().Begin()
		//获得当前用户
		user, err := model.GetUserById(userId)

		//查询视频作者
		author, err2 := model.GetUserById(data.AuthorId)
		if err2 != nil {
			panic(err2)
		}
		actionType := data.ActionType

		if actionType == 1 {
			//喜欢数量+一
			user.FavoriteCount = user.FavoriteCount + 1
			//如果是同一个作者，在同一个事务中必须保证针对同一行的操作只出现一次
			if user.Id == author.Id {
				user.TotalFavorited++
			}
			err = model.UpdateUser(tx, user)
			if err != nil {
				log.Println("err:", err)
				tx.Rollback()
				panic(err)
			}
			if user.Id != author.Id {
				//总点赞数+1
				author.TotalFavorited = author.TotalFavorited + 1
				err = model.UpdateUser(tx, author)
				if err != nil {
					log.Println("err:", err)
					tx.Rollback()
					panic(err)
				}
			}

		} else {
			//喜欢数量-1
			user.FavoriteCount = user.FavoriteCount - 1
			//如果是同一个作者，在同一个事务中必须保证针对同一行的操作只出现一次
			if user.Id == author.Id {
				user.TotalFavorited--
			}
			err = model.UpdateUser(tx, user)
			if err != nil {
				log.Println("err:", err)
				tx.Rollback()
				panic(err)
			}
			//总点赞数-1
			if user.Id != author.Id {
				author.TotalFavorited = author.TotalFavorited - 1
				err = model.UpdateUser(tx, author)
				if err != nil {
					log.Println("err:", err)
					tx.Rollback()
					panic(err)
				}
			}
		}
		tx.Commit()
	}
}

// 创建点赞消费者协程
func (userService UserServiceImpl) MakeLikeConsumers() {
	numConsumers := 20
	for i := 0; i < numConsumers; i++ {
		go userService.LikeConsume(mq.LikeRMQ)
	}
}

func (userService UserServiceImpl) UserInfo(userId int64, token string) (*model.User, error) {
	//userClaims, err := utils.AnalyseToken(token)
	//if err != nil || userClaims == nil {
	//	return nil, errors.New("用户未登录")
	//}
	user, err1 := userService.GetUserById(userId)
	if token != "" {
		userClaim, err := utils.AnalyseToken(token)
		if err != nil {
			return &user, nil
		}
		var req pb.CheckFollowRequest
		req.UserId = userClaim.CommonEntity.Id
		req.ToUserId = userId
		resp, _ := rpc.RelationClient.CheckFollowForUser(context.Background(), &req)
		user.IsFollow = resp.IsFollow
	}

	if err1 != nil {
		return nil, errors.New("用户不存在！")
	}
	return &user, nil
}

//type UserResponse struct {
//	models.Response
//	User models.User `json:"user"`
//}
//
//type UserRequest struct {
//	Username string `json:"username"`
//	Password string `json:"password"`
//}
//
//type UserLoginResponse struct {
//	models.Response
//	UserId int64  `json:"user_id,omitempty"`
//	Token  string `json:"token"`
//}

// FollowConsume  消费"followMQ"中的消息
func (userService UserServiceImpl) FollowConsume(followMQ *mq.FollowMQ) {
	_, err := followMQ.Channel.QueueDeclare(followMQ.QueueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	messages, err1 := followMQ.Channel.Consume(
		followMQ.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err1 != nil {
		panic(err1)
	}
	go userService.followConsume(messages)
	//forever := make(chan bool)
	//log.Println(messages)

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")
}

// 关注具体消费逻辑
func (userService UserServiceImpl) followConsume(message <-chan amqp.Delivery) {
	for d := range message {
		jsonData := string(d.Body)
		log.Printf("user收到的消息为 %s\n", jsonData)
		data := mq.FollowMQToUser{}
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			panic(err)
		}
		userId := data.UserId
		tx := utils.GetMysqlDB().Begin()
		//获得当前用户
		user, err := userService.GetUserById(userId)

		//查询视频作者
		toUser, err2 := userService.GetUserById(data.FollowUserId)
		if err2 != nil {
			panic(err2)
		}
		actionType := data.ActionType

		if actionType == 1 {
			//喜欢数量+一
			user.FollowCount = user.FollowCount + 1
			err = model.UpdateUser(tx, user)
			if err != nil {
				log.Println("err:", err)
				tx.Rollback()
				panic(err)
			}
			//总点赞数+1
			toUser.FollowerCount = toUser.FollowerCount + 1
			err = model.UpdateUser(tx, toUser)
			if err != nil {
				log.Println("err:", err)
				tx.Rollback()
				panic(err)
			}

		} else {
			//喜欢数量-1
			user.FollowCount = user.FollowCount - 1

			err = model.UpdateUser(tx, user)
			if err != nil {
				log.Println("err:", err)
				tx.Rollback()
				panic(err)
			}
			//总点赞数-1
			toUser.FollowerCount = toUser.FollowerCount - 1
			err = model.UpdateUser(tx, toUser)
			if err != nil {
				log.Println("err:", err)
				tx.Rollback()
				panic(err)
			}
		}
		tx.Commit()
	}
}

// 创建关注消费者协程
func (userService UserServiceImpl) MakeFollowConsumers() {
	numConsumers := 20
	for i := 0; i < numConsumers; i++ {
		go userService.FollowConsume(mq.FollowRMQ)
	}
}
