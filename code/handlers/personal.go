package handlers

import (
	"context"
	"fmt"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"start-feishubot/services"
)

type PersonalMessageHandler struct {
	userCache services.UserCacheInterface
	msgCache  services.MsgCacheInterface
}

func (p PersonalMessageHandler) handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {

	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	if p.msgCache.IfProcessed(*msgId) {
		fmt.Println("msgId", *msgId, "processed")
		return nil
	}
	p.msgCache.TagProcessed(*msgId)
	qParsed := parseContent(*content)
	sender := event.Event.Sender
	openId := sender.SenderId.OpenId
	cacheContent := p.userCache.Get(*openId)
	qEnd := qParsed
	if cacheContent != "" {
		qEnd = cacheContent + qParsed
	}
	ok := true
	completions, err := services.Completions(qEnd)
	if err != nil {
		return err
	}
	if len(completions) == 0 {
		ok = false
	}
	if ok {
		p.userCache.Set(*openId, qParsed, completions)
		sendMsg(ctx, completions, event.Event.Message.ChatId)
	}
	return nil

}

var _ MessageHandlerInterface = (*PersonalMessageHandler)(nil)

func NewPersonalMessageHandler() MessageHandlerInterface {
	return &PersonalMessageHandler{
		userCache: services.GetUserCache(),
		msgCache:  services.GetMsgCache(),
	}
}
