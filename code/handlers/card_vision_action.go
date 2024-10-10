package handlers

import (
	"context"
	"fmt"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"start-feishubot/services"
)

func NewVisionResolutionHandler(cardMsg CardMsg,
	m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		if cardMsg.Kind == VisionStyleKind {
			CommonProcessVisionStyle(cardMsg, cardAction, m.sessionCache)
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}
func NewVisionModeChangeHandler(cardMsg CardMsg,
	m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		if cardMsg.Kind == VisionModeChangeKind {
			newCard, err, done := CommonProcessVisionModeChange(cardMsg, m.sessionCache)
			if done {
				return newCard, err
			}
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func CommonProcessVisionStyle(msg CardMsg,
	cardAction *larkcard.CardAction,
	cache services.SessionServiceCacheInterface) {
	option := cardAction.Action.Option
	fmt.Println(larkcore.Prettify(msg))
	cache.SetVisionDetail(msg.SessionId, services.VisionDetail(option))
	//send text
	replyMsg(context.Background(), "图片解析度调整为："+option,
		&msg.MsgId)
}

func CommonProcessVisionModeChange(cardMsg CardMsg,
	session services.SessionServiceCacheInterface) (
	interface{}, error, bool) {
	if cardMsg.Value == "1" {

		sessionId := cardMsg.SessionId
		session.Clear(sessionId)
		session.SetMode(sessionId,
			services.ModeVision)
		session.SetVisionDetail(sessionId,
			services.VisionDetailLow)

		newCard, _ :=
			newSendCard(
				withHeader("🕵️️ 已进入图片推理模式", larkcard.TemplateBlue),
				withVisionDetailLevelBtn(&sessionId),
				withNote("提醒：回复图片，让LLM和你一起推理图片的内容。"))
		return newCard, nil, true
	}
	if cardMsg.Value == "0" {
		newCard, _ := newSendCard(
			withHeader("️🎒 机器人提醒", larkcard.TemplateGreen),
			withMainMd("依旧保留此话题的上下文信息"),
			withNote("我们可以继续探讨这个话题,期待和您聊天。如果您有其他问题或者想要讨论的话题，请告诉我哦"),
		)
		return newCard, nil, true
	}
	return nil, nil, false
}
