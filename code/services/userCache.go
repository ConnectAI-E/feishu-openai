package services

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type UserService struct {
	cache *cache.Cache
}

var userServices *UserService

func (u UserService) Get(userId string) string {
	get, b := u.cache.Get(userId)
	if !b {
		return ""
	}
	return get.(string)
}

func (u UserService) Set(userId string, question, reply string) {
	value := question + "\n" + reply
	u.cache.Set(userId, value, time.Minute*5)
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
		userServices = &UserService{cache: cache.New(10*time.Minute, 10*time.Minute)}
	}
	return userServices
}
