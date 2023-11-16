package handlers

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"start-feishubot/logger"

	"start-feishubot/services"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

func NewPicResolutionHandler(cardMsg CardMsg, m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		if cardMsg.Kind == PicResolutionKind {
			CommonProcessPicResolution(cardMsg, cardAction, m.sessionCache)
			return nil, nil
		}
		if cardMsg.Kind == PicStyleKind {
			CommonProcessPicStyle(cardMsg, cardAction, m.sessionCache)
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func NewPicModeChangeHandler(cardMsg CardMsg, m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		if cardMsg.Kind == PicModeChangeKind {
			newCard, err, done := CommonProcessPicModeChange(cardMsg, m.sessionCache)
			if done {
				return newCard, err
			}
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func NewPicTextMoreHandler(cardMsg CardMsg, m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		if cardMsg.Kind == PicTextMoreKind {
			go func() {
				m.CommonProcessPicMore(cardMsg)
			}()
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func CommonProcessPicResolution(msg CardMsg,
	cardAction *larkcard.CardAction,
	cache services.SessionServiceCacheInterface) {
	option := cardAction.Action.Option
	fmt.Println(larkcore.Prettify(msg))
	cache.SetPicResolution(msg.SessionId, services.Resolution(option))
	//send text
	replyMsg(context.Background(), "已更新图片分辨率为"+option,
		&msg.MsgId)
}

func CommonProcessPicStyle(msg CardMsg,
	cardAction *larkcard.CardAction,
	cache services.SessionServiceCacheInterface) {
	option := cardAction.Action.Option
	fmt.Println(larkcore.Prettify(msg))
	cache.SetPicStyle(msg.SessionId, services.PicStyle(option))
	//send text
	replyMsg(context.Background(), "已更新图片风格为"+option,
		&msg.MsgId)
}

func (m MessageHandler) CommonProcessPicMore(msg CardMsg) {
	resolution := m.sessionCache.GetPicResolution(msg.SessionId)
	style := m.sessionCache.GetPicStyle(msg.SessionId)

	logger.Debugf("resolution: %v", resolution)
	logger.Debug("msg: %v", msg)
	question := msg.Value.(string)
	bs64, _ := m.gpt.GenerateOneImage(question, resolution, style)
	replayImageCardByBase64(context.Background(), bs64, &msg.MsgId,
		&msg.SessionId, question)
}

func CommonProcessPicModeChange(cardMsg CardMsg,
	session services.SessionServiceCacheInterface) (
	interface{}, error, bool) {
	if cardMsg.Value == "1" {

		sessionId := cardMsg.SessionId
		session.Clear(sessionId)
		session.SetMode(sessionId,
			services.ModePicCreate)
		session.SetPicResolution(sessionId,
			services.Resolution256)

		newCard, _ :=
			newSendCard(
				withHeader("🖼️ 已进入图片创作模式", larkcard.TemplateBlue),
				withPicResolutionBtn(&sessionId),
				withNote("提醒：回复文本或图片，让AI生成相关的图片。"))
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
