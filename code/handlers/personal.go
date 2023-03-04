package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"start-feishubot/services"
	"start-feishubot/utils"
	"strings"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type PersonalMessageHandler struct {
	sessionCache services.SessionServiceCacheInterface
	msgCache     services.MsgCacheInterface
}

func (p PersonalMessageHandler) cardHandler(
	_ context.Context,
	cardAction *larkcard.CardAction) (interface{}, error) {
	var cardMsg CardMsg
	actionValue := cardAction.Action.Value
	actionValueJson, _ := json.Marshal(actionValue)
	json.Unmarshal(actionValueJson, &cardMsg)
	if cardMsg.Kind == ClearCardKind {
		newCard, err, done := CommonProcessClearCache(cardMsg, p.sessionCache)
		if done {
			return newCard, err
		}
	}
	return nil, nil
}

func CommonProcessClearCache(cardMsg CardMsg, session services.SessionServiceCacheInterface) (interface{},
	error,
	bool) {
	if cardMsg.Value == "1" {
		newCard, _ := newSendCard(
			withHeader("️👻 机器人提醒", larkcard.TemplateRed),
			withMainMd("已删除此话题的上下文信息"),
			withNote("我们可以开始一个全新的话题，继续找我聊天吧"),
		)
		session.Clear(cardMsg.SessionId)
		return newCard, nil, true
	}
	if cardMsg.Value == "0" {
		newCard, _ := newSendCard(
			withHeader("️👻 机器人提醒", larkcard.TemplateGreen),
			withMainMd("依旧保留此话题的上下文信息"),
			withNote("我们可以继续探讨这个话题,期待和您聊天。如果您有其他问题或者想要讨论的话题，请告诉我哦"),
		)
		return newCard, nil, true
	}
	return nil, nil, false
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
		sendClearCacheCheckCard(ctx, sessionId, msgId)
		return nil
	}

	system, foundSystem := utils.EitherCutPrefix(qParsed, "/system ",
		"角色扮演 ")
	if foundSystem {
		p.sessionCache.Clear(*sessionId)
		systemMsg := append([]services.Messages{}, services.Messages{
			Role: "system", Content: system,
		})
		p.sessionCache.Set(*sessionId, systemMsg)
		sendSystemInstructionCard(ctx, sessionId, msgId, system)
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
	//if new topic
	if len(msg) == 2 {
		sendNewTopicCard(ctx, sessionId, msgId, completions.Content)
		return nil
	}
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
