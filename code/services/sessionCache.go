package services

import (
	"encoding/json"
	"time"

	"github.com/patrickmn/go-cache"
)

type SessionService struct {
	cache *cache.Cache
}

var sessionServices *SessionService

func (u SessionService) Get(sessionId string) (msg []Messages) {
	// 获取用户的会话上下文
	sessionContext, ok := u.cache.Get(sessionId)
	if !ok {
		return msg
	}
	return sessionContext.([]Messages)
}

func (u SessionService) Set(sessionId string, msg []Messages) {
	maxLength := 4096
	maxCacheTime := time.Hour * 12

	//限制对话上下文长度
	for getStrPoolTotalLength(msg) > maxLength {
		msg = msg[2:]
	}
	u.cache.Set(sessionId, msg, maxCacheTime)
}

func (u SessionService) Clear(sessionId string) bool {
	u.cache.Delete(sessionId)
	return true
}

type SessionServiceCacheInterface interface {
	Get(sessionId string) []Messages
	Set(sessionId string, msg []Messages)
	Clear(sessionId string) bool
}

func GetSessionCache() SessionServiceCacheInterface {
	if sessionServices == nil {
		sessionServices = &SessionService{cache: cache.New(time.Hour*12, time.Hour*1)}
	}
	return sessionServices
}

func getStrPoolTotalLength(strPool []Messages) int {
	var total int
	for _, v := range strPool {
		bytes, _ := json.Marshal(v)
		total += len(string(bytes))
	}
	return total
}
