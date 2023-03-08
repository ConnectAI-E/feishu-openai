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

type Resolution string

const (
	Resolution256  Resolution = "256x256"
	Resolution512  Resolution = "512x512"
	Resolution1024 Resolution = "1024x1024"
)

type PicSetting struct {
	resolution Resolution
}

type SessionMeta struct {
	Mode       SessionMode
	Msg        []Messages
	PicSetting PicSetting
}

var sessionServices *SessionService

func (s *SessionService) GetMode(sessionID string) SessionMode {
	// Get the session mode from the cache.
	sessionContext, ok := s.cache.Get(sessionID)
	if !ok {
		return ModeGPT
	}
	sessionMeta := sessionContext.(*SessionMeta)
	return sessionMeta.Mode
}

func (s *SessionService) SetMode(sessionID string, mode SessionMode) {
	// Update the session mode in the cache.
	maxCacheTime := time.Hour * 12

	sessionContext, ok := s.cache.Get(sessionID)
	if !ok {
		sessionMeta := &SessionMeta{Mode: mode}
		s.cache.Set(sessionID, sessionMeta, maxCacheTime)
		return
	}
	sessionMeta := sessionContext.(*SessionMeta)
	sessionMeta.Mode = mode
	s.cache.Set(sessionID, sessionMeta, maxCacheTime)
}

func (s *SessionService) GetMsg(sessionId string) (msg []Messages) {
	sessionContext, ok := s.cache.Get(sessionId)
	if !ok {
		return nil
	}
	sessionMeta := sessionContext.(*SessionMeta)
	return sessionMeta.Msg
}

func (s *SessionService) SetMsg(sessionId string, msg []Messages) {
	maxLength := 4096
	maxCacheTime := time.Hour * 12

	//限制对话上下文长度
	for getStrPoolTotalLength(msg) > maxLength {
		msg = append(msg[:1], msg[2:]...)
	}

	sessionContext, ok := s.cache.Get(sessionId)
	if !ok {
		sessionMeta := &SessionMeta{Msg: msg}
		s.cache.Set(sessionId, sessionMeta, maxCacheTime)
		return
	}
	sessionMeta := sessionContext.(*SessionMeta)
	sessionMeta.Msg = msg
	s.cache.Set(sessionId, sessionMeta, maxCacheTime)
}

func (s *SessionService) SetPicResolution(sessionId string,
	resolution Resolution) {
	maxCacheTime := time.Hour * 12

	//if not in [Resolution256, Resolution512, Resolution1024] then set
	//to Resolution256
	switch resolution {
	case Resolution256, Resolution512, Resolution1024:
	default:
		resolution = Resolution256
	}

	sessionContext, ok := s.cache.Get(sessionId)
	if !ok {
		sessionMeta := &SessionMeta{PicSetting: PicSetting{resolution: resolution}}
		s.cache.Set(sessionId, sessionMeta, maxCacheTime)
		return
	}
	sessionMeta := sessionContext.(*SessionMeta)
	sessionMeta.PicSetting.resolution = resolution
	s.cache.Set(sessionId, sessionMeta, maxCacheTime)
}

func (s *SessionService) GetPicResolution(sessionId string) string {
	sessionContext, ok := s.cache.Get(sessionId)
	if !ok {
		return string(Resolution256)
	}
	sessionMeta := sessionContext.(*SessionMeta)
	return string(sessionMeta.PicSetting.resolution)

}

func (s *SessionService) Clear(sessionID string) {
	// Delete the session context from the cache.
	s.cache.Delete(sessionID)
}

type SessionServiceCacheInterface interface {
	GetMsg(sessionId string) []Messages
	SetMsg(sessionId string, msg []Messages)
	SetMode(sessionId string, mode SessionMode)
	GetMode(sessionId string) SessionMode
	SetPicResolution(sessionId string, resolution Resolution)
	GetPicResolution(sessionId string) string
	Clear(sessionId string)
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
