package handlers

import (
	"context"

	"start-feishubot/services"
	"start-feishubot/services/openai"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

// AIModeChooseKind is the kind of card action for choosing AI mode
func NewAIModeCardHandler(cardMsg CardMsg,
	m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {

		if cardMsg.Kind == AIModeChooseKind {
			newCard, err, done := CommonProcessAIMode(cardMsg, cardAction,
				m.sessionCache)
			if done {
				return newCard, err
			}
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

// CommonProcessAIMode is the common process for choosing AI mode
func CommonProcessAIMode(msg CardMsg, cardAction *larkcard.CardAction,
	cache services.SessionServiceCacheInterface) (interface{},
	error, bool) {
	option := cardAction.Action.Option
	replyMsg(context.Background(), "已选择发散模式:"+option,
		&msg.MsgId)
	cache.SetAIMode(msg.SessionId, openai.AIModeMap[option])
	return nil, nil, true
}
