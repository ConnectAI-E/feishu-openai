package handlers

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"start-feishubot/services"
)

type PersonalMessageHandler struct {
	cache services.UserCacheInterface
}

func (p PersonalMessageHandler) handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {

	fmt.Println(larkcore.Prettify(event))
	content := event.Event.Message.Content
	q := parseContent(*content)
	fmt.Println("q", q)
	//sender := event.Event.Sender
	//openId := sender.SenderId.OpenId
	//cacheContent := p.cache.Get(*openId)
	qEnd := q
	//if cacheContent != "" {
	//	qEnd = cacheContent + q
	//}
	fmt.Println("qEnd", qEnd)
	ok := true
	reply, ok := services.GetAnswer(qEnd)
	fmt.Println("reply", reply, ok)
	if ok {
		sendMsg(ctx, reply, event.Event.Message.ChatId)
		//p.cache.Set(*openId, q, "nihao")
		return nil
	}

	return nil
}

var _ MessageHandlerInterface = (*PersonalMessageHandler)(nil)

func NewPersonalMessageHandler() MessageHandlerInterface {
	return &PersonalMessageHandler{
		cache: services.GetUserCache(),
	}
}
