package handlers

import (
	"context"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"start-feishubot/logger"
	"start-feishubot/services"
)

func NewClearCardHandler(cardMsg CardMsg, m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		if cardMsg.Kind == ClearCardKind {
			newCard, err, done := CommonProcessClearCache(cardMsg, m.sessionCache)
			if done {
				return newCard, err
			}
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func CommonProcessClearCache(cardMsg CardMsg, session services.SessionServiceCacheInterface) (
	interface{}, error, bool) {
	logger.Debugf("card msg value %v", cardMsg.Value)
	if cardMsg.Value == "1" {
		session.Clear(cardMsg.SessionId)
		newCard, _ := newSendCard(
			withHeader("️🆑 机器人提醒", larkcard.TemplateGrey),
			withMainMd("已删除此话题的上下文信息"),
			withNote("我们可以开始一个全新的话题，继续找我聊天吧"),
		)
		logger.Debugf("session %v", newCard)
		return newCard, nil, true
	}
	if cardMsg.Value == "0" {
		newCard, _ := newSendCard(
			withHeader("️🆑 机器人提醒", larkcard.TemplateGreen),
			withMainMd("依旧保留此话题的上下文信息"),
			withNote("我们可以继续探讨这个话题,期待和您聊天。如果您有其他问题或者想要讨论的话题，请告诉我哦"),
		)
		return newCard, nil, true
	}
	return nil, nil, false
}
