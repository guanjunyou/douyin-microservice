package models

import (
	"context"
	"douyin-microservice/config"
	utils2 "douyin-microservice/pkg/utils"
	"encoding/json"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type User struct {
	utils2.CommonEntity
	//Id            int64  `json:"id,omitempty"`
	Name            string `json:"name"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	Avatar          string `json:"avatar"`
	Gender          int    `json:"gender"`
	Age             int    `json:"age"`
	Nickname        string `json:"nickname"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
	IsFollow        bool   `json:"is_follow"`
	BackgroundImage string `json:"background_image"`
}

func (table *User) TableName() string {
	return "user"
}

func GetUserById(Id int64) (User, error) {
	var user User
	userKey := config.UserKey + strconv.FormatInt(Id, 10)
	userStr, errfind := utils2.GetRedisDB().Get(context.Background(), userKey).Result()
	if errfind == nil {
		errUnmarshal := json.Unmarshal([]byte(userStr), &user)
		if errUnmarshal != nil {
			return user, errUnmarshal
		}
		return user, nil
	}
	// 传参禁止直接字符串拼接，防止SQL注入
	err := utils2.GetMysqlDB().Where("id = ? AND is_deleted != ?", Id, 1).First(&user).Error
	if err != nil {
		return user, err
	}
	jsonStr, _ := json.Marshal(user)
	utils2.GetRedisDB().Set(context.Background(), userKey, jsonStr, time.Duration(config.UsedrKeyTTL)*time.Second)
	return user, nil
}

func GetUserByName(name string) (User, error) {
	var user User
	// 传参禁止直接字符串拼接，防止SQL注入
	err := utils2.GetMysqlDB().Where("name = ? AND is_deleted != ?", name, 1).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func SaveUser(user User) error {
	err := utils2.GetMysqlDB().Create(&user).Error
	if err != nil {
		return err
	}
	userStr, _ := json.Marshal(user)
	userKey := config.UserKey + strconv.FormatInt(user.Id, 10)
	utils2.GetRedisDB().Set(context.Background(), userKey, userStr, time.Duration(config.UsedrKeyTTL)*time.Second)
	return nil
}

func UpdateUser(tx *gorm.DB, user User) error {
	err := tx.Save(&user).Error
	if err != nil {
		return err
	}
	userStr, _ := json.Marshal(user)
	userKey := config.UserKey + strconv.FormatInt(user.Id, 10)
	utils2.GetRedisDB().Set(context.Background(), userKey, userStr, time.Duration(config.UsedrKeyTTL)*time.Second)
	return nil
}
