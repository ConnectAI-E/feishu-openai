package services

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

type UserService struct {
	cache *cache.Cache
}

var userServices *UserService

func (u UserService) Get(userId string) string {
	// 获取用户的会话上下文
	sessionContext, ok := u.cache.Get(userId)
	if !ok {
		return ""
	}
	//list to string
	list := sessionContext.([]string)
	var result string
	for _, v := range list {
		result += v
	}
	return result
}

func (u UserService) Set(userId string, question, reply string) {
	// 列表，最多保存8个
	//如果满了，删除最早的一个
	//如果没有满，直接添加
	maxCache := 8
	maxLength := 2048
	maxCacheTime := time.Minute * 30
	listOut := make([]string, maxCache)
	value := fmt.Sprintf("Q:%s\nA:%s\n\n", question, reply)
	raw, ok := u.cache.Get(userId)
	if ok {
		listOut = raw.([]string)
		if len(listOut) == maxCache {
			listOut = listOut[1:]
		}
		listOut = append(listOut, value)
	} else {
		listOut = append(listOut, value)
	}

	//限制对话上下文长度
	if getStrPoolTotalLength(listOut) > maxLength {
		listOut = listOut[1:]
	}
	u.cache.Set(userId, listOut, maxCacheTime)
}

func (u UserService) Clear(userId string) bool {
	u.cache.Delete(userId)
	return true
}

type UserCacheInterface interface {
	Get(userId string) string
	Set(userId string, question, reply string)
	Clear(userId string) bool
}

func GetUserCache() UserCacheInterface {
	if userServices == nil {
		userServices = &UserService{cache: cache.New(30*time.Minute, 30*time.Minute)}
	}
	return userServices
}

func getStrPoolTotalLength(strPool []string) int {
	var total int
	for _, v := range strPool {
		total += len(v)
	}
	return total
}
