package services

import (
	"encoding/json"
	"time"

	"github.com/patrickmn/go-cache"
)

type SessionMode string

var (
	ModePicCreate SessionMode = "pic_create"
	ModePicVary   SessionMode = "pic_vary"
	ModeGPT       SessionMode = "gpt"
)

type SessionService struct {
	cache *cache.Cache
}

func (u SessionService) GetMode(sessionId string) SessionMode {
	// 获取用户的会话上下文
	sessionContext, ok := u.cache.Get(sessionId)
	if !ok {
		return ModeGPT
	}
	return sessionContext.(SessionMode)
}

func (u SessionService) SetMode(sessionId string, mode SessionMode) {
	maxCacheTime := time.Hour * 12

	u.cache.Set(sessionId, mode, maxCacheTime)
}

var sessionServices *SessionService

func (u SessionService) GetMsg(sessionId string) (msg []Messages) {
	// 获取用户的会话上下文
	sessionContext, ok := u.cache.Get(sessionId)
	if !ok {
		return msg
	}
	return sessionContext.([]Messages)
}

func (u SessionService) SetMsg(sessionId string, msg []Messages) {
	maxLength := 4096
	maxCacheTime := time.Hour * 12

	//限制对话上下文长度
	for getStrPoolTotalLength(msg) > maxLength {
		msg = append(msg[:1], msg[3:]...)
	}
	u.cache.Set(sessionId, msg, maxCacheTime)
}

func (u SessionService) Clear(sessionId string) bool {
	u.cache.Delete(sessionId)
	return true
}

type SessionServiceCacheInterface interface {
	GetMsg(sessionId string) []Messages
	SetMsg(sessionId string, msg []Messages)
	SetMode(sessionId string, mode SessionMode)
	GetMode(sessionId string) SessionMode
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
