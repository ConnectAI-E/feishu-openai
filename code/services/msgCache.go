package services

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type MsgService struct {
	cache *cache.Cache
}

var msgService *MsgService

func (u MsgService) IfProcessed(msgId string) bool {
	get, b := u.cache.Get(msgId)
	if !b {
		return false
	}
	return get.(bool)
}
func (u MsgService) TagProcessed(msgId string) {
	u.cache.Set(msgId, true, time.Minute*30)
}

func (u MsgService) Clear(userId string) bool {
	u.cache.Delete(userId)
	return true
}

type MsgCacheInterface interface {
	IfProcessed(msg string) bool
	TagProcessed(msg string)
}

func GetMsgCache() MsgCacheInterface {
	if msgService == nil {
		msgService = &MsgService{cache: cache.New(30*time.Minute, 30*time.Minute)}
	}
	return msgService
}
