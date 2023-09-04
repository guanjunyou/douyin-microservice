package utils

import (
	"github.com/zhenghaoz/gorse/client"
	"log"
)

var GorseClient *client.GorseClient

func InitGorse() {
	GorseClient = client.NewGorseClient("http://127.0.0.1:8088", "api_key")

	log.Println("gorse 连接成功")
}
