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
