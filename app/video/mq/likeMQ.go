package mq

import (
	"douyin-microservice/app/video/models"
	"douyin-microservice/config"
	"github.com/streadway/amqp"
)

type LikeMQ struct {
	RabbitMQ
	Channel        *amqp.Channel
	QueueUserName  string
	QueueVideoName string
	exchange       string
	key            string
}

type LikeMQToUser struct {
	UserId     int64 `json:"user_id"`
	VideoId    int64 `json:"video_id"`
	AuthorId   int64 `json:"author_id"`
	ActionType int   `json:"action_type"`
}

// 初始化 channel
// var LikeChannel chan models.LikeMQToVideo
var LikeChannel chan models.LikeMQToUser

func MakeLikeChannel() {
	ch := make(chan models.LikeMQToUser, config.BufferSize)
	LikeChannel = ch
}

// NewLikeRabbitMQ 获取likeMQ的对应队列。
func NewLikeRabbitMQ() *LikeMQ {
	likeMQ := &LikeMQ{
		RabbitMQ:       *Rmq,
		QueueUserName:  "userLikeMQ",
		QueueVideoName: "videoLikeMQ",
		exchange:       "likeExchange",
	}
	ch, err := likeMQ.conn.Channel()
	likeMQ.Channel = ch
	Rmq.failOnErr(err, "获取通道失败")
	return likeMQ
}

// Publish like操作的发布配置。
func (l *LikeMQ) Publish(message string) {
	//声明交换机
	err := l.Channel.ExchangeDeclare(
		//1.交换机名称
		l.exchange,
		//2、kind:交换机类型
		//	//amqp.ExchangeDirect 定向
		//	//amqp.ExchangeFanout 扇形（广播），发送消息到每个队列
		//	//amqp.ExchangeTopic 通配符的方式
		//	//amqp.ExchangeHeaders 参数匹配
		amqp.ExchangeFanout,
		//是否持久化
		true,
		//自动删除
		false,
		//内部使用
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	_, err = l.Channel.QueueDeclare(
		l.QueueUserName,
		//是否持久化
		true,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		panic(err)
	}
	//_, err = l.Channel.QueueDeclare(
	//	l.QueueVideoName,
	//	//是否持久化
	//	true,
	//	//是否为自动删除
	//	false,
	//	//是否具有排他性
	//	false,
	//	//是否阻塞
	//	false,
	//	//额外属性
	//	nil,
	//)
	//if err != nil {
	//	panic(err)
	//}
	//绑定队列和交换机
	err = l.Channel.QueueBind(l.QueueUserName, "", l.exchange, false, nil)
	if err != nil {
		panic(err)
	}
	//err = l.Channel.QueueBind(l.QueueVideoName, "", l.exchange, false, nil)
	//if err != nil {
	//	panic(err)
	//}

	err1 := l.Channel.Publish(
		l.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err1 != nil {
		panic(err)
	}

}

var LikeRMQ *LikeMQ

// InitLikeRabbitMQ 初始化rabbitMQ连接。
func InitLikeRabbitMQ() {
	LikeRMQ = NewLikeRabbitMQ()
	//LikeRMQ.Publish("hello word !")
	//go LikeRMQ.Consumer()
}
