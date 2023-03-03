package handlers

import (
	"context"
	"fmt"
	"start-feishubot/services"
	"strings"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type PersonalMessageHandler struct {
	userCache services.UserCacheInterface
	msgCache  services.MsgCacheInterface
}

func (p PersonalMessageHandler) handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	sender := event.Event.Sender
	openId := sender.SenderId.OpenId
	chatId := event.Event.Message.ChatId
	if p.msgCache.IfProcessed(*msgId) {
		fmt.Println("msgId", *msgId, "processed")
		return nil
	}
	p.msgCache.TagProcessed(*msgId)
	qParsed := strings.Trim(parseContent(*content), " ")
	//// todo: test
	//if true {
	//	replyMarkdown(ctx, `![](https://open.feishu.cn/open-apis/block-kit/image/img_v2_041b28e3-5680-48c2-9af2-497ace79333g)`, msgId)
	//	return nil
	//}

	if len(qParsed) == 0 {
		sendMsg(ctx, "🤖️：你想知道什么呢~", chatId)
		fmt.Println("msgId", *msgId, "message.text is empty")
		return nil
	}

	if qParsed == "/clear" || qParsed == "清除" {
		p.userCache.Clear(*openId)
		sendMsg(ctx, "🤖️：AI机器人已清除记忆", chatId)
		return nil
	}

	msg := p.userCache.Get(*openId)
	msg = append(msg, services.Messages{
		Role: "user", Content: qParsed,
	})
	completions, err := services.Completions(msg)
	if err != nil {
		replyMsg(ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), msgId)
		return nil
	}
	msg = append(msg, completions)
	p.userCache.Set(*openId, msg)
	err = replyMsg(ctx, completions.Content, msgId)
	if err != nil {
		replyMsg(ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), msgId)
		return nil
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
