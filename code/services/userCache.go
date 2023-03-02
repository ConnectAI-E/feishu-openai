package services

import (
	"encoding/json"
	"time"

	"github.com/patrickmn/go-cache"
)

type UserService struct {
	cache *cache.Cache
}

var userServices *UserService

func (u UserService) Get(userId string) (msg []Messages) {
	// 获取用户的会话上下文
	sessionContext, ok := u.cache.Get(userId)
	if !ok {
		return msg
	}
	return sessionContext.([]Messages)
}

func (u UserService) Set(userId string, msg []Messages) {
	// 列表，最多保存8个
	//如果满了，删除最早的一个
	//如果没有满，直接添加
	maxCache := 16
	maxLength := 2048
	maxCacheTime := time.Minute * 30

	if len(msg) == maxCache {
		msg = msg[2:]
	}

	//限制对话上下文长度
	for getStrPoolTotalLength(msg) > maxLength {
		msg = msg[2:]
	}
	u.cache.Set(userId, msg, maxCacheTime)
}

func (u UserService) Clear(userId string) bool {
	u.cache.Delete(userId)
	return true
}

type UserCacheInterface interface {
	Get(userId string) []Messages
	Set(userId string, msg []Messages)
	Clear(userId string) bool
}

func GetUserCache() UserCacheInterface {
	if userServices == nil {
		userServices = &UserService{cache: cache.New(30*time.Minute, 30*time.Minute)}
	}
	return userServices
}

func getStrPoolTotalLength(strPool []Messages) int {
	var total int
	for _, v := range strPool {
		bytes, _ := json.Marshal(v)
		total += len(string(bytes))
	}
	return total
}
