package services

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type MsgService struct {
	cache *cache.Cache
}
type MsgCacheInterface interface {
	IfProcessed(msgId string) bool
	TagProcessed(msgId string)
	Clear(userId string) bool
}

var msgService *MsgService

func (u MsgService) IfProcessed(msgId string) bool {
	_, found := u.cache.Get(msgId)
	return found
}
func (u MsgService) TagProcessed(msgId string) {
	u.cache.Set(msgId, true, time.Minute*30)
}

func (u MsgService) Clear(userId string) bool {
	u.cache.Delete(userId)
	return true
}

func GetMsgCache() MsgCacheInterface {
	if msgService == nil {
		msgService = &MsgService{cache: cache.New(30*time.Minute, 30*time.Minute)}
	}
	return msgService
}
