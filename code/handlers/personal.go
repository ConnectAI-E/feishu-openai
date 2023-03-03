package handlers

import (
	"context"
	"fmt"
	"start-feishubot/services"
	"strings"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type PersonalMessageHandler struct {
	sessionCache services.SessionServiceCacheInterface
	msgCache     services.MsgCacheInterface
}

func (p PersonalMessageHandler) handle(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	rootId := event.Event.Message.RootId
	chatId := event.Event.Message.ChatId
	sessionId := rootId
	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}
	if p.msgCache.IfProcessed(*msgId) {
		fmt.Println("msgId", *msgId, "processed")
		return nil
	}
	p.msgCache.TagProcessed(*msgId)
	qParsed := strings.Trim(parseContent(*content), " ")
	if len(qParsed) == 0 {
		sendMsg(ctx, "🤖️：你想知道什么呢~", chatId)
		fmt.Println("msgId", *msgId, "message.text is empty")
		return nil
	}

	if qParsed == "/clear" || qParsed == "清除" {
		p.sessionCache.Clear(*sessionId)
		sendMsg(ctx, "🤖️：AI机器人已清除记忆", chatId)
		//p.userCache.Clear(*openId)
		//sendMsg(ctx, "🤖️：AI机器人已清除记忆", chatId)

		sendClearCacheCheckCard(ctx, chatId)

		return nil
	}

	msg := p.sessionCache.Get(*sessionId)
	msg = append(msg, services.Messages{
		Role: "user", Content: qParsed,
	})
	completions, err := services.Completions(msg)
	if err != nil {
		replyMsg(ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), msgId)
		return nil
	}
	msg = append(msg, completions)
	p.sessionCache.Set(*sessionId, msg)
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
		sessionCache: services.GetSessionCache(),
		msgCache:     services.GetMsgCache(),
	}
}
